package profile

import (
	"maps"
	"slices"
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

func newFavorites() map[string]any {
	favorites := make(map[string]any, len(slots))
	for _, s := range slots {
		if s.size > 0 {
			favorites[s.stat] = make([]string, s.size)
		} else {
			favorites[s.stat] = ""
		}
	}
	return favorites
}

func (a *account) snapshotFavorites() map[string]any {
	favorites := maps.Clone(a.favorites)
	for name, value := range favorites {
		if list, ok := value.([]string); ok {
			favorites[name] = slices.Clone(list)
		}
	}
	return favorites
}

func (a *account) equip(s slot, itemID string, index int) any {
	var value any = itemID
	if s.size > 0 {
		list := slices.Clone(a.favorites[s.stat].([]string))
		if index < 0 {
			for i := range list {
				list[i] = itemID
			}
		} else {
			list[index] = itemID
		}
		value = list
	}
	a.favorites[s.stat] = value
	return value
}
