package xmpp

import (
	"errors"
	"log/slog"
	"net"
	"strings"
)

const adminUser = "xmpp-admin"

type Authenticator func(token string) (accountID, displayName string, ok bool)

type Server struct {
	listener     net.Listener
	hub          *hub
	log          *slog.Logger
	authenticate Authenticator
}

func Listen(addr string, log *slog.Logger, authenticate Authenticator) (*Server, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{listener: listener, hub: newHub(), log: log, authenticate: authenticate}, nil
}

func (s *Server) Serve() error {
	for {
		raw, err := s.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		if err != nil {
			return err
		}
		go (&session{server: s, conn: &conn{raw: raw}}).run()
	}
}

func (s *Server) Close() {
	s.listener.Close()
	s.hub.closeAll()
}

func (s *Server) Send(accountID, body string) {
	c := s.hub.client(accountID)
	if c == nil {
		return
	}
	_, host, _ := strings.Cut(c.jid, "@")
	domain, _, _ := strings.Cut(host, "/")
	c.send(message(c.jid, adminUser+"@"+domain, "", "", body))
}
