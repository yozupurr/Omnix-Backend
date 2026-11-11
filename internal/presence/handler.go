package presence

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /presence/api/v1/_/{id}/last-online", lastOnline)
}

func lastOnline(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, map[string]any{})
}
