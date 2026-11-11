package storefront

import (
	"net/http"
	"time"

	"omnix/internal/catalog"
	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

type handler struct {
	cosmetics *catalog.Catalog
}

func Register(mux *http.ServeMux, cosmetics *catalog.Catalog) {
	h := handler{cosmetics: cosmetics}
	mux.HandleFunc("GET /fortnite/api/storefront/v2/catalog", h.catalog)
}

func (h handler) catalog(w http.ResponseWriter, r *http.Request) {
	season := version.From(r.Context()).Season()
	httpx.JSON(w, http.StatusOK, newShop(h.cosmetics.Available(max(season-1, 1)), season, time.Now()))
}
