package account

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"sync"
	"time"
)

const emailDomain = "omnix.local"

type Account struct {
	ID          string
	DisplayName string
	Email       string
	LastLogin   time.Time
}

type Store struct {
	mu       sync.RWMutex
	accounts map[string]Account
}

func NewStore() *Store {
	return &Store{accounts: make(map[string]Account)}
}

func (s *Store) FindOrCreate(username string, now time.Time) Account {
	id := idFor(username)

	s.mu.Lock()
	defer s.mu.Unlock()

	account, ok := s.accounts[id]
	if !ok {
		account = Account{
			ID:          id,
			DisplayName: displayNameFor(username),
			Email:       emailFor(username),
		}
	}
	account.LastLogin = now
	s.accounts[id] = account
	return account
}

func (s *Store) Find(id string) (Account, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	account, ok := s.accounts[id]
	return account, ok
}

func idFor(username string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(username)))
	return hex.EncodeToString(sum[:16])
}

func displayNameFor(username string) string {
	if name, _, _ := strings.Cut(username, "@"); name != "" {
		return name
	}
	return username
}

func emailFor(username string) string {
	if strings.Contains(username, "@") {
		return username
	}
	return username + "@" + emailDomain
}
