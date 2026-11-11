package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const endpoint = "https://fortnite-api.com/v2/cosmetics/br"

var cosmeticTypes = map[string]struct{}{
	"AthenaCharacter":       {},
	"AthenaBackpack":        {},
	"AthenaPickaxe":         {},
	"AthenaGlider":          {},
	"AthenaSkyDiveContrail": {},
	"AthenaDance":           {},
	"AthenaItemWrap":        {},
	"AthenaLoadingScreen":   {},
	"AthenaMusicPack":       {},
}

type apiResponse struct {
	Data []struct {
		ID   string `json:"id"`
		Type struct {
			BackendValue string `json:"backendValue"`
		} `json:"type"`
		Rarity struct {
			Value string `json:"value"`
		} `json:"rarity"`
		Introduction *struct {
			BackendValue int `json:"backendValue"`
		} `json:"introduction"`
	} `json:"data"`
}

func Fetch(ctx context.Context) ([]Item, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog: %s returned %d", endpoint, res.StatusCode)
	}

	var body apiResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}

	items := make([]Item, 0, len(body.Data))
	for _, c := range body.Data {
		if _, ok := cosmeticTypes[c.Type.BackendValue]; !ok {
			continue
		}
		if c.Introduction == nil || c.Introduction.BackendValue <= 0 {
			continue
		}
		items = append(items, Item{
			ID:     c.ID,
			Type:   c.Type.BackendValue,
			Season: c.Introduction.BackendValue,
			Rarity: c.Rarity.Value,
		})
	}
	return New(items).items, nil
}
