package profile

import (
	"maps"
	"net/http"
	"time"

	"omnix/internal/catalog"
	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

const commandPath = "POST /fortnite/api/game/v2/profile/{id}/client/"

var acknowledgedCommands = []string{
	"GetMcpTimeForLogin",
	"RefreshExpeditions",
	"IncrementNamedCounterStat",
	"ClientQuestLogin",
}

type handler struct {
	items    *itemCache
	loadouts *loadouts
}

func Register(mux *http.ServeMux, cosmetics *catalog.Catalog) {
	h := &handler{items: &itemCache{catalog: cosmetics}, loadouts: newLoadouts()}
	mux.HandleFunc(commandPath+"QueryProfile", h.queryProfile)
	mux.HandleFunc(commandPath+"EquipBattleRoyaleCustomization", h.equip)
	for _, command := range acknowledgedCommands {
		mux.HandleFunc(commandPath+command, acknowledge)
	}
}

func (h *handler) queryProfile(w http.ResponseWriter, r *http.Request) {
	profileID := r.URL.Query().Get("profileId")
	template, ok := templates[profileID]
	if !ok {
		httpx.OperationForbidden(profileID).Write(w)
		return
	}

	season := version.From(r.Context()).Season()
	now := time.Now()
	attributes := template.attributes(season)
	var items map[string]any
	revision := 0
	if template.cosmetics {
		items = h.items.forSeason(season)
		var favorites map[string]any
		revision, favorites = h.loadouts.snapshot(r.PathValue("id"))
		maps.Copy(attributes, favorites)
	}
	p := newProfile(r.PathValue("id"), profileID, attributes, items, now)
	p.Revision = revision
	httpx.JSON(w, http.StatusOK, newResponse(p, now))
}

func acknowledge(w http.ResponseWriter, r *http.Request) {
	profileID := r.URL.Query().Get("profileId")
	if _, ok := templates[profileID]; !ok {
		httpx.OperationForbidden(profileID).Write(w)
		return
	}
	httpx.JSON(w, http.StatusOK, newAcknowledgement(profileID, time.Now()))
}
