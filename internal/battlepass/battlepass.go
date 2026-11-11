package battlepass

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
)

//go:embed seasons/*.json
var files embed.FS

type Rewards map[string]int

type Season struct {
	Number        int       `json:"-"`
	PassOfferID   string    `json:"battlePassOfferId"`
	BundleOfferID string    `json:"battleBundleOfferId"`
	TierOfferID   string    `json:"tierOfferId"`
	Free          []Rewards `json:"freeRewards"`
	Paid          []Rewards `json:"paidRewards"`
}

type Catalog struct {
	seasons map[int]Season
}

func Load() (*Catalog, error) {
	names, err := fs.Glob(files, "seasons/*.json")
	if err != nil {
		return nil, err
	}
	seasons := make(map[int]Season, len(names))
	for _, name := range names {
		season, err := parse(name)
		if err != nil {
			return nil, fmt.Errorf("battlepass: %s: %w", name, err)
		}
		seasons[season.Number] = season
	}
	return &Catalog{seasons: seasons}, nil
}

func parse(name string) (Season, error) {
	var season Season
	if _, err := fmt.Sscanf(path.Base(name), "season%d.json", &season.Number); err != nil {
		return season, err
	}
	data, err := files.ReadFile(name)
	if err != nil {
		return season, err
	}
	if err := json.Unmarshal(data, &season); err != nil {
		return season, err
	}
	if season.PassOfferID == "" || season.TierOfferID == "" {
		return season, fmt.Errorf("missing offer ids")
	}
	if len(season.Free) != len(season.Paid) {
		return season, fmt.Errorf("%d free tiers for %d paid tiers", len(season.Free), len(season.Paid))
	}
	return season, nil
}

func (c *Catalog) Season(number int) (Season, bool) {
	season, ok := c.seasons[number]
	return season, ok
}

func (s Season) Tiers() int {
	return len(s.Paid)
}
