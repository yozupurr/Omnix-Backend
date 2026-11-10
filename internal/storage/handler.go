package storage

import (
	"fmt"
	"net/http"

	"omnix/internal/platform/httpx"
)

func Register(mux *http.ServeMux, dir string) error {
	system, err := load(dir)
	if err != nil {
		return fmt.Errorf("storage: load %q: %w", dir, err)
	}
	users := newUserStore()

	mux.HandleFunc("GET /fortnite/api/cloudstorage/system", system.list)
	mux.HandleFunc("GET /fortnite/api/cloudstorage/system/{file}", system.file)
	mux.HandleFunc("GET /fortnite/api/cloudstorage/user/{id}", users.list)
	mux.HandleFunc("GET /fortnite/api/cloudstorage/user/{id}/{file}", users.get)
	mux.HandleFunc("PUT /fortnite/api/cloudstorage/user/{id}/{file}", users.put)
	return nil
}

func (c catalog) list(w http.ResponseWriter, _ *http.Request) {
	httpx.JSON(w, http.StatusOK, c.entries)
}

func (c catalog) file(w http.ResponseWriter, r *http.Request) {
	data, ok := c.content[r.PathValue("file")]
	if !ok {
		httpx.NotFound.Write(w)
		return
	}
	w.Header().Set("Content-Type", contentType)
	_, _ = w.Write(data)
}
