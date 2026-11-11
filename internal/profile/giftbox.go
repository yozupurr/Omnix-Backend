package profile

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
)

const (
	battlePassGiftBox           = "GiftBox:gb_battlepass"
	battlePassPurchasedGiftBox  = "GiftBox:gb_battlepasspurchased"
	firstPurchasedGiftBoxSeason = 5
)

type giftBox struct {
	templateID string
	loot       []lootItem
	giftedOn   time.Time
}

type giftBoxAttributes struct {
	FromAccountID string         `json:"fromAccountId"`
	LootList      []lootItem     `json:"lootList"`
	Params        map[string]any `json:"params"`
	Level         int            `json:"level"`
	GiftedOn      string         `json:"giftedOn"`
}

type itemRemoved struct {
	ChangeType string `json:"changeType"`
	ItemID     string `json:"itemId"`
}

type removeGiftBoxRequest struct {
	GiftBoxItemID  string   `json:"giftBoxItemId"`
	GiftBoxItemIDs []string `json:"giftBoxItemIds"`
}

func (r removeGiftBoxRequest) ids() []string {
	if r.GiftBoxItemID == "" {
		return r.GiftBoxItemIDs
	}
	return append(r.GiftBoxItemIDs, r.GiftBoxItemID)
}

func passGiftBox(season int) string {
	if season < firstPurchasedGiftBoxSeason {
		return battlePassGiftBox
	}
	return battlePassPurchasedGiftBox
}

func (g giftBox) item() item {
	return item{
		TemplateID: g.templateID,
		Attributes: giftBoxAttributes{
			LootList: g.loot,
			Params:   map[string]any{},
			Level:    1,
			GiftedOn: httpx.Timestamp(g.giftedOn),
		},
		Quantity: 1,
	}
}

func (t *tx) giftBox(templateID string, loot []lootItem, now time.Time) {
	if len(loot) == 0 {
		return
	}
	id := newGUID()
	box := giftBox{templateID: templateID, loot: loot, giftedOn: now}
	t.account.giftBoxes[id] = box
	t.record(commonCoreProfile, itemAdded{ChangeType: "itemAdded", ItemID: id, Item: box.item()})
}

func (h *handler) removeGiftBox(w http.ResponseWriter, r *http.Request) {
	if profileID := r.URL.Query().Get("profileId"); profileID != commonCoreProfile {
		httpx.OperationForbidden(profileID).Write(w)
		return
	}
	var req removeGiftBoxRequest
	if !decode(w, r, &req) {
		return
	}
	h.accounts.with(r.PathValue("id"), func(a *account) {
		t := newTx(a, commonCoreProfile)
		for _, id := range req.ids() {
			if _, ok := a.giftBoxes[id]; !ok {
				continue
			}
			delete(a.giftBoxes, id)
			t.record(commonCoreProfile, itemRemoved{ChangeType: "itemRemoved", ItemID: id})
		}
		httpx.JSON(w, http.StatusOK, t.commit(commonCoreProfile, nil, time.Now()))
	})
}

func newGUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
