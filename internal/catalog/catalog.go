package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"
)

type Item struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Season int    `json:"season"`
	Rarity string `json:"rarity,omitempty"`
}

type Catalog struct {
	items []Item
}

func New(items []Item) *Catalog {
	sorted := slices.Clone(items)
	slices.SortStableFunc(sorted, func(a, b Item) int { return a.Season - b.Season })
	return &Catalog{items: sorted}
}

func Load(path string) (*Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}
	var items []Item
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("catalog: %s: %w", path, err)
	}
	return New(items), nil
}

func Save(path string, items []Item) error {
	var b strings.Builder
	b.WriteString("[\n")
	for i, item := range items {
		line, err := json.Marshal(item)
		if err != nil {
			return err
		}
		b.Write(line)
		if i < len(items)-1 {
			b.WriteByte(',')
		}
		b.WriteByte('\n')
	}
	b.WriteString("]\n")
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func (c *Catalog) Len() int {
	return len(c.items)
}

func (c *Catalog) Available(season int) []Item {
	n := sort.Search(len(c.items), func(i int) bool { return c.items[i].Season > season })
	return c.items[:n]
}
