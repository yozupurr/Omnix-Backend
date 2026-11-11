package xmpp

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"
)

const (
	defaultDomain = "prod.ol.epicgames.com"
	roomPrefix    = "muc."
	partyPrefix   = "party-"
	maxChatLength = 300
)

type session struct {
	server *Server
	conn   *conn
	domain string

	streamID      string
	accountID     string
	displayName   string
	authenticated bool
	client        *client
}

func (s *session) run() {
	defer s.cleanup()
	decoder := xml.NewDecoder(s.conn.raw)
	for {
		token, err := decoder.Token()
		if err != nil {
			s.server.log.Info("xmpp stream ended", "reason", err.Error(), "account", s.accountID)
			return
		}
		switch t := token.(type) {
		case xml.EndElement:
			if t.Name.Local == "stream" {
				return
			}
		case xml.StartElement:
			if t.Name.Local == "stream" {
				s.open(t)
				continue
			}
			var n node
			if err := decoder.DecodeElement(&n, &t); err != nil {
				s.server.log.Warn("xmpp invalid stanza", "name", t.Name.Local, "reason", err.Error(), "account", s.accountID)
				s.fail()
				return
			}
			if !s.dispatch(n) {
				s.server.log.Warn("xmpp stanza rejected", "name", n.XMLName.Local, "type", n.attr("type"), "to", n.attr("to"), "account", s.accountID)
				s.fail()
				return
			}
		}
	}
}

func (s *session) fail() {
	_ = s.conn.write(streamFooter)
}

func (s *session) cleanup() {
	if s.client != nil {
		for _, room := range s.server.hub.remove(s.client) {
			s.announceExit(room)
		}
		s.server.log.Info("xmpp disconnected", "account", s.accountID, "name", s.displayName)
	}
	s.conn.close()
}

func (s *session) send(e element) {
	_ = s.conn.write(e.String())
}

func (s *session) dispatch(n node) bool {
	switch n.XMLName.Local {
	case "auth":
		return s.auth(n)
	case "iq":
		return s.iq(n)
	case "message":
		return s.message(n)
	case "presence":
		return s.presence(n)
	}
	s.server.log.Warn("xmpp unhandled stanza", "name", n.XMLName.Local, "account", s.accountID)
	return true
}

func (s *session) open(start xml.StartElement) {
	if s.streamID == "" {
		s.streamID = randomID()
		s.domain = defaultDomain
		for _, a := range start.Attr {
			if a.Name.Local == "to" && a.Value != "" {
				s.domain = a.Value
			}
		}
	}
	s.server.log.Debug("xmpp stream opened", "remote", s.conn.raw.RemoteAddr().String(), "authenticated", s.authenticated)
	_ = s.conn.write(streamHeader(s.domain, s.streamID))
	s.send(features(s.authenticated))
}

func (s *session) auth(n node) bool {
	if s.authenticated {
		return false
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimSpace(n.Text))
	if err != nil {
		return s.reject("invalid base64", n)
	}
	parts := strings.Split(string(decoded), "\x00")
	if len(parts) != 3 {
		return s.reject("unexpected credentials layout", n)
	}
	accountID, displayName, ok := s.server.authenticate(parts[2])
	if !ok {
		return s.reject("unknown token", n)
	}
	s.accountID, s.displayName, s.authenticated = accountID, displayName, true
	s.send(success())
	return true
}

func (s *session) reject(reason string, n node) bool {
	s.server.log.Warn("xmpp auth rejected", "reason", reason, "mechanism", n.attr("mechanism"))
	s.send(authFailure())
	return false
}

func (s *session) iq(n node) bool {
	if !s.authenticated {
		return false
	}
	id := n.attr("id")
	if bind, ok := n.child("bind"); ok {
		return s.bind(id, bind)
	}
	if s.client == nil {
		return false
	}
	s.send(iqResult(s.domain, s.client.jid, id))
	return true
}

func (s *session) bind(id string, bind node) bool {
	resource, ok := bind.child("resource")
	name := strings.TrimSpace(resource.Text)
	if s.client != nil || !ok || name == "" {
		return false
	}
	s.client = &client{
		conn:        s.conn,
		accountID:   s.accountID,
		displayName: s.displayName,
		resource:    name,
		jid:         s.accountID + "@" + s.domain + "/" + name,
	}
	s.server.hub.add(s.client)
	s.server.log.Info("xmpp connected", "account", s.accountID, "name", s.displayName)
	s.send(bindResult(s.client.jid, id))
	return true
}

