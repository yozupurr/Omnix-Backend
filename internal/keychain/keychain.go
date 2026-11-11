package keychain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"
	"time"
)

const source = "https://raw.githubusercontent.com/Lawin0129/LawinServer/main/responses/keychain.json"

type Keychain struct {
	keys []string
}

func Load(path string) (*Keychain, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return &Keychain{keys: []string{}}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("keychain: %w", err)
	}
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("keychain: %s: %w", path, err)
	}
	return &Keychain{keys: keys}, nil
}

func (k *Keychain) Keys() []string {
	return k.keys
}

func Fetch(ctx context.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keychain: %s returned %d", source, res.StatusCode)
	}

	var raw []string
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("keychain: %w", err)
	}
	keys := make([]string, 0, len(raw))
	for _, key := range raw {
		if strings.Count(key, ":") >= 2 {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

func Save(path string, keys []string) error {
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
