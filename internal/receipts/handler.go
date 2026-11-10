package receipts

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/receipts/v1/account/{id}/receipts", list)
}

func list(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}
