package party

import (
	"crypto/rand"
	"encoding/hex"
	"maps"
	"strings"
	"time"

	"omnix/internal/platform/httpx"
)

const (
	captainRole = "CAPTAIN"
	memberRole  = "MEMBER"
)

type meta map[string]string

type metaPatch struct {
	Delete []string `json:"delete"`
	Update meta     `json:"update"`
}

func (m meta) apply(p metaPatch) {
	for _, key := range p.Delete {
		delete(m, key)
	}
	maps.Copy(m, p.Update)
}

type connection struct {
	ID              string `json:"id"`
	ConnectedAt     string `json:"connected_at"`
	UpdatedAt       string `json:"updated_at"`
	YieldLeadership bool   `json:"yield_leadership"`
	Meta            meta   `json:"meta"`
}

type member struct {
	AccountID   string       `json:"account_id"`
	Meta        meta         `json:"meta"`
	Connections []connection `json:"connections"`
	Revision    int          `json:"revision"`
	UpdatedAt   string       `json:"updated_at"`
	JoinedAt    string       `json:"joined_at"`
	Role        string       `json:"role"`
}

type party struct {
	ID         string         `json:"id"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	Config     map[string]any `json:"config"`
	Members    []*member      `json:"members"`
	Applicants []any          `json:"applicants"`
	Meta       meta           `json:"meta"`
	Invites    []any          `json:"invites"`
	Revision   int            `json:"revision"`
	Intentions []any          `json:"intentions"`
}

type joinInfo struct {
	Connection struct {
		ID   string `json:"id"`
		Meta meta   `json:"meta"`
	} `json:"connection"`
	Meta meta `json:"meta"`
}

func (j joinInfo) accountID() string {
	id, _, _ := strings.Cut(j.Connection.ID, "@")
	return id
}

func newParty(config map[string]any, partyMeta meta, now time.Time) *party {
	stamp := httpx.Timestamp(now)
	if config == nil {
		config = map[string]any{}
	}
	if partyMeta == nil {
		partyMeta = meta{}
	}
	return &party{
		ID:         newID(),
		CreatedAt:  stamp,
		UpdatedAt:  stamp,
		Config:     config,
		Members:    []*member{},
		Applicants: []any{},
		Meta:       partyMeta,
		Invites:    []any{},
		Intentions: []any{},
	}
}

func (p *party) join(accountID string, info joinInfo, now time.Time) {
	stamp := httpx.Timestamp(now)
	role := memberRole
	if len(p.Members) == 0 {
		role = captainRole
	}
	memberMeta := info.Meta
	if memberMeta == nil {
		memberMeta = meta{}
	}
	connectionMeta := info.Connection.Meta
	if connectionMeta == nil {
		connectionMeta = meta{}
	}
	p.leave(accountID)
	p.Members = append(p.Members, &member{
		AccountID: accountID,
		Meta:      memberMeta,
		Connections: []connection{{
			ID:          info.Connection.ID,
			ConnectedAt: stamp,
			UpdatedAt:   stamp,
			Meta:        connectionMeta,
		}},
		UpdatedAt: stamp,
		JoinedAt:  stamp,
		Role:      role,
	})
	p.touch(now)
}

func (p *party) leave(accountID string) bool {
	for i, m := range p.Members {
		if m.AccountID != accountID {
			continue
		}
		p.Members = append(p.Members[:i], p.Members[i+1:]...)
		if m.Role == captainRole && len(p.Members) > 0 {
			p.Members[0].Role = captainRole
		}
		return true
	}
	return false
}

func (p *party) member(accountID string) (*member, bool) {
	for _, m := range p.Members {
		if m.AccountID == accountID {
			return m, true
		}
	}
	return nil, false
}

func (p *party) promote(accountID string) bool {
	target, ok := p.member(accountID)
	if !ok {
		return false
	}
	for _, m := range p.Members {
		m.Role = memberRole
	}
	target.Role = captainRole
	return true
}

func (p *party) touch(now time.Time) {
	p.Revision++
	p.UpdatedAt = httpx.Timestamp(now)
}

func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
