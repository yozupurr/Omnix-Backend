package profile

import (
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
)

type equipRequest struct {
	SlotName        string `json:"slotName"`
	ItemToSlot      string `json:"itemToSlot"`
	IndexWithinSlot int    `json:"indexWithinSlot"`
}

func (h *handler) equip(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("profileId") != athenaProfile {
		httpx.OperationForbidden(r.URL.Query().Get("profileId")).Write(w)
		return
	}
	var req equipRequest
	if !decode(w, r, &req) {
		return
	}
	s, ok := slots[req.SlotName]
	if !ok || req.IndexWithinSlot >= max(s.size, 1) {
		httpx.InvalidRequest.Write(w)
		return
	}
	h.accounts.with(r.PathValue("id"), func(a *account) {
		if _, owned := a.athena[req.ItemToSlot]; req.ItemToSlot != "" && !owned {
			httpx.InvalidRequest.Write(w)
			return
		}
		t := newTx(a, athenaProfile)
		t.stat(s.stat, a.equip(s, req.ItemToSlot, req.IndexWithinSlot))
		httpx.JSON(w, http.StatusOK, t.commit(athenaProfile, nil, time.Now()))
	})
}
