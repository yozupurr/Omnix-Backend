package profile

import (
	"time"

	"omnix/internal/platform/httpx"
)

const (
	profileVersion  = "no_version"
	responseVersion = 1
)

type profile struct {
	ID              string          `json:"_id"`
	Created         string          `json:"created"`
	Updated         string          `json:"updated"`
	Revision        int             `json:"rvn"`
	WipeNumber      int             `json:"wipeNumber"`
	AccountID       string          `json:"accountId"`
	ProfileID       string          `json:"profileId"`
	Version         string          `json:"version"`
	Items           map[string]item `json:"items"`
	Stats           stats           `json:"stats"`
	CommandRevision int             `json:"commandRevision"`
}

type stats struct {
	Attributes map[string]any `json:"attributes"`
}

type fullUpdate struct {
	ChangeType string  `json:"changeType"`
	Profile    profile `json:"profile"`
}

type statModified struct {
	ChangeType string `json:"changeType"`
	Name       string `json:"name"`
	Value      any    `json:"value"`
}

type itemAdded struct {
	ChangeType string `json:"changeType"`
	ItemID     string `json:"itemId"`
	Item       item   `json:"item"`
}

type itemQuantityChanged struct {
	ChangeType string `json:"changeType"`
	ItemID     string `json:"itemId"`
	Quantity   int    `json:"quantity"`
}

type itemAttrChanged struct {
	ChangeType     string `json:"changeType"`
	ItemID         string `json:"itemId"`
	AttributeName  string `json:"attributeName"`
	AttributeValue any    `json:"attributeValue"`
}

type update struct {
	ProfileRevision            int    `json:"profileRevision"`
	ProfileID                  string `json:"profileId"`
	ProfileChangesBaseRevision int    `json:"profileChangesBaseRevision"`
	ProfileChanges             []any  `json:"profileChanges"`
	ProfileCommandRevision     int    `json:"profileCommandRevision"`
}

type response struct {
	update
	Notifications   []any    `json:"notifications,omitempty"`
	MultiUpdate     []update `json:"multiUpdate,omitempty"`
	ServerTime      string   `json:"serverTime"`
	ResponseVersion int      `json:"responseVersion"`
}

func render(a *account, accountID, profileID string, season int, now time.Time) profile {
	t := templates[profileID]
	timestamp := httpx.Timestamp(now)
	revision := a.revisions[profileID]
	return profile{
		ID:              accountID,
		Created:         timestamp,
		Updated:         timestamp,
		Revision:        revision,
		AccountID:       accountID,
		ProfileID:       profileID,
		Version:         profileVersion,
		Items:           t.items(a),
		Stats:           stats{Attributes: t.stats(a, season)},
		CommandRevision: revision,
	}
}

type tx struct {
	account *account
	wallet  string
	changes map[string][]any
	order   []string
}

func newTx(a *account, wallet string) *tx {
	return &tx{account: a, wallet: wallet, changes: map[string][]any{}}
}

func (t *tx) record(profileID string, change any) {
	if _, seen := t.changes[profileID]; !seen {
		t.order = append(t.order, profileID)
	}
	t.changes[profileID] = append(t.changes[profileID], change)
}

func (t *tx) stat(name string, value any) {
	t.record(athenaProfile, statModified{ChangeType: "statModified", Name: name, Value: value})
}

func (t *tx) vbucksChanged() {
	t.record(t.wallet, itemQuantityChanged{ChangeType: "itemQuantityChanged", ItemID: vbucksTemplate, Quantity: t.account.vbucks})
}

func (t *tx) commit(primary string, notifications []any, now time.Time) response {
	res := response{
		update:          t.bump(primary),
		Notifications:   notifications,
		ServerTime:      httpx.Timestamp(now),
		ResponseVersion: responseVersion,
	}
	for _, profileID := range t.order {
		if profileID != primary {
			res.MultiUpdate = append(res.MultiUpdate, t.bump(profileID))
		}
	}
	return res
}

func (t *tx) bump(profileID string) update {
	changes := t.changes[profileID]
	base := t.account.revisions[profileID]
	if len(changes) > 0 {
		t.account.revisions[profileID]++
	} else {
		changes = []any{}
	}
	revision := t.account.revisions[profileID]
	return update{
		ProfileRevision:            revision,
		ProfileID:                  profileID,
		ProfileChangesBaseRevision: base,
		ProfileChanges:             changes,
		ProfileCommandRevision:     revision,
	}
}
