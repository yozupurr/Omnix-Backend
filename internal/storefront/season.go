package storefront

import (
	"fmt"

	"omnix/internal/battlepass"
)

const (
	passPrice          = 950
	bundlePrice        = 2800
	bundleRegularPrice = 4700
	tierPrice          = 150
	bundleTiers        = 25
	displayAssets      = "/Game/Catalog/DisplayAssets/"
)

func seasonStorefront(pass battlepass.Season) storefront {
	passOffer := newStaticOffer(pass.PassOfferID, seasonDevName(pass, "BattlePass"), passPrice, 1)
	passOffer.Title = "Battle Pass"
	passOffer.ShortDescription = fmt.Sprintf("Season %d", pass.Number)
	passOffer.DisplayAssetPath = displayAsset(fmt.Sprintf("DA_BR_Season%d_BattlePass", pass.Number))

	tierOffer := newStaticOffer(pass.TierOfferID, seasonDevName(pass, "SingleTier"), tierPrice, 0)
	tierOffer.Title = "Battle Pass Level"

	entries := []offer{passOffer, tierOffer}
	if pass.BundleOfferID != "" {
		bundleOffer := newStaticOffer(pass.BundleOfferID, seasonDevName(pass, "BattleBundle"), bundlePrice, 0)
		bundleOffer.Prices[0].RegularPrice = bundleRegularPrice
		bundleOffer.Title = "Battle Bundle"
		bundleOffer.ShortDescription = fmt.Sprintf("Battle Pass + %d tiers!", bundleTiers)
		bundleOffer.DisplayAssetPath = displayAsset(fmt.Sprintf("DA_BR_Season%d_BattlePassWithLevels", pass.Number))
		entries = append(entries, bundleOffer)
	}
	return storefront{Name: fmt.Sprintf("BRSeason%d", pass.Number), CatalogEntries: entries}
}

func passOffer(pass battlepass.Season, offerID string) (Offer, bool) {
	switch offerID {
	case pass.PassOfferID:
		return Offer{Kind: KindPass, Price: passPrice, Pass: pass}, true
	case pass.TierOfferID:
		return Offer{Kind: KindTier, Price: tierPrice, Tiers: 1, Pass: pass}, true
	case pass.BundleOfferID:
		if pass.BundleOfferID != "" {
			return Offer{Kind: KindBundle, Price: bundlePrice, Tiers: bundleTiers, Pass: pass}, true
		}
	}
	return Offer{}, false
}

func seasonDevName(pass battlepass.Season, kind string) string {
	return fmt.Sprintf("BR.Season%d.%s.01", pass.Number, kind)
}

func displayAsset(name string) string {
	return displayAssets + name + "." + name
}
