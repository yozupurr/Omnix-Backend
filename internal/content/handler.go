package content

import (
	"net/http"
	"time"

	"omnix/internal/platform/httpx"
)

const locale = "en-US"

type page struct {
	Title            string  `json:"_title"`
	ActiveDate       string  `json:"_activeDate"`
	LastModified     string  `json:"lastModified"`
	Locale           string  `json:"_locale"`
	BattleRoyaleNews section `json:"battleroyalenews"`
	EmergencyNotice  section `json:"emergencynotice"`
}

type section struct {
	Title        string `json:"_title"`
	ActiveDate   string `json:"_activeDate"`
	LastModified string `json:"lastModified"`
	Locale       string `json:"_locale"`
	News         feed   `json:"news"`
}

type feed struct {
	Type             string `json:"_type"`
	Messages         []any  `json:"messages"`
	PlatformMessages []any  `json:"platform_messages"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /content/api/pages/fortnite-game", fortniteGame)
}

func fortniteGame(w http.ResponseWriter, _ *http.Request) {
	now := httpx.Timestamp(time.Now())
	httpx.JSON(w, http.StatusOK, page{
		Title:            "Fortnite Game",
		ActiveDate:       now,
		LastModified:     now,
		Locale:           locale,
		BattleRoyaleNews: newSection("battleroyalenews", now),
		EmergencyNotice:  newSection("emergencynotice", now),
	})
}

func newSection(title, now string) section {
	return section{
		Title:        title,
		ActiveDate:   now,
		LastModified: now,
		Locale:       locale,
		News: feed{
			Type:             "Battle Royale News",
			Messages:         []any{},
			PlatformMessages: []any{},
		},
	}
}
