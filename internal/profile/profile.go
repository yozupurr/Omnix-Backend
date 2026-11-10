package profile

import (
	"time"

	"omnix/internal/platform/httpx"
)

const (
	profileVersion  = "no_version"
	responseVersion = 1
)

type stats struct {
	Attributes map[string]any `json:"attributes"`
}

type profile struct {
	ID              string         `json:"_id"`
	Created         string         `json:"created"`
	Updated         string         `json:"updated"`
	Revision        int            `json:"rvn"`
	WipeNumber      int            `json:"wipeNumber"`
	AccountID       string         `json:"accountId"`
	ProfileID       string         `json:"profileId"`
	Version         string         `json:"version"`
	Items           map[string]any `json:"items"`
	Stats           stats          `json:"stats"`
	CommandRevision int            `json:"commandRevision"`
}

type change struct {
	ChangeType string  `json:"changeType"`
	Profile    profile `json:"profile"`
}

type statModified struct {
	ChangeType string `json:"changeType"`
	Name       string `json:"name"`
	Value      any    `json:"value"`
}

type response struct {
	ProfileRevision            int    `json:"profileRevision"`
	ProfileID                  string `json:"profileId"`
	ProfileChangesBaseRevision int    `json:"profileChangesBaseRevision"`
	ProfileChanges             []any  `json:"profileChanges"`
	ProfileCommandRevision     int    `json:"profileCommandRevision"`
	ServerTime                 string `json:"serverTime"`
	ResponseVersion            int    `json:"responseVersion"`
}

func newProfile(accountID, profileID string, attributes, items map[string]any, now time.Time) profile {
	if items == nil {
		items = map[string]any{}
	}
	timestamp := httpx.Timestamp(now)
	return profile{
		ID:        accountID,
		Created:   timestamp,
		Updated:   timestamp,
		AccountID: accountID,
		ProfileID: profileID,
		Version:   profileVersion,
		Items:     items,
		Stats:     stats{Attributes: attributes},
	}
}

func newAcknowledgement(profileID string, now time.Time) response {
	return response{
		ProfileID:       profileID,
		ProfileChanges:  []any{},
		ServerTime:      httpx.Timestamp(now),
		ResponseVersion: responseVersion,
	}
}

func newResponse(p profile, now time.Time) response {
	return response{
		ProfileRevision: p.Revision,
		ProfileID:       p.ProfileID,
		ProfileChanges:  []any{change{ChangeType: "fullProfileUpdate", Profile: p}},
		ServerTime:      httpx.Timestamp(now),
		ResponseVersion: responseVersion,
	}
}

func newStatUpdate(profileID, name string, value any, revision int, now time.Time) response {
	return response{
		ProfileRevision:            revision,
		ProfileID:                  profileID,
		ProfileChangesBaseRevision: revision - 1,
		ProfileChanges:             []any{statModified{ChangeType: "statModified", Name: name, Value: value}},
		ProfileCommandRevision:     revision,
		ServerTime:                 httpx.Timestamp(now),
		ResponseVersion:            responseVersion,
	}
}
