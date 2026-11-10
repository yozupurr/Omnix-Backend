package datarouter

import "net/http"

func Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /datarouter/api/v1/public/data", ingest)
}

func ingest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
