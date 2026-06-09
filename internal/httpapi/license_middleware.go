package httpapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"streaming-golang/internal/app/authz"
	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
)

func licenseValidation(validator authz.LicenseValidator, fallbackStage string) func(http.Handler) http.Handler {
	if validator == nil {
		validator = authz.NoopLicenseValidator{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && isLiteRoute(r.URL.Path) {
				identifiers := uniqueLiteIdentifiers(r)
				if len(identifiers) == 0 {
					next.ServeHTTP(w, r)
					return
				}
				if err := validateLicense(r, validator, identifiers, fallbackStage); err != nil {
					writeAppError(w, r, err)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if r.Method != http.MethodPost || r.Body == nil {
				next.ServeHTTP(w, r)
				return
			}

			body, err := io.ReadAll(r.Body)
			_ = r.Body.Close()
			if err != nil {
				writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", "Invalid or empty request body.")
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(body))

			requests, ok := decodeLicenseRequests(body)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			identifiers := uniqueIdentifiers(requests)
			if len(identifiers) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			if err := validateLicense(r, validator, identifiers, fallbackStage); err != nil {
				writeAppError(w, r, err)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
			next.ServeHTTP(w, r)
		})
	}
}

func decodeLicenseRequests(body []byte) ([]transactional.Request, bool) {
	if len(bytes.TrimSpace(body)) == 0 {
		return nil, false
	}

	var requests []transactional.Request
	if err := json.Unmarshal(body, &requests); err == nil {
		return requests, true
	}

	var generic genericRequest
	if err := json.Unmarshal(body, &generic); err != nil {
		return nil, false
	}
	return []transactional.Request{generic.toTransactionalRequest()}, true
}

func uniqueIdentifiers(requests []transactional.Request) []int64 {
	seen := make(map[domain.Identifier]struct{})
	identifiers := make([]int64, 0)
	for _, request := range requests {
		for _, id := range request.IDs {
			if _, exists := seen[id]; exists {
				continue
			}
			seen[id] = struct{}{}
			identifiers = append(identifiers, int64(id))
		}
	}
	return identifiers
}

func uniqueLiteIdentifiers(r *http.Request) []int64 {
	rawID := strings.TrimSpace(r.URL.Query().Get("id"))
	if rawID == "" {
		return nil
	}
	request, err := transactionalRequestFromLiteQuery(r)
	if err != nil || len(request.IDs) == 0 {
		return nil
	}
	return []int64{int64(request.IDs[0])}
}

func isLiteRoute(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return len(parts) > 0 && parts[len(parts)-1] == "lite"
}

func validateLicense(r *http.Request, validator authz.LicenseValidator, identifiers []int64, fallbackStage string) error {
	requestContext := transactionalRequestContext(r.URL.Path, transactional.ModeJSON, fallbackStage)
	return validator.ValidateReadAccess(r.Context(), authz.LicenseRequest{
		Identifiers:           identifiers,
		Stage:                 requestContext.Stage,
		InternalCorrelationID: correlationIDFromContext(r.Context()),
	})
}
