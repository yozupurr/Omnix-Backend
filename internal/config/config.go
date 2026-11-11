package config

import "os"

const (
	defaultPort       = "1337"
	defaultXMPPPort   = "5222"
	defaultStorageDir = "internal/cloudstorage"
	defaultCatalog    = "internal/cosmetics/cosmetics.json"
	defaultKeychain   = "internal/cosmetics/keychain.json"
)

type Config struct {
	Addr         string
	XMPPAddr     string
	StorageDir   string
	CatalogFile  string
	KeychainFile string
}

func Load() Config {
	return Config{
		Addr:         ":" + env("PORT", defaultPort),
		XMPPAddr:     ":" + env("XMPP_PORT", defaultXMPPPort),
		StorageDir:   env("STORAGE_DIR", defaultStorageDir),
		CatalogFile:  env("CATALOG_FILE", defaultCatalog),
		KeychainFile: env("KEYCHAIN_FILE", defaultKeychain),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
