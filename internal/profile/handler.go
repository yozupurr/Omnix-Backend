package profile

import (
	"net/http"
	"time"

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

func Register(mux *http.ServeMux) {
	mux.HandleFunc(commandPath+"QueryProfile", queryProfile)
	for _, command := range acknowledgedCommands {
		mux.HandleFunc(commandPath+command, acknowledge)
	}
}

func queryProfile(w http.ResponseWriter, r *http.Request) {
	profileID := r.URL.Query().Get("profileId")
	template, ok := templates[profileID]
	if !ok {
		httpx.OperationForbidden(profileID).Write(w)
		return
	}

	season := version.From(r.Context()).Season()
	now := time.Now()
	p := newProfile(r.PathValue("id"), profileID, template(season), now)
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
