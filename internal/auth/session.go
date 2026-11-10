package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"omnix/internal/account"
	"omnix/internal/platform/httpx"
)

const (
	service         = "fortnite"
	accessLifetime  = 8 * time.Hour
	refreshLifetime = 32 * time.Hour
)

type clientToken struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	ExpiresAt      string `json:"expires_at"`
	TokenType      string `json:"token_type"`
	ClientID       string `json:"client_id"`
	InternalClient bool   `json:"internal_client"`
	ClientService  string `json:"client_service"`
}

type userToken struct {
	clientToken
	RefreshToken     string `json:"refresh_token"`
	RefreshExpires   int    `json:"refresh_expires"`
	RefreshExpiresAt string `json:"refresh_expires_at"`
	AccountID        string `json:"account_id"`
	App              string `json:"app"`
	InAppID          string `json:"in_app_id"`
	DeviceID         string `json:"device_id"`
	DisplayName      string `json:"displayName"`
}

func newClientToken(clientID string, now time.Time) clientToken {
	return clientToken{
		AccessToken:    randomHex(16),
		ExpiresIn:      int(accessLifetime.Seconds()),
		ExpiresAt:      httpx.Timestamp(now.Add(accessLifetime)),
		TokenType:      "bearer",
		ClientID:       clientID,
		InternalClient: true,
		ClientService:  service,
	}
}

func newUserToken(clientID string, user account.Account, now time.Time) userToken {
	return userToken{
		clientToken:      newClientToken(clientID, now),
		RefreshToken:     randomHex(16),
		RefreshExpires:   int(refreshLifetime.Seconds()),
		RefreshExpiresAt: httpx.Timestamp(now.Add(refreshLifetime)),
		AccountID:        user.ID,
		App:              service,
		InAppID:          user.ID,
		DeviceID:         randomHex(16),
		DisplayName:      user.DisplayName,
	}
}

func randomHex(size int) string {
	buf := make([]byte, size)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
