package storefront

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

type catalog struct {
	RefreshIntervalHrs int    `json:"refreshIntervalHrs"`
	DailyPurchaseHrs   int    `json:"dailyPurchaseHrs"`
	Expiration         string `json:"expiration"`
	Storefronts        []any  `json:"storefronts"`
}

var empty = catalog{
	RefreshIntervalHrs: 1,
	DailyPurchaseHrs:   24,
	Expiration:         "9999-12-31T23:59:59.999Z",
	Storefronts:        []any{},
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/storefront/v2/catalog", getCatalog)
}

func getCatalog(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, empty)
}
