package party

import (
	"encoding/json"
	"time"

	"omnix/internal/platform/httpx"
)

const (
	namespace          = "Fortnite"
	notificationPrefix = "com.epicgames.social.party.notification.v0."
	displayNameKey     = "urn:epic:member:dn_s"
)

var configDefaults = map[string]any{
	"type":        "DEFAULT",
	"sub_type":    "default",
	"joinability": "OPEN",
	"max_size":    16,
	"invite_ttl":  14400,
}

type Notifier interface {
	Send(accountID, body string)
}

type envelope struct {
	Type    string `json:"type"`
	Sent    string `json:"sent"`
	Ns      string `json:"ns"`
	PartyID string `json:"party_id"`
}

type partyUpdated struct {
	envelope
	Revision             int               `json:"revision"`
	CaptainID            string            `json:"captain_id"`
	PartyType            any               `json:"party_type"`
	PartySubType         any               `json:"party_sub_type"`
	PartyPrivacyType     any               `json:"party_privacy_type"`
	MaxNumberOfMembers   any               `json:"max_number_of_members"`
	InviteTTLSeconds     any               `json:"invite_ttl_seconds"`
	PartyStateUpdated    meta              `json:"party_state_updated"`
	PartyStateRemoved    []string          `json:"party_state_removed"`
	PartyStateOverridden map[string]string `json:"party_state_overriden"`
	CreatedAt            string            `json:"created_at"`
	UpdatedAt            string            `json:"updated_at"`
}

type memberStateUpdated struct {
	envelope
	AccountID             string            `json:"account_id"`
	AccountDN             string            `json:"account_dn"`
	Revision              int               `json:"revision"`
	MemberStateUpdated    meta              `json:"member_state_updated"`
	MemberStateRemoved    []string          `json:"member_state_removed"`
	MemberStateOverridden map[string]string `json:"member_state_overridden"`
	UpdatedAt             string            `json:"updated_at"`
}

type memberJoined struct {
	envelope
	AccountID          string     `json:"account_id"`
	AccountDN          string     `json:"account_dn"`
	Revision           int        `json:"revision"`
	Connection         connection `json:"connection"`
	MemberStateUpdated meta       `json:"member_state_updated"`
	JoinedAt           string     `json:"joined_at"`
	UpdatedAt          string     `json:"updated_at"`
}

type memberLeft struct {
	envelope
	AccountID         string `json:"account_id"`
	Revision          int    `json:"revision"`
	MemberStateUpdate meta   `json:"member_state_update"`
}

func newEnvelope(kind string, p *party, now time.Time) envelope {
	return envelope{Type: notificationPrefix + kind, Sent: httpx.Timestamp(now), Ns: namespace, PartyID: p.ID}
}

func (p *party) config(key string) any {
	if value, ok := p.Config[key]; ok && value != nil {
		return value
	}
	return configDefaults[key]
}

func (p *party) captain() string {
	for _, m := range p.Members {
		if m.Role == captainRole {
			return m.AccountID
		}
	}
	return ""
}

func (p *party) recipients(extra ...string) []string {
	ids := make([]string, 0, len(p.Members)+len(extra))
	for _, m := range p.Members {
		ids = append(ids, m.AccountID)
	}
	return append(ids, extra...)
}

func (m *member) displayName() string {
	return m.Meta[displayNameKey]
}

func partyUpdate(p *party, patch metaPatch, now time.Time) any {
	return partyUpdated{
		envelope:             newEnvelope("PARTY_UPDATED", p, now),
		Revision:             p.Revision,
		CaptainID:            p.captain(),
		PartyType:            p.config("type"),
		PartySubType:         p.config("sub_type"),
		PartyPrivacyType:     p.config("joinability"),
		MaxNumberOfMembers:   p.config("max_size"),
		InviteTTLSeconds:     p.config("invite_ttl"),
		PartyStateUpdated:    orEmpty(patch.Update),
		PartyStateRemoved:    orNone(patch.Delete),
		PartyStateOverridden: map[string]string{},
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}
}

func memberUpdate(p *party, m *member, patch metaPatch, now time.Time) any {
	return memberStateUpdated{
		envelope:              newEnvelope("MEMBER_STATE_UPDATED", p, now),
		AccountID:             m.AccountID,
		AccountDN:             m.displayName(),
		Revision:              m.Revision,
		MemberStateUpdated:    orEmpty(patch.Update),
		MemberStateRemoved:    orNone(patch.Delete),
		MemberStateOverridden: map[string]string{},
		UpdatedAt:             m.UpdatedAt,
	}
}

func memberJoin(p *party, m *member, now time.Time) any {
	return memberJoined{
		envelope:           newEnvelope("MEMBER_JOINED", p, now),
		AccountID:          m.AccountID,
		AccountDN:          m.displayName(),
		Revision:           m.Revision,
		Connection:         m.Connections[0],
		MemberStateUpdated: m.Meta,
		JoinedAt:           m.JoinedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func memberLeave(p *party, accountID string, now time.Time) any {
	return memberLeft{
		envelope:          newEnvelope("MEMBER_LEFT", p, now),
		AccountID:         accountID,
		Revision:          p.Revision,
		MemberStateUpdate: meta{},
	}
}

func (s *store) deliver(payload any, recipients []string) {
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	for _, accountID := range recipients {
		s.notifier.Send(accountID, string(body))
	}
}

func orEmpty(m meta) meta {
	if m == nil {
		return meta{}
	}
	return m
}

func orNone(list []string) []string {
	if list == nil {
		return []string{}
	}
	return list
}
