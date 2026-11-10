package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"omnix/internal/account"
	"omnix/internal/auth"
	"omnix/internal/calendar"
	"omnix/internal/catalog"
	"omnix/internal/config"
	"omnix/internal/content"
	"omnix/internal/datarouter"
	"omnix/internal/eulatracking"
	"omnix/internal/friends"
	"omnix/internal/game"
	"omnix/internal/lightswitch"
	"omnix/internal/matchmaking"
	"omnix/internal/platform/httpx"
	"omnix/internal/profile"
	"omnix/internal/receipts"
	"omnix/internal/storage"
	"omnix/internal/storefront"
	"omnix/internal/version"
	"omnix/internal/versioncheck"
	"omnix/internal/xmpp"
)

type Server struct {
	http *http.Server
	xmpp *xmpp.Server
	log  *slog.Logger
	addr string
}

func New(cfg config.Config, log *slog.Logger) (*Server, error) {
	cosmetics, err := catalog.Load(cfg.CatalogFile)
	if err != nil {
		return nil, err
	}
	log.Info("catalog loaded", "items", cosmetics.Len())

	accounts := account.NewStore()

	mux := http.NewServeMux()
	account.Register(mux, accounts)
	authenticate := auth.Register(mux, accounts)
	calendar.Register(mux)
	content.Register(mux)
	datarouter.Register(mux)
	eulatracking.Register(mux)
	friends.Register(mux)
	game.Register(mux)
	lightswitch.Register(mux)
	matchmaking.Register(mux)
	profile.Register(mux, cosmetics)
	receipts.Register(mux)
	storefront.Register(mux)
	versioncheck.Register(mux)
	if err := storage.Register(mux, cfg.StorageDir); err != nil {
		return nil, err
	}
	mux.HandleFunc("/", notFound)

	chat, err := xmpp.Listen(cfg.XMPPAddr, log, authenticate)
	if err != nil {
		return nil, err
	}

	return &Server{
		http: &http.Server{
			Addr:              cfg.Addr,
			Handler:           version.Middleware(requestLogger(log)(mux)),
			ReadHeaderTimeout: 10 * time.Second,
		},
		xmpp: chat,
		log:  log,
		addr: cfg.Addr,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	failed := make(chan error, 2)
	go func() {
		s.log.Info("listening", "addr", s.addr)
		if err := s.http.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			failed <- err
		}
	}()
	go func() {
		if err := s.xmpp.Serve(); err != nil {
			failed <- err
		}
	}()

	var result error
	select {
	case <-ctx.Done():
	case result = <-failed:
	}

	s.xmpp.Close()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.http.Shutdown(shutdownCtx); err != nil && result == nil {
		result = err
	}
	return result
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	httpx.NotFound.Write(w)
}
