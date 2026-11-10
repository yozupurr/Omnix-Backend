package config

import "os"

const (
	defaultPort       = "1337"
	defaultStorageDir = "internal/cloudstorage"
)

type Config struct {
	Addr       string
	StorageDir string
}

func Load() Config {
	return Config{
		Addr:       ":" + env("PORT", defaultPort),
		StorageDir: env("STORAGE_DIR", defaultStorageDir),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
