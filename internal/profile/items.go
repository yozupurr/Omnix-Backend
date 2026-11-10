package profile

import (
	"strings"
	"sync"

	"omnix/internal/catalog"
)

var noVariants = []any{}

type itemAttributes struct {
	MaxLevelBonus int   `json:"max_level_bonus"`
	Level         int   `json:"level"`
	ItemSeen      int   `json:"item_seen"`
	XP            int   `json:"xp"`
	Variants      []any `json:"variants"`
	Favorite      bool  `json:"favorite"`
}

type item struct {
	TemplateID string         `json:"templateId"`
	Attributes itemAttributes `json:"attributes"`
	Quantity   int            `json:"quantity"`
}

type itemCache struct {
	catalog *catalog.Catalog
	seasons sync.Map
}

func (c *itemCache) forSeason(season int) map[string]any {
	if cached, ok := c.seasons.Load(season); ok {
		return cached.(map[string]any)
	}
	available := c.catalog.Available(season)
	items := make(map[string]any, len(available))
	for _, cosmetic := range available {
		id := cosmetic.Type + ":" + strings.ToLower(cosmetic.ID)
		items[id] = item{
			TemplateID: id,
			Attributes: itemAttributes{Level: 1, ItemSeen: 1, Variants: noVariants},
			Quantity:   1,
		}
	}
	cached, _ := c.seasons.LoadOrStore(season, items)
	return cached.(map[string]any)
}
