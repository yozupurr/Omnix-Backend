package friends

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

type blocklist struct {
	BlockedUsers []string `json:"blockedUsers"`
}

type preferences struct {
	AcceptInvites string `json:"acceptInvites"`
}

type summary struct {
	Friends   []any       `json:"friends"`
	Incoming  []any       `json:"incoming"`
	Outgoing  []any       `json:"outgoing"`
	Suggested []any       `json:"suggested"`
	Blocklist []any       `json:"blocklist"`
	Settings  preferences `json:"settings"`
}

var publicPreferences = preferences{AcceptInvites: "public"}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /friends/api/public/friends/{id}", list)
	mux.HandleFunc("GET /friends/api/public/blocklist/{id}", blocked)
	mux.HandleFunc("GET /friends/api/public/list/fortnite/{id}/recentPlayers", recentPlayers)
	mux.HandleFunc("GET /friends/api/v1/{id}/settings", settings)
	mux.HandleFunc("GET /friends/api/v1/{id}/summary", summarize)
	mux.HandleFunc("GET /friends/api/v1/{id}/blocklist", blocklistV1)
	mux.HandleFunc("GET /friends/api/v1/{id}/recent/fortnite", recentPlayers)
}

func settings(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, publicPreferences)
}

func summarize(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, summary{
		Friends:   []any{},
		Incoming:  []any{},
		Outgoing:  []any{},
		Suggested: []any{},
		Blocklist: []any{},
		Settings:  publicPreferences,
	})
}

func blocklistV1(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}

func list(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}

func blocked(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, blocklist{BlockedUsers: []string{}})
}

func recentPlayers(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}
