package profile

import (
	"maps"
	"slices"
	"sync"
)

type slot struct {
	stat string
	size int
}

var slots = map[string]slot{
	"Character":       {stat: "favorite_character"},
	"Backpack":        {stat: "favorite_backpack"},
	"Pickaxe":         {stat: "favorite_pickaxe"},
	"Glider":          {stat: "favorite_glider"},
	"SkyDiveContrail": {stat: "favorite_skydivecontrail"},
	"LoadingScreen":   {stat: "favorite_loadingscreen"},
	"MusicPack":       {stat: "favorite_musicpack"},
	"Dance":           {stat: "favorite_dance", size: 6},
	"ItemWrap":        {stat: "favorite_itemwraps", size: 7},
}

type loadout struct {
	revision int
	stats    map[string]any
}

type loadouts struct {
	mu        sync.Mutex
	byAccount map[string]*loadout
}

func newLoadouts() *loadouts {
	return &loadouts{byAccount: map[string]*loadout{}}
}

func newLoadout() *loadout {
	stats := make(map[string]any, len(slots))
	for _, s := range slots {
		if s.size > 0 {
			stats[s.stat] = make([]string, s.size)
		} else {
			stats[s.stat] = ""
		}
	}
	return &loadout{stats: stats}
}

func (l *loadouts) get(accountID string) *loadout {
	current, ok := l.byAccount[accountID]
	if !ok {
		current = newLoadout()
		l.byAccount[accountID] = current
	}
	return current
}

func (l *loadouts) snapshot(accountID string) (int, map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.get(accountID)
	stats := maps.Clone(current.stats)
	for name, value := range stats {
		if list, ok := value.([]string); ok {
			stats[name] = slices.Clone(list)
		}
	}
	return current.revision, stats
}

func (l *loadouts) equip(accountID string, s slot, itemID string, index int) (any, int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.get(accountID)
	var value any = itemID
	if s.size > 0 {
		list := slices.Clone(current.stats[s.stat].([]string))
		if index < 0 {
			for i := range list {
				list[i] = itemID
			}
		} else {
			list[index] = itemID
		}
		value = list
	}
	current.stats[s.stat] = value
	current.revision++
	return value, current.revision
}
