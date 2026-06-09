package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"streaming-golang/internal/app/apperr"
)

type problemDetails struct {
	Type          string            `json:"type"`
	Title         string            `json:"title"`
	Status        int               `json:"status"`
	Detail        string            `json:"detail"`
	Instance      string            `json:"instance"`
	CorrelationID string            `json:"correlationId,omitempty"`
	Errors        map[string]string `json:"errors,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeProblem(w http.ResponseWriter, r *http.Request, status int, title, detail string) {
	w.Header().Set("Content-Type", "application/problem+json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(problemDetails{
		Type:          "about:blank",
		Title:         title,
		Status:        status,
		Detail:        detail,
		Instance:      r.URL.Path,
		CorrelationID: correlationIDFromContext(r.Context()),
	})
}

func writeAppError(w http.ResponseWriter, r *http.Request, err error) {
	var appError *apperr.Error
	if !errors.As(err, &appError) {
		writeProblem(w, r, http.StatusInternalServerError, "internal-error", "An unexpected error occurred.")
		return
	}

	status := http.StatusInternalServerError
	switch appError.Kind {
	case apperr.Invalid:
		status = http.StatusBadRequest
	case apperr.NotFound:
		status = http.StatusNotFound
	case apperr.Unauthorized:
		status = http.StatusUnauthorized
	case apperr.Forbidden:
		status = http.StatusForbidden
	case apperr.Unavailable:
		status = http.StatusServiceUnavailable
	}

	writeProblem(w, r, status, string(appError.Kind), appError.Message)
}
