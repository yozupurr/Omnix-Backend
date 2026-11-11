package auth

import (
	"sync"
	"time"
)

type session struct {
	token   userToken
	expires time.Time
}

type sessions struct {
	mu        sync.RWMutex
	tokens    map[string]session
	refreshes map[string]session
}

func newSessions() *sessions {
	return &sessions{tokens: make(map[string]session), refreshes: make(map[string]session)}
}

func (s *sessions) save(token userToken, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, active := range s.tokens {
		if now.After(active.expires) {
			delete(s.tokens, key)
		}
	}
	for key, active := range s.refreshes {
		if now.After(active.expires) {
			delete(s.refreshes, key)
		}
	}
	s.tokens[token.AccessToken] = session{token: token, expires: now.Add(accessLifetime)}
	s.refreshes[token.RefreshToken] = session{token: token, expires: now.Add(refreshLifetime)}
}

func (s *sessions) redeem(refreshToken string, now time.Time) (userToken, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	active, ok := s.refreshes[refreshToken]
	if !ok || now.After(active.expires) {
		return userToken{}, false
	}
	delete(s.refreshes, refreshToken)
	return active.token, true
}

func (s *sessions) find(accessToken string, now time.Time) (userToken, bool) {
	s.mu.RLock()
	active, ok := s.tokens[accessToken]
	s.mu.RUnlock()

	if !ok || now.After(active.expires) {
		return userToken{}, false
	}
	token := active.token
	token.ExpiresIn = int(active.expires.Sub(now).Seconds())
	return token, true
}

func (s *sessions) revoke(accessToken string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if active, ok := s.tokens[accessToken]; ok {
		delete(s.refreshes, active.token.RefreshToken)
		delete(s.tokens, accessToken)
	}
}
