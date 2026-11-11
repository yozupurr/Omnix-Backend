package profile

import (
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

type equipRequest struct {
	SlotName        string `json:"slotName"`
	ItemToSlot      string `json:"itemToSlot"`
	IndexWithinSlot int    `json:"indexWithinSlot"`
}

type bulkEquipRequest struct {
	LoadoutData []equipRequest `json:"loadoutData"`
}

func (h *handler) equip(w http.ResponseWriter, r *http.Request) {
	var req equipRequest
	if !requireAthena(w, r) || !decode(w, r, &req) {
		return
	}
	h.applyEquips(w, r, []equipRequest{req})
}

func (h *handler) bulkEquip(w http.ResponseWriter, r *http.Request) {
	var req bulkEquipRequest
	if !requireAthena(w, r) || !decode(w, r, &req) {
		return
	}
	h.applyEquips(w, r, req.LoadoutData)
}

func (h *handler) applyEquips(w http.ResponseWriter, r *http.Request, requests []equipRequest) {
	for _, req := range requests {
		s, ok := slots[req.SlotName]
		if !ok || req.IndexWithinSlot >= max(s.size, 1) {
			httpx.InvalidRequest.Write(w)
			return
		}
	}
	season := version.From(r.Context()).Season()
	h.accounts.with(r.PathValue("id"), func(a *account) {
		for _, req := range requests {
			if req.ItemToSlot != "" && !a.owns(req.ItemToSlot, season) {
				httpx.InvalidRequest.Write(w)
				return
			}
		}
		t := newTx(a, athenaProfile)
		for _, req := range requests {
			s := slots[req.SlotName]
			t.stat(s.stat, a.equip(s, req.ItemToSlot, req.IndexWithinSlot))
		}
		httpx.JSON(w, http.StatusOK, t.commit(athenaProfile, nil, time.Now()))
	})
}

func requireAthena(w http.ResponseWriter, r *http.Request) bool {
	if profileID := r.URL.Query().Get("profileId"); profileID != athenaProfile {
		httpx.OperationForbidden(profileID).Write(w)
		return false
	}
	return true
}
