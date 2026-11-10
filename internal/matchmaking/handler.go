package matchmaking

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/matchmaking/session/findPlayer/{id}", findPlayer)
}

func findPlayer(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}