func (s *session) message(n node) bool {
	if s.client == nil {
		return false
	}
	body, ok := n.child("body")
	if !ok || body.Text == "" {
		return true
	}
	to := n.attr("to")
	switch n.attr("type") {
	case "chat":
		s.chat(to, body.Text)
	case "groupchat":
		s.groupChat(to, body.Text)
	default:
		s.relay(to, n.attr("id"), body.Text)
	}
	return true
}

func (s *session) chat(to, body string) {
	if to == "" || len(body) >= maxChatLength {
		return
	}
	receiver := s.server.hub.find(to)
	if receiver == nil || receiver == s.client {
		return
	}
	receiver.send(message(receiver.jid, s.client.jid, "", "chat", body))
}

func (s *session) groupChat(to, body string) {
	room, ok := s.roomName(to)
	if !ok || len(body) >= maxChatLength {
		return
	}
	members := s.server.hub.members(room)
	if !contains(members, s.client) {
		return
	}
	from := s.roomJID(room, s.client.nick())
	for _, m := range members {
		m.send(message(m.jid, from, "", "groupchat", body))
	}
}

func (s *session) relay(to, id, body string) {
	if to == "" || !isObject(body) {
		return
	}
	if receiver := s.server.hub.find(to); receiver != nil {
		receiver.send(message(receiver.jid, s.client.jid, id, "", body))
	}
}

func (s *session) presence(n node) bool {
	if s.client == nil {
		return false
	}
	to := n.attr("to")
	if n.attr("type") == "unavailable" {
		if room, ok := s.roomName(to); ok {
			s.leave(room)
		}
		return true
	}
	if _, ok := n.child("x"); ok {
		if room, ok := s.roomName(to); ok {
			s.join(room)
		}
		return true
	}
	status, ok := n.child("status")
	if !ok || !isObject(status.Text) {
		return true
	}
	_, away := n.child("show")
	s.send(selfPresence(s.client.jid, status.Text, away))
	return true
}

func (s *session) join(room string) {
	members, joined := s.server.hub.join(room, s.client)
	if !joined {
		return
	}
	me := s.client
	s.send(roomPresence(me.jid, s.roomJID(room, me.nick()), "participant", me.nick(), me.jid, false, "110", "100", "170", "201"))
	for _, m := range members {
		s.send(roomPresence(me.jid, s.roomJID(room, m.nick()), "participant", m.nick(), m.jid, false))
		if m != me {
			m.send(roomPresence(m.jid, s.roomJID(room, me.nick()), "participant", me.nick(), me.jid, false))
		}
	}
}

func (s *session) leave(room string) {
	if !s.server.hub.leave(room, s.client) {
		return
	}
	me := s.client
	s.send(roomPresence(me.jid, s.roomJID(room, me.nick()), "none", me.nick(), me.jid, true, "110", "100", "170"))
}

func (s *session) announceExit(room string) {
	if !strings.HasPrefix(strings.ToLower(room), partyPrefix) {
		return
	}
	body, err := json.Marshal(memberExited{
		Type: "com.epicgames.party.memberexited",
		Payload: memberExitedPayload{
			PartyID:  room[len(partyPrefix):],
			MemberID: s.client.accountID,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
	})
	if err != nil {
		return
	}
	for _, m := range s.server.hub.members(room) {
		m.send(message(m.jid, s.client.jid, strings.ToUpper(randomID()), "", string(body)))
	}
}

type memberExitedPayload struct {
	PartyID   string `json:"partyId"`
	MemberID  string `json:"memberId"`
	WasKicked bool   `json:"wasKicked"`
}

type memberExited struct {
	Type      string              `json:"type"`
	Payload   memberExitedPayload `json:"payload"`
	Timestamp string              `json:"timestamp"`
}

func (s *session) roomName(jid string) (string, bool) {
	bare, _, _ := strings.Cut(jid, "/")
	name, host, ok := strings.Cut(bare, "@")
	return name, ok && name != "" && host == roomPrefix+s.domain
}

func (s *session) roomJID(room, nick string) string {
	return room + "@" + roomPrefix + s.domain + "/" + nick
}

func contains(members []*client, c *client) bool {
	for _, m := range members {
		if m == c {
			return true
		}
	}
	return false
}

func isObject(raw string) bool {
	var v map[string]any
	return json.Unmarshal([]byte(raw), &v) == nil
}

func randomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
