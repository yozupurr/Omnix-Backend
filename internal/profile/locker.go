package profile

import (
	"net/http"
	"slices"
	"time"

	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

const (
	cosmeticLockerSeason = 12
	lockerItemID         = "loadout_1"
	lockerTemplate       = "CosmeticLocker:cosmeticlocker_athena"
	lockerSlotsAttribute = "locker_slots_data"
	bannerIconAttribute  = "banner_icon_template"
	bannerColorAttribute = "banner_color_template"
)

type lockerSlot struct {
	Items          []string `json:"items"`
	ActiveVariants []any    `json:"activeVariants"`
}

type lockerSlots struct {
	Slots map[string]lockerSlot `json:"slots"`
}

type lockerAttributes struct {
	LockerSlotsData     lockerSlots `json:"locker_slots_data"`
	UseCount            int         `json:"use_count"`
	BannerIconTemplate  string      `json:"banner_icon_template"`
	BannerColorTemplate string      `json:"banner_color_template"`
	LockerName          string      `json:"locker_name"`
	ItemSeen            bool        `json:"item_seen"`
	Favorite            bool        `json:"favorite"`
}

type setLockerSlotRequest struct {
	LockerItem string `json:"lockerItem"`
	Category   string `json:"category"`
	ItemToSlot string `json:"itemToSlot"`
	SlotIndex  int    `json:"slotIndex"`
}

type setLockerBannerRequest struct {
	LockerItem              string `json:"lockerItem"`
	BannerIconTemplateName  string `json:"bannerIconTemplateName"`
	BannerColorTemplateName string `json:"bannerColorTemplateName"`
}

func usesCosmeticLocker(season int) bool {
	return season >= cosmeticLockerSeason
}

func lockerStats() map[string]any {
	return map[string]any{
		"loadouts":             []string{lockerItemID},
		"active_loadout_index": 0,
		"last_applied_loadout": lockerItemID,
	}
}

func (a *account) lockerSlots() lockerSlots {
	result := lockerSlots{Slots: make(map[string]lockerSlot, len(slots))}
	for category, s := range slots {
		var items []string
		switch value := a.favorites[s.stat].(type) {
		case []string:
			items = slices.Clone(value)
		case string:
			items = []string{value}
		}
		result.Slots[category] = lockerSlot{
			Items:          items,
			ActiveVariants: make([]any, len(items)),
		}
	}
	return result
}

func (a *account) locker() item {
	return item{
		TemplateID: lockerTemplate,
		Attributes: lockerAttributes{
			LockerSlotsData:     a.lockerSlots(),
			BannerIconTemplate:  a.bannerIcon,
			BannerColorTemplate: a.bannerColor,
			ItemSeen:            true,
		},
		Quantity: 1,
	}
}

func (h *handler) setLockerSlot(w http.ResponseWriter, r *http.Request) {
	var req setLockerSlotRequest
	if !requireAthena(w, r) || !decode(w, r, &req) {
		return
	}
	s, ok := slots[req.Category]
	if !ok || req.LockerItem != lockerItemID || req.SlotIndex >= max(s.size, 1) {
		httpx.InvalidRequest.Write(w)
		return
	}
	season := version.From(r.Context()).Season()
	h.accounts.with(r.PathValue("id"), func(a *account) {
		if req.ItemToSlot != "" && !a.owns(req.ItemToSlot, season) {
			httpx.InvalidRequest.Write(w)
			return
		}
		t := newTx(a, athenaProfile)
		t.stat(s.stat, a.equip(s, req.ItemToSlot, req.SlotIndex))
		t.record(athenaProfile, itemAttrChanged{
			ChangeType:     "itemAttrChanged",
			ItemID:         lockerItemID,
			AttributeName:  lockerSlotsAttribute,
			AttributeValue: a.lockerSlots(),
		})
		httpx.JSON(w, http.StatusOK, t.commit(athenaProfile, nil, time.Now()))
	})
}

func (h *handler) setLockerBanner(w http.ResponseWriter, r *http.Request) {
	var req setLockerBannerRequest
	if !requireAthena(w, r) || !decode(w, r, &req) {
		return
	}
	if req.LockerItem != lockerItemID {
		httpx.InvalidRequest.Write(w)
		return
	}
	h.accounts.with(r.PathValue("id"), func(a *account) {
		a.bannerIcon, a.bannerColor = req.BannerIconTemplateName, req.BannerColorTemplateName
		t := newTx(a, athenaProfile)
		t.record(athenaProfile, itemAttrChanged{ChangeType: "itemAttrChanged", ItemID: lockerItemID, AttributeName: bannerIconAttribute, AttributeValue: a.bannerIcon})
		t.record(athenaProfile, itemAttrChanged{ChangeType: "itemAttrChanged", ItemID: lockerItemID, AttributeName: bannerColorAttribute, AttributeValue: a.bannerColor})
		httpx.JSON(w, http.StatusOK, t.commit(athenaProfile, nil, time.Now()))
	})
}
