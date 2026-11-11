package storefront

import (
	"math/rand/v2"
	"time"

	"omnix/internal/catalog"
)

type section struct {
	name  string
	slots []string
}

var sections = []section{
	{name: "BRWeeklyStorefront", slots: []string{"AthenaCharacter", "AthenaCharacter"}},
	{name: "BRDailyStorefront", slots: []string{"AthenaCharacter", "AthenaCharacter", "AthenaGlider", "AthenaPickaxe", "AthenaDance", "AthenaDance"}},
}

type rotation struct {
	day     time.Time
	expires time.Time
}

func NextRotation(now time.Time) time.Time {
	return currentRotation(now).expires
}

func currentRotation(now time.Time) rotation {
	local := now.Local()
	day := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, time.Local)
	return rotation{day: day, expires: day.AddDate(0, 0, 1)}
}

func (r rotation) pick(pool []catalog.Item, season int) map[string][]catalog.Item {
	byType := make(map[string][]catalog.Item)
	for _, item := range pool {
		if _, sold := prices[item.Type]; sold && item.Rarity != "common" {
			byType[item.Type] = append(byType[item.Type], item)
		}
	}
	rng := rand.New(rand.NewPCG(uint64(r.day.Unix()), uint64(season)))
	for _, items := range byType {
		rng.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })
	}
	picked := make(map[string][]catalog.Item, len(sections))
	for _, s := range sections {
		for _, slot := range s.slots {
			candidates := byType[slot]
			if len(candidates) == 0 {
				continue
			}
			picked[s.name] = append(picked[s.name], candidates[0])
			byType[slot] = candidates[1:]
		}
	}
	return picked
}
