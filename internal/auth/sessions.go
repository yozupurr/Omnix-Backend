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
	mu     sync.RWMutex
	tokens map[string]session
}

func newSessions() *sessions {
	return &sessions{tokens: make(map[string]session)}
}

func (s *sessions) save(token userToken, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for key, active := range s.tokens {
		if now.After(active.expires) {
			delete(s.tokens, key)
		}
	}
	s.tokens[token.AccessToken] = session{token: token, expires: now.Add(accessLifetime)}
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
