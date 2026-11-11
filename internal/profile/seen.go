package profile

import (
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
)

type markSeenRequest struct {
	ItemIDs []string `json:"itemIds"`
}

func (h *handler) markSeen(w http.ResponseWriter, r *http.Request) {
	profileID, ok := requireProfile(w, r)
	if !ok {
		return
	}
	var req markSeenRequest
	if !decode(w, r, &req) {
		return
	}
	h.accounts.with(r.PathValue("id"), func(a *account) {
		t := newTx(a, profileID)
		for _, id := range req.ItemIDs {
			if _, owned := a.athena[id]; !owned {
				continue
			}
			if _, seen := a.seen[id]; seen {
				continue
			}
			a.seen[id] = struct{}{}
			t.record(athenaProfile, itemAttrChanged{ChangeType: "itemAttrChanged", ItemID: id, AttributeName: "item_seen", AttributeValue: 1})
		}
		httpx.JSON(w, http.StatusOK, t.commit(profileID, nil, time.Now()))
	})
}
