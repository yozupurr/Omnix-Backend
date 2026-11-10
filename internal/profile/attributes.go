package profile

var templates = map[string]func(season int) map[string]any{
	"athena":   athena,
	"profile0": profile0,
}

func athena(season int) map[string]any {
	return map[string]any{
		"season_num":     season,
		"level":          1,
		"xp":             0,
		"accountLevel":   1,
		"book_level":     1,
		"book_xp":        0,
		"book_purchased": false,
		"lifetime_wins":  0,
		"past_seasons":   []any{},
	}
}

func profile0(_ int) map[string]any {
	return map[string]any{
		"level": 1,
		"xp":    0,
	}
}
