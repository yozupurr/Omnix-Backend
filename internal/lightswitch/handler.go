package lightswitch

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

type launcherInfo struct {
	AppName       string `json:"appName"`
	CatalogItemID string `json:"catalogItemId"`
	Namespace     string `json:"namespace"`
}

type status struct {
	ServiceInstanceID  string       `json:"serviceInstanceId"`
	Status             string       `json:"status"`
	Message            string       `json:"message"`
	MaintenanceURI     *string      `json:"maintenanceUri"`
	OverrideCatalogIDs []string     `json:"overrideCatalogIds"`
	AllowedActions     []string     `json:"allowedActions"`
	Banned             bool         `json:"banned"`
	Launcher           launcherInfo `json:"launcherInfoDTO"`
}

var fortnite = []status{{
	ServiceInstanceID:  "fortnite",
	Status:             "UP",
	Message:            "Fortnite is online",
	OverrideCatalogIDs: []string{"a7f138b2e51945ffbfdacc1af0541053"},
	AllowedActions:     []string{"PLAY", "DOWNLOAD"},
	Launcher: launcherInfo{
		AppName:       "Fortnite",
		CatalogItemID: "4fe75bbc5a674f4f9b356b5c90567da5",
		Namespace:     "fn",
	},
}}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /lightswitch/api/service/bulk/status", bulkStatus)
}

func bulkStatus(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, fortnite)
}
