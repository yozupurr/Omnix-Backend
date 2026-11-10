package game

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

type worldInfo struct {
	Theaters      []any `json:"theaters"`
	Missions      []any `json:"missions"`
	MissionAlerts []any `json:"missionAlerts"`
}

var emptyWorld = worldInfo{
	Theaters:      []any{},
	Missions:      []any{},
	MissionAlerts: []any{},
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/game/v2/world/info", worldInformation)
	mux.HandleFunc("POST /fortnite/api/game/v2/tryPlayOnPlatform/account/{id}", tryPlayOnPlatform)
	mux.HandleFunc("GET /fortnite/api/game/v2/enabled_features", enabledFeatures)
	mux.HandleFunc("POST /fortnite/api/game/v2/grant_access/{id}", grantAccess)
}

func worldInformation(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, emptyWorld)
}

func tryPlayOnPlatform(w http.ResponseWriter, _ *http.Request) {
	httpx.Text(w, http.StatusOK, "true")
}

func enabledFeatures(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, []any{})
}

func grantAccess(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
