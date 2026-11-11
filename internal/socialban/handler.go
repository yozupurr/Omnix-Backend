package socialban

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

type status struct {
	Bans     []any `json:"bans"`
	Warnings []any `json:"warnings"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /socialban/api/public/v1/{id}", bans)
}

func bans(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, status{Bans: []any{}, Warnings: []any{}})
}
