package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"runtime/debug"
	"time"
)

type contextKey string

const correlationIDKey contextKey = "correlation_id"

func correlationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("correlationId")
		if id == "" {
			id = newCorrelationID()
		}

		w.Header().Set("correlationId", id)
		ctx := context.WithValue(r.Context(), correlationIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", time.Since(start)),
				slog.String("correlation_id", correlationIDFromContext(r.Context())),
			)
		})
	}
}

func recoverer(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered",
						slog.Any("panic", recovered),
						slog.String("stack", string(debug.Stack())),
						slog.String("correlation_id", correlationIDFromContext(r.Context())),
					)
					writeProblem(w, r, http.StatusInternalServerError, "internal-error", "An unexpected error occurred.")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func correlationIDFromContext(ctx context.Context) string {
	value, _ := ctx.Value(correlationIDKey).(string)
	return value
}

func newCorrelationID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return time.Now().UTC().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(bytes[:])
}
