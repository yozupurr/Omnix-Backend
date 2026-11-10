package storage

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"omnix/internal/platform/httpx"
)

const (
	extension   = ".ini"
	contentType = "application/octet-stream"
)

type entry struct {
	UniqueFilename string            `json:"uniqueFilename"`
	Filename       string            `json:"filename"`
	Hash           string            `json:"hash"`
	Hash256        string            `json:"hash256"`
	Length         int               `json:"length"`
	ContentType    string            `json:"contentType"`
	Uploaded       string            `json:"uploaded"`
	StorageType    string            `json:"storageType"`
	StorageIDs     map[string]string `json:"storageIds"`
	DoNotCache     bool              `json:"doNotCache"`
}

type catalog struct {
	entries []entry
	content map[string][]byte
}

func load(dir string) (catalog, error) {
	items, err := os.ReadDir(dir)
	if err != nil {
		return catalog{}, err
	}

	c := catalog{
		entries: make([]entry, 0, len(items)),
		content: make(map[string][]byte, len(items)),
	}
	for _, item := range items {
		if item.IsDir() || filepath.Ext(item.Name()) != extension {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, item.Name()))
		if err != nil {
			return catalog{}, err
		}
		info, err := item.Info()
		if err != nil {
			return catalog{}, err
		}
		c.entries = append(c.entries, newEntry(item.Name(), data, info.ModTime()))
		c.content[item.Name()] = data
	}
	return c, nil
}

func newEntry(name string, data []byte, modified time.Time) entry {
	sum1 := sha1.Sum(data)
	sum256 := sha256.Sum256(data)
	return entry{
		UniqueFilename: name,
		Filename:       name,
		Hash:           hex.EncodeToString(sum1[:]),
		Hash256:        hex.EncodeToString(sum256[:]),
		Length:         len(data),
		ContentType:    contentType,
		Uploaded:       httpx.Timestamp(modified),
		StorageType:    "S3",
		StorageIDs:     map[string]string{},
		DoNotCache:     true,
	}
}
