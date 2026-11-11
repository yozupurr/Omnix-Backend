package storefront

import (
	"net/http"
	"time"

	"omnix/internal/battlepass"
	"omnix/internal/catalog"
	"omnix/internal/keychain"
	"omnix/internal/platform/httpx"
	"omnix/internal/version"
)

type Kind int

const (
	KindItem Kind = iota
	KindPass
	KindBundle
	KindTier
)

type Offer struct {
	Kind       Kind
	Price      int
	TemplateID string
	Tiers      int
	Pass       battlepass.Season
}

type Service struct {
	cosmetics *catalog.Catalog
	passes    *battlepass.Catalog
	keys      *keychain.Keychain
}

func New(cosmetics *catalog.Catalog, passes *battlepass.Catalog, keys *keychain.Keychain) *Service {
	return &Service{cosmetics: cosmetics, passes: passes, keys: keys}
}

func (s *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /fortnite/api/storefront/v2/catalog", s.catalog)
	mux.HandleFunc("GET /fortnite/api/storefront/v2/keychain", s.keychain)
}

func (s *Service) keychain(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, s.keys.Keys())
}

func (s *Service) catalog(w http.ResponseWriter, r *http.Request) {
	season := version.From(r.Context()).Season()
	var pass *battlepass.Season
	if p, ok := s.passes.Season(season); ok {
		pass = &p
	}
	httpx.JSON(w, http.StatusOK, newShop(s.pool(season), season, pass, time.Now()))
}

func (s *Service) Find(offerID string, season int, now time.Time) (Offer, bool) {
	if pass, ok := s.passes.Season(season); ok {
		if offer, ok := passOffer(pass, offerID); ok {
			return offer, true
		}
	}
	r := currentRotation(now)
	for _, items := range r.pick(s.pool(season), season) {
		for _, item := range items {
			templateID := templateIDOf(item)
			if offerID == offerIDFor(templateID, r.day) {
				return Offer{Kind: KindItem, Price: priceOf(item), TemplateID: templateID}, true
			}
		}
	}
	return Offer{}, false
}

func (s *Service) pool(season int) []catalog.Item {
	return s.cosmetics.Available(max(season-1, 1))
}
