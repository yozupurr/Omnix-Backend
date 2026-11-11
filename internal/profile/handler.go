package profile

import (
	"encoding/json"
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
	"omnix/internal/storefront"
	"omnix/internal/version"
)

const (
	commandPath    = "POST /fortnite/api/game/v2/profile/{id}/client/"
	maxCommandBody = 1 << 16
)

var acknowledgedCommands = []string{
	"GetMcpTimeForLogin",
	"RefreshExpeditions",
	"IncrementNamedCounterStat",
	"ClientQuestLogin",
}

type handler struct {
	accounts *store
	shop     *storefront.Service
}

func Register(mux *http.ServeMux, shop *storefront.Service) {
	h := &handler{accounts: newStore(), shop: shop}
	mux.HandleFunc(commandPath+"QueryProfile", h.queryProfile)
	mux.HandleFunc(commandPath+"EquipBattleRoyaleCustomization", h.equip)
	mux.HandleFunc(commandPath+"BulkEquipBattleRoyaleCustomization", h.bulkEquip)
	mux.HandleFunc(commandPath+"PurchaseCatalogEntry", h.purchase)
	mux.HandleFunc(commandPath+"MarkItemSeen", h.markSeen)
	for _, command := range acknowledgedCommands {
		mux.HandleFunc(commandPath+command, h.acknowledge)
	}
}

func (h *handler) queryProfile(w http.ResponseWriter, r *http.Request) {
	profileID, ok := requireProfile(w, r)
	if !ok {
		return
	}
	season := version.From(r.Context()).Season()
	now := time.Now()
	h.accounts.with(r.PathValue("id"), func(a *account) {
		p := render(a, r.PathValue("id"), profileID, season, now)
		httpx.JSON(w, http.StatusOK, response{
			update: update{
				ProfileRevision:        p.Revision,
				ProfileID:              profileID,
				ProfileChanges:         []any{fullUpdate{ChangeType: "fullProfileUpdate", Profile: p}},
				ProfileCommandRevision: p.Revision,
			},
			ServerTime:      httpx.Timestamp(now),
			ResponseVersion: responseVersion,
		})
	})
}

func (h *handler) acknowledge(w http.ResponseWriter, r *http.Request) {
	profileID, ok := requireProfile(w, r)
	if !ok {
		return
	}
	h.accounts.with(r.PathValue("id"), func(a *account) {
		httpx.JSON(w, http.StatusOK, newTx(a, profileID).commit(profileID, nil, time.Now()))
	})
}

func requireProfile(w http.ResponseWriter, r *http.Request) (string, bool) {
	profileID := r.URL.Query().Get("profileId")
	if _, ok := templates[profileID]; !ok {
		httpx.OperationForbidden(profileID).Write(w)
		return "", false
	}
	return profileID, true
}

func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxCommandBody)).Decode(v); err != nil {
		httpx.InvalidRequest.Write(w)
		return false
	}
	return true
}
