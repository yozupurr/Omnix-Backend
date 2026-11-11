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

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /friends/api/public/friends/{id}", list)
	mux.HandleFunc("GET /friends/api/public/blocklist/{id}", blocked)
	mux.HandleFunc("GET /friends/api/public/list/fortnite/{id}/recentPlayers", recentPlayers)
	mux.HandleFunc("GET /friends/api/v1/{id}/settings", settings)
}

func settings(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, preferences{AcceptInvites: "public"})
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
