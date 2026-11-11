package storefront

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"omnix/internal/catalog"
	"omnix/internal/platform/httpx"
)

const (
	refreshIntervalHrs = 24
	dailyPurchaseHrs   = 24
	currency           = "MtxCurrency"
	noExpiration       = "9999-12-31T23:59:59.999Z"
)

var prices = map[string]map[string]int{
	"AthenaCharacter": {"uncommon": 800, "rare": 1200, "epic": 1500, "legendary": 2000},
	"AthenaPickaxe":   {"uncommon": 500, "rare": 800, "epic": 1200, "legendary": 1500},
	"AthenaGlider":    {"uncommon": 500, "rare": 800, "epic": 1200, "legendary": 1500},
	"AthenaDance":     {"uncommon": 200, "rare": 500, "epic": 800, "legendary": 800},
}

type shop struct {
	RefreshIntervalHrs int          `json:"refreshIntervalHrs"`
	DailyPurchaseHrs   int          `json:"dailyPurchaseHrs"`
	Expiration         string       `json:"expiration"`
	Storefronts        []storefront `json:"storefronts"`
}

type storefront struct {
	Name           string  `json:"name"`
	CatalogEntries []offer `json:"catalogEntries"`
}

type offer struct {
	OfferID              string        `json:"offerId"`
	DevName              string        `json:"devName"`
	OfferType            string        `json:"offerType"`
	Prices               []price       `json:"prices"`
	Categories           []string      `json:"categories"`
	DailyLimit           int           `json:"dailyLimit"`
	WeeklyLimit          int           `json:"weeklyLimit"`
	MonthlyLimit         int           `json:"monthlyLimit"`
	AppStoreID           []string      `json:"appStoreId"`
	Requirements         []requirement `json:"requirements"`
	MetaInfo             []any         `json:"metaInfo"`
	CatalogGroupPriority int           `json:"catalogGroupPriority"`
	SortPriority         int           `json:"sortPriority"`
	Title                string        `json:"title"`
	Description          string        `json:"description"`
	DisplayAssetPath     string        `json:"displayAssetPath"`
	ItemGrants           []grant       `json:"itemGrants"`
	FulfillmentIDs       []string      `json:"fulfillmentIds"`
	MatchFilter          string        `json:"matchFilter"`
	FilterWeight         float64       `json:"filterWeight"`
	Refundable           bool          `json:"refundable"`
}

type price struct {
	CurrencyType    string `json:"currencyType"`
	CurrencySubType string `json:"currencySubType"`
	RegularPrice    int    `json:"regularPrice"`
	FinalPrice      int    `json:"finalPrice"`
	BasePrice       int    `json:"basePrice"`
	SaleExpiration  string `json:"saleExpiration"`
}

type requirement struct {
	RequirementType string `json:"requirementType"`
	RequiredID      string `json:"requiredId"`
	MinQuantity     int    `json:"minQuantity"`
}

type grant struct {
	TemplateID string `json:"templateId"`
	Quantity   int    `json:"quantity"`
}

func newShop(pool []catalog.Item, season int, now time.Time) shop {
	r := currentRotation(now)
	picked := r.pick(pool, season)
	storefronts := make([]storefront, 0, len(sections))
	for _, s := range sections {
		entries := make([]offer, 0, len(picked[s.name]))
		for i, item := range picked[s.name] {
			entries = append(entries, newOffer(item, r.day, len(picked[s.name])-i))
		}
		storefronts = append(storefronts, storefront{Name: s.name, CatalogEntries: entries})
	}
	return shop{
		RefreshIntervalHrs: refreshIntervalHrs,
		DailyPurchaseHrs:   dailyPurchaseHrs,
		Expiration:         httpx.Timestamp(r.expires),
		Storefronts:        storefronts,
	}
}

func newOffer(item catalog.Item, day time.Time, priority int) offer {
	templateID := item.Type + ":" + strings.ToLower(item.ID)
	cost := priceOf(item)
	return offer{
		OfferID:        offerID(templateID, day),
		DevName:        fmt.Sprintf("[VIRTUAL]1 x %s for %d %s", item.ID, cost, currency),
		OfferType:      "StaticPrice",
		Prices:         []price{{CurrencyType: currency, RegularPrice: cost, FinalPrice: cost, BasePrice: cost, SaleExpiration: noExpiration}},
		Categories:     []string{},
		DailyLimit:     -1,
		WeeklyLimit:    -1,
		MonthlyLimit:   -1,
		AppStoreID:     []string{},
		Requirements:   []requirement{{RequirementType: "DenyOnItemOwnership", RequiredID: templateID, MinQuantity: 1}},
		MetaInfo:       []any{},
		SortPriority:   priority,
		ItemGrants:     []grant{{TemplateID: templateID, Quantity: 1}},
		FulfillmentIDs: []string{},
	}
}

func priceOf(item catalog.Item) int {
	byRarity := prices[item.Type]
	if cost, ok := byRarity[item.Rarity]; ok {
		return cost
	}
	return byRarity["epic"]
}

func offerID(templateID string, day time.Time) string {
	sum := sha256.Sum256([]byte(day.Format(time.DateOnly) + "|" + templateID))
	return "v2:/" + hex.EncodeToString(sum[:])
}
