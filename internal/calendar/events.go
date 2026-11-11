package calendar

import (
	"fmt"

	"omnix/internal/version"
)

const eventFlag = "EventFlag."

type activeEvent struct {
	EventType   string `json:"eventType"`
	ActiveUntil string `json:"activeUntil"`
	ActiveSince string `json:"activeSince"`
}

type phase struct {
	major, minor int
	flags        []string
}

var seasonalEvents = map[int][]phase{
	3: {
		{major: 3, minor: 0, flags: []string{"Spring2018Phase1"}},
		{major: 3, minor: 1, flags: []string{"Spring2018Phase2"}},
		{major: 3, minor: 3, flags: []string{"Spring2018Phase3"}},
		{major: 3, minor: 4, flags: []string{"Spring2018Phase4"}},
	},
	4: {
		{major: 4, minor: 0, flags: []string{"Blockbuster2018", "Blockbuster2018Phase1"}},
		{major: 4, minor: 3, flags: []string{"Blockbuster2018Phase2"}},
		{major: 4, minor: 4, flags: []string{"Blockbuster2018Phase3"}},
		{major: 4, minor: 5, flags: []string{"Blockbuster2018Phase4"}},
	},
	5: {
		{major: 5, minor: 0, flags: []string{"RoadTrip2018", "Horde", "Anniversary2018_BR", "LTM_Heist"}},
	},
}

func activeEvents(build version.Build) []activeEvent {
	season := build.Season()
	flags := []string{fmt.Sprintf("Season%d", season), fmt.Sprintf("LobbySeason%d", season)}
	for _, p := range seasonalEvents[season] {
		if build.AtLeast(p.major, p.minor) {
			flags = append(flags, p.flags...)
		}
	}
	events := make([]activeEvent, len(flags))
	for i, flag := range flags {
		events[i] = activeEvent{EventType: eventFlag + flag, ActiveUntil: farFuture, ActiveSince: farPast}
	}
	return events
}
