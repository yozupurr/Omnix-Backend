package calendar

import (
	"fmt"
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

const (
	farPast           = "2017-01-01T00:00:00.000Z"
	farFuture         = "9999-12-31T23:59:59.999Z"
	cacheIntervalMins = 10
)

type timeline struct {
	Channels            map[string]channel `json:"channels"`
	EventsTimeOffsetHrs int                `json:"eventsTimeOffsetHrs"`
	CacheIntervalMins   int                `json:"cacheIntervalMins"`
	CurrentTime         string             `json:"currentTime"`
}

type channel struct {
	States      []state `json:"states"`
	CacheExpire string  `json:"cacheExpire"`
}

type state struct {
	ValidFrom    string `json:"validFrom"`
	ActiveEvents []any  `json:"activeEvents"`
	State        any    `json:"state"`
}

type matchmakingState struct {
	Region map[string]any `json:"region"`
}

type eventsState struct {
	ActiveStorefronts  []any          `json:"activeStorefronts"`
	EventNamedWeights  map[string]any `json:"eventNamedWeights"`
	SeasonNumber       int            `json:"seasonNumber"`
	SeasonTemplateID   string         `json:"seasonTemplateId"`
	MatchXPBonusPoints int            `json:"matchXpBonusPoints"`
	SeasonBegin        string         `json:"seasonBegin"`
	SeasonEnd          string         `json:"seasonEnd"`
	SeasonDisplayedEnd string         `json:"seasonDisplayedEnd"`
	WeeklyStoreEnd     string         `json:"weeklyStoreEnd"`
	STWEventStoreEnd   string         `json:"stwEventStoreEnd"`
	STWWeeklyStoreEnd  string         `json:"stwWeeklyStoreEnd"`
	DailyStoreEnd      string         `json:"dailyStoreEnd"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/calendar/v1/timeline", getTimeline)
}

func getTimeline(w http.ResponseWriter, r *http.Request) {
	season := version.From(r.Context()).Season()
	httpx.JSON(w, http.StatusOK, newTimeline(season, time.Now()))
}

func newTimeline(season int, now time.Time) timeline {
	return timeline{
		Channels: map[string]channel{
			"client-matchmaking": newChannel(matchmakingState{Region: map[string]any{}}),
			"client-events": newChannel(eventsState{
				ActiveStorefronts:  []any{},
				EventNamedWeights:  map[string]any{},
				SeasonNumber:       season,
				SeasonTemplateID:   fmt.Sprintf("AthenaSeason:athenaseason%d", season),
				SeasonBegin:        farPast,
				SeasonEnd:          farFuture,
				SeasonDisplayedEnd: farFuture,
				WeeklyStoreEnd:     farFuture,
				STWEventStoreEnd:   farFuture,
				STWWeeklyStoreEnd:  farFuture,
				DailyStoreEnd:      farFuture,
			}),
		},
		CacheIntervalMins: cacheIntervalMins,
		CurrentTime:       httpx.Timestamp(now),
	}
}

func newChannel(content any) channel {
	return channel{
		States:      []state{{ValidFrom: farPast, ActiveEvents: []any{}, State: content}},
		CacheExpire: farFuture,
	}
}
