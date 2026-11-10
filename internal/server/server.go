package server

import (
	"log/slog"
	"net/http"
	"time"

	"omnix/internal/account"
	"omnix/internal/auth"
	"omnix/internal/calendar"
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
)

func New(cfg config.Config, log *slog.Logger) (*http.Server, error) {
	accounts := account.NewStore()

	mux := http.NewServeMux()
	account.Register(mux, accounts)
	auth.Register(mux, accounts)
	calendar.Register(mux)
	content.Register(mux)
	datarouter.Register(mux)
	eulatracking.Register(mux)
	friends.Register(mux)
	game.Register(mux)
	lightswitch.Register(mux)
	matchmaking.Register(mux)
	profile.Register(mux)
	receipts.Register(mux)
	storefront.Register(mux)
	versioncheck.Register(mux)
	if err := storage.Register(mux, cfg.StorageDir); err != nil {
		return nil, err
	}
	mux.HandleFunc("/", notFound)

	return &http.Server{
		Addr:              cfg.Addr,
		Handler:           version.Middleware(requestLogger(log)(mux)),
		ReadHeaderTimeout: 10 * time.Second,
	}, nil
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	httpx.NotFound.Write(w)
}
