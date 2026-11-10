package httpx

import (
	"encoding/json"
	"net/http"
	"time"
)

const timestampLayout = "2006-01-02T15:04:05.000Z"

func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func Timestamp(t time.Time) string {
	return t.UTC().Format(timestampLayout)
}
