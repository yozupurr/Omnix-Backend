package account

import (
	"net/http"

	"omnix/internal/platform/httpx"
)

const (
	country  = "US"
	language = "en"
)

type response struct {
	ID                         string `json:"id"`
	DisplayName                string `json:"displayName"`
	Name                       string `json:"name"`
	Email                      string `json:"email"`
	FailedLoginAttempts        int    `json:"failedLoginAttempts"`
	LastLogin                  string `json:"lastLogin"`
	NumberOfDisplayNameChanges int    `json:"numberOfDisplayNameChanges"`
	AgeGroup                   string `json:"ageGroup"`
	Headless                   bool   `json:"headless"`
	Country                    string `json:"country"`
	LastName                   string `json:"lastName"`
	PreferredLanguage          string `json:"preferredLanguage"`
	CanUpdateDisplayName       bool   `json:"canUpdateDisplayName"`
	TFAEnabled                 bool   `json:"tfaEnabled"`
	EmailVerified              bool   `json:"emailVerified"`
	MinorVerified              bool   `json:"minorVerified"`
	MinorExpected              bool   `json:"minorExpected"`
	MinorStatus                string `json:"minorStatus"`
}

var ssoDomains = []string{"unrealengine.com", "unrealtournament.com", "fortnite.com", "epicgames.com"}

type handler struct {
	store *Store
}

func Register(mux *http.ServeMux, store *Store) {
	h := handler{store: store}
	mux.HandleFunc("GET /account/api/public/account/{id}", h.get)
	mux.HandleFunc("GET /account/api/public/account/{id}/externalAuths", h.externalAuths)
	mux.HandleFunc("GET /account/api/epicdomains/ssodomains", domains)
}

func domains(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, ssoDomains)
}

func (h handler) get(w http.ResponseWriter, r *http.Request) {
	account, ok := h.store.Find(r.PathValue("id"))
	if !ok {
		httpx.AccountNotFound.Write(w)
		return
	}
	httpx.JSON(w, http.StatusOK, newResponse(account))
}

func (h handler) externalAuths(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.store.Find(r.PathValue("id")); !ok {
		httpx.AccountNotFound.Write(w)
		return
	}
	httpx.JSON(w, http.StatusOK, []any{})
}

func newResponse(account Account) response {
	return response{
		ID:                account.ID,
		DisplayName:       account.DisplayName,
		Name:              account.DisplayName,
		Email:             account.Email,
		LastLogin:         httpx.Timestamp(account.LastLogin),
		AgeGroup:          "UNKNOWN",
		Country:           country,
		PreferredLanguage: language,
		EmailVerified:     true,
		MinorStatus:       "UNKNOWN",
	}
}
