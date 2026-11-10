package auth

import (
	"net/http"
	"strings"
	"time"

	"omnix/internal/account"
	"omnix/internal/platform/httpx"
)

type handler struct {
	accounts *account.Store
	sessions *sessions
}

func Register(mux *http.ServeMux, accounts *account.Store) func(accessToken string) (accountID, displayName string, ok bool) {
	h := handler{accounts: accounts, sessions: newSessions()}
	mux.HandleFunc("POST /account/api/oauth/token", h.token)
	mux.HandleFunc("GET /account/api/oauth/verify", h.verify)
	mux.HandleFunc("DELETE /account/api/oauth/sessions/kill", killSessions)
	return h.authenticate
}

func (h handler) authenticate(accessToken string) (string, string, bool) {
	token, ok := h.sessions.find(accessToken, time.Now())
	return token.AccountID, token.DisplayName, ok
}

func (h handler) token(w http.ResponseWriter, r *http.Request) {
	clientID, _, ok := r.BasicAuth()
	if !ok || clientID == "" {
		httpx.InvalidClient.Write(w)
		return
	}
	if r.ParseForm() != nil {
		httpx.InvalidRequest.Write(w)
		return
	}

	now := time.Now()
	switch r.PostForm.Get("grant_type") {
	case "client_credentials":
		httpx.JSON(w, http.StatusOK, newClientToken(clientID, now))
	case "password":
		username := strings.TrimSpace(r.PostForm.Get("username"))
		if username == "" {
			httpx.InvalidRequest.Write(w)
			return
		}
		user := h.accounts.FindOrCreate(username, now)
		issued := newUserToken(clientID, user, now)
		h.sessions.save(issued, now)
		httpx.JSON(w, http.StatusOK, issued)
	default:
		httpx.UnsupportedGrant.Write(w)
	}
}

func (h handler) verify(w http.ResponseWriter, r *http.Request) {
	accessToken, ok := bearerToken(r)
	if ok {
		if token, found := h.sessions.find(accessToken, time.Now()); found {
			httpx.JSON(w, http.StatusOK, token)
			return
		}
	}
	httpx.InvalidToken.Write(w)
}

func killSessions(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func bearerToken(r *http.Request) (string, bool) {
	scheme, value, ok := strings.Cut(r.Header.Get("Authorization"), " ")
	return value, ok && strings.EqualFold(scheme, "bearer")
}
