package channels

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/user/setting", userSettings)
}

func userSettings(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}
