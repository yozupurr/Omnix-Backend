package versioncheck

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

type response struct {
	Type string `json:"type"`
}

var noUpdate = response{Type: "NO_UPDATE"}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/versioncheck", check)
}

func check(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, noUpdate)
}
