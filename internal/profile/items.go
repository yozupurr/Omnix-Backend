package profile

import "strings"

const (
	sharedPlatform   = "Shared"
	giveawayTemplate = "currency:mtxgiveaway"
	bannerPrefix     = "HomebaseBannerIcon:"
)

type item struct {
	TemplateID string `json:"templateId"`
	Attributes any    `json:"attributes"`
	Quantity   int    `json:"quantity"`
}

type cosmeticAttributes struct {
	MaxLevelBonus int   `json:"max_level_bonus"`
	Level         int   `json:"level"`
	ItemSeen      int   `json:"item_seen"`
	XP            int   `json:"xp"`
	Variants      []any `json:"variants"`
	Favorite      bool  `json:"favorite"`
}

type currencyAttributes struct {
	Platform string `json:"platform"`
}

func cosmetic(templateID string, quantity int, seen bool) item {
	attributes := cosmeticAttributes{Level: 1, Variants: []any{}}
	if seen {
		attributes.ItemSeen = 1
	}
	return item{TemplateID: templateID, Attributes: attributes, Quantity: quantity}
}

func currency(amount int) item {
	return item{TemplateID: vbucksTemplate, Attributes: currencyAttributes{Platform: sharedPlatform}, Quantity: amount}
}

func isAthenaItem(templateID string) bool {
	return strings.HasPrefix(templateID, "Athena") || strings.HasPrefix(templateID, "Token:")
}
