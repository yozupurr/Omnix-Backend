package config

import "os"

const (
	defaultPort       = "1337"
	defaultXMPPPort   = "5222"
	defaultStorageDir = "internal/cloudstorage"
	defaultCatalog    = "internal/cosmetics/cosmetics.json"
)

type Config struct {
	Addr        string
	XMPPAddr    string
	StorageDir  string
	CatalogFile string
}

func Load() Config {
	return Config{
		Addr:        ":" + env("PORT", defaultPort),
		XMPPAddr:    ":" + env("XMPP_PORT", defaultXMPPPort),
		StorageDir:  env("STORAGE_DIR", defaultStorageDir),
		CatalogFile: env("CATALOG_FILE", defaultCatalog),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
