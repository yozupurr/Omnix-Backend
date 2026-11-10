package server

import (
	"log/slog"
	"net/http"
	"time"

	"omnix/internal/version"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func requestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			build := version.From(r.Context())
			log.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"cl", build.CL,
				"release", build.Release,
				"duration", time.Since(start),
				"query", r.URL.RawQuery,
				"ua", r.UserAgent(),
			)
		})
	}
}
