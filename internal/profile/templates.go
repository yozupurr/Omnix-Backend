package profile

import "maps"

const (
	athenaProfile     = "athena"
	commonCoreProfile = "common_core"
	commonPublic      = "common_public"
	legacyProfile     = "profile0"
	vbucksTemplate    = "Currency:MtxPurchased"
	mtxPlatform       = "EpicPC"
)

type template struct {
	stats    func(a *account, season int) map[string]any
	items    func(a *account) map[string]item
	currency bool
}

var templates = map[string]template{
	athenaProfile:     {stats: athenaStats, items: athenaItems},
	commonCoreProfile: {stats: commonCoreStats, items: commonCoreItems, currency: true},
	legacyProfile:     {stats: legacyStats, items: currencyItems, currency: true},
	commonPublic:      {stats: publicStats, items: noItems},
}

func athenaStats(a *account, season int) map[string]any {
	b := a.book(season)
	stats := map[string]any{
		"season_num":     season,
		"level":          1,
		"xp":             0,
		"accountLevel":   1,
		"book_level":     b.level,
		"book_xp":        0,
		"book_purchased": b.purchased,
		"lifetime_wins":  0,
		"past_seasons":   []any{},
	}
	maps.Copy(stats, a.snapshotFavorites())
	return stats
}

func commonCoreStats(_ *account, _ int) map[string]any {
	return map[string]any{
		"current_mtx_platform": mtxPlatform,
		"mtx_affiliate":        "",
		"mtx_purchase_history": map[string]any{"refundsUsed": 0, "refundCredits": 3, "purchases": []any{}},
	}
}

func legacyStats(a *account, season int) map[string]any {
	stats := commonCoreStats(a, season)
	stats["level"] = 1
	stats["xp"] = 0
	return stats
}

func publicStats(_ *account, _ int) map[string]any {
	return map[string]any{"banner_icon": "", "banner_color": "", "homebase_name": ""}
}

func noItems(_ *account) map[string]item {
	return map[string]item{}
}

func athenaItems(a *account) map[string]item {
	items := make(map[string]item, len(a.athena))
	for id, quantity := range a.athena {
		_, seen := a.seen[id]
		items[id] = cosmetic(id, quantity, seen)
	}
	return items
}

func commonCoreItems(a *account) map[string]item {
	items := currencyItems(a)
	for id, quantity := range a.banners {
		items[id] = cosmetic(id, quantity, true)
	}
	return items
}

func currencyItems(a *account) map[string]item {
	return map[string]item{vbucksTemplate: currency(a.vbucks)}
}
