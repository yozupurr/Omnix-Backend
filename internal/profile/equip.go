package profile

import (
	"encoding/json"
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

const maxCommandBody = 1 << 16

type equipRequest struct {
	SlotName        string `json:"slotName"`
	ItemToSlot      string `json:"itemToSlot"`
	IndexWithinSlot int    `json:"indexWithinSlot"`
}

func (h *handler) equip(w http.ResponseWriter, r *http.Request) {
	profileID := r.URL.Query().Get("profileId")
	if template, ok := templates[profileID]; !ok || !template.cosmetics {
		httpx.OperationForbidden(profileID).Write(w)
		return
	}

	var req equipRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxCommandBody)).Decode(&req); err != nil {
		httpx.InvalidRequest.Write(w)
		return
	}
	s, ok := slots[req.SlotName]
	if !ok || req.IndexWithinSlot >= max(s.size, 1) {
		httpx.InvalidRequest.Write(w)
		return
	}
	if req.ItemToSlot != "" {
		if _, owned := h.items.forSeason(version.From(r.Context()).Season())[req.ItemToSlot]; !owned {
			httpx.InvalidRequest.Write(w)
			return
		}
	}

	value, revision := h.loadouts.equip(r.PathValue("id"), s, req.ItemToSlot, req.IndexWithinSlot)
	httpx.JSON(w, http.StatusOK, newStatUpdate(profileID, s.stat, value, revision, time.Now()))
}
