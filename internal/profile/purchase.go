package profile

import (
	"maps"
	"net/http"
	"slices"
	"strings"
	"time"

	"omnix/internal/battlepass"
	"omnix/internal/platform/httpx"
	"omnix/internal/storefront"
	"omnix/internal/version"
)

type purchaseRequest struct {
	OfferID          string `json:"offerId"`
	PurchaseQuantity int    `json:"purchaseQuantity"`
}

type lootItem struct {
	ItemType    string `json:"itemType"`
	ItemGUID    string `json:"itemGuid"`
	ItemProfile string `json:"itemProfile"`
	Quantity    int    `json:"quantity"`
}

type purchaseNotification struct {
	Type       string     `json:"type"`
	Primary    bool       `json:"primary"`
	LootResult lootResult `json:"lootResult"`
}

type lootResult struct {
	Items []lootItem `json:"items"`
}

func (h *handler) purchase(w http.ResponseWriter, r *http.Request) {
	wallet, ok := requireProfile(w, r)
	if !ok {
		return
	}
	if !templates[wallet].currency {
		httpx.OperationForbidden(wallet).Write(w)
		return
	}
	var req purchaseRequest
	if !decode(w, r, &req) {
		return
	}
	now := time.Now()
	offer, ok := h.shop.Find(req.OfferID, version.From(r.Context()).Season(), now)
	if !ok {
		httpx.CatalogOutOfDate.Write(w)
		return
	}
	quantity := 1
	if offer.Kind == storefront.KindTier {
		quantity = max(req.PurchaseQuantity, 1)
	}

	h.accounts.with(r.PathValue("id"), func(a *account) {
		total := offer.Price * quantity
		if a.vbucks < total {
			httpx.NotEnoughVBucks.Write(w)
			return
		}
		if !allowed(a, offer) {
			httpx.PurchaseNotAllowed.Write(w)
			return
		}
		t := newTx(a, wallet)
		a.vbucks -= total
		t.vbucksChanged()
		loot := t.fulfill(offer, quantity)
		switch offer.Kind {
		case storefront.KindPass, storefront.KindBundle:
			t.giftBox(passGiftBox(offer.Pass.Number), loot, now)
		case storefront.KindTier:
			t.giftBox(battlePassGiftBox, loot, now)
		}
		notification := purchaseNotification{Type: "CatalogPurchase", Primary: true, LootResult: lootResult{Items: loot}}
		httpx.JSON(w, http.StatusOK, t.commit(wallet, []any{notification}, now))
	})
}

func allowed(a *account, offer storefront.Offer) bool {
	switch offer.Kind {
	case storefront.KindItem:
		_, owned := a.athena[offer.TemplateID]
		return !owned
	case storefront.KindPass, storefront.KindBundle:
		return !a.book(offer.Pass.Number).purchased
	case storefront.KindTier:
		return a.book(offer.Pass.Number).level < offer.Pass.Tiers()
	}
	return false
}

func (t *tx) fulfill(offer storefront.Offer, quantity int) []lootItem {
	var loot []lootItem
	switch offer.Kind {
	case storefront.KindItem:
		loot = t.grant(offer.TemplateID, 1, loot)
	case storefront.KindPass, storefront.KindBundle:
		b := t.account.book(offer.Pass.Number)
		b.purchased = true
		t.stat("book_purchased", true)
		for tier := 1; tier <= b.level; tier++ {
			loot = t.grantAll(offer.Pass.Paid[tier-1], loot)
		}
		loot = t.levelUp(offer.Pass, offer.Tiers, loot)
	case storefront.KindTier:
		loot = t.levelUp(offer.Pass, offer.Tiers*quantity, loot)
	}
	return loot
}

func (t *tx) levelUp(pass battlepass.Season, tiers int, loot []lootItem) []lootItem {
	b := t.account.book(pass.Number)
	target := min(b.level+tiers, pass.Tiers())
	for tier := b.level + 1; tier <= target; tier++ {
		loot = t.grantAll(pass.Free[tier-1], loot)
		if b.purchased {
			loot = t.grantAll(pass.Paid[tier-1], loot)
		}
	}
	if target != b.level {
		b.level = target
		t.stat("book_level", target)
	}
	return loot
}

func (t *tx) grantAll(rewards battlepass.Rewards, loot []lootItem) []lootItem {
	for _, templateID := range slices.Sorted(maps.Keys(rewards)) {
		loot = t.grant(templateID, rewards[templateID], loot)
	}
	return loot
}

func (t *tx) grant(templateID string, quantity int, loot []lootItem) []lootItem {
	a := t.account
	switch {
	case strings.EqualFold(templateID, giveawayTemplate):
		a.vbucks += quantity
		t.vbucksChanged()
		return loot
	case strings.HasPrefix(templateID, bannerPrefix):
		if _, owned := a.banners[templateID]; owned {
			return loot
		}
		a.banners[templateID] = 1
		t.record(commonCoreProfile, itemAdded{ChangeType: "itemAdded", ItemID: templateID, Item: cosmetic(templateID, 1, true)})
		return append(loot, lootItem{ItemType: templateID, ItemGUID: templateID, ItemProfile: commonCoreProfile, Quantity: 1})
	case isAthenaItem(templateID):
		current, owned := a.athena[templateID]
		switch {
		case !owned:
			a.athena[templateID] = quantity
			t.record(athenaProfile, itemAdded{ChangeType: "itemAdded", ItemID: templateID, Item: cosmetic(templateID, quantity, false)})
		case strings.HasPrefix(templateID, "Token:"):
			a.athena[templateID] = current + quantity
			t.record(athenaProfile, itemQuantityChanged{ChangeType: "itemQuantityChanged", ItemID: templateID, Quantity: current + quantity})
		default:
			return loot
		}
		return append(loot, lootItem{ItemType: templateID, ItemGUID: templateID, ItemProfile: athenaProfile, Quantity: quantity})
	}
	return loot
}
