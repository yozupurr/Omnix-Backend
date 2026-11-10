package eulatracking

import "net/http"

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /eulatracking/api/shared/agreements/fn", agreement)
}

func agreement(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
