package storage

import (
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"omnix/internal/platform/httpx"
)

const maxUserFileSize = 4 << 20

type userFile struct {
	data     []byte
	modified time.Time
}

type userStore struct {
	mu    sync.RWMutex
	files map[string]map[string]userFile
}

func newUserStore() *userStore {
	return &userStore{files: make(map[string]map[string]userFile)}
}

func (s *userStore) list(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	owned := s.files[r.PathValue("id")]
	entries := make([]entry, 0, len(owned))
	for name, file := range owned {
		entries = append(entries, newEntry(name, file.data, file.modified))
	}
	s.mu.RUnlock()

	slices.SortFunc(entries, func(a, b entry) int {
		return strings.Compare(a.Filename, b.Filename)
	})
	httpx.JSON(w, http.StatusOK, entries)
}

func (s *userStore) get(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	file, ok := s.files[r.PathValue("id")][r.PathValue("file")]
	s.mu.RUnlock()

	if !ok {
		httpx.NotFound.Write(w)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(file.data)
}

func (s *userStore) put(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxUserFileSize))
	if err != nil {
		httpx.ValidationFailed.Write(w)
		return
	}

	id, name := r.PathValue("id"), r.PathValue("file")
	s.mu.Lock()
	if s.files[id] == nil {
		s.files[id] = make(map[string]userFile)
	}
	s.files[id][name] = userFile{data: data, modified: time.Now()}
	s.mu.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
