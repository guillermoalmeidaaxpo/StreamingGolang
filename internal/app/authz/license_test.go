package authz

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHttpLicenseValidatorMatchesCSharpAuthorizationFlow(t *testing.T) {
	seen := make([]string, 0, 2)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.URL.Path)
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("authorization header = %q, want bearer token", got)
		}

		switch r.URL.Path {
		case "/api/v1/DataUniverse/BulkAuthorize":
			var body bulkDataUniverseRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode universe request: %v", err)
			}
			if body.Type != "TransactionalDataOutbound" || body.Action != "Read" || body.StageID != 3 {
				t.Fatalf("universe body = %#v", body)
			}
			if len(body.MDOIDs) != 1 || body.MDOIDs[0] != 40 {
				t.Fatalf("universe mdoIds = %#v, want [40]", body.MDOIDs)
			}
			_, _ = w.Write([]byte(`{"40":true}`))
		case "/api/v1/TimeSeries/Authorize":
			var body timeSeriesRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode license request: %v", err)
			}
			if len(body.Identifiers) != 1 || body.Identifiers[0] != 40 || body.StageID != 3 {
				t.Fatalf("license body = %#v", body)
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	validator := NewHttpLicenseValidator(server.URL, "api/v1/TimeSeries/Authorize", "api/v1/DataUniverse/BulkAuthorize", time.Second, slog.Default())
	ctx := contextWithRawToken("test-token")

	err := validator.ValidateReadAccess(ctx, LicenseRequest{
		Identifiers:           []int64{40},
		Stage:                 "development",
		InternalCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatalf("validate read access failed: %v", err)
	}
	if got := strings.Join(seen, ","); got != "/api/v1/DataUniverse/BulkAuthorize,/api/v1/TimeSeries/Authorize" {
		t.Fatalf("paths = %s", got)
	}
}

func TestHttpLicenseValidatorRejectsUniverseDenial(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/DataUniverse/BulkAuthorize" {
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"40":false}`))
	}))
	defer server.Close()

	validator := NewHttpLicenseValidator(server.URL, "api/v1/TimeSeries/Authorize", "api/v1/DataUniverse/BulkAuthorize", time.Second, slog.Default())

	err := validator.ValidateReadAccess(contextWithRawToken("test-token"), LicenseRequest{
		Identifiers: []int64{40},
		Stage:       "productive",
	})
	if err == nil {
		t.Fatal("expected forbidden error")
	}
}

func TestHttpLicenseValidatorStopsWhenUniverseNotFound(t *testing.T) {
	seen := make([]string, 0, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.URL.Path)
		switch r.URL.Path {
		case "/api/v1/DataUniverse/BulkAuthorize":
			http.NotFound(w, r)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	validator := NewHttpLicenseValidator(server.URL, "api/v1/TimeSeries/Authorize", "api/v1/DataUniverse/BulkAuthorize", time.Second, slog.Default())

	err := validator.ValidateReadAccess(contextWithRawToken("test-token"), LicenseRequest{
		Identifiers: []int64{40},
		Stage:       "productive",
	})
	if err == nil {
		t.Fatal("expected universe authorization error")
	}
	if got := strings.Join(seen, ","); got != "/api/v1/DataUniverse/BulkAuthorize" {
		t.Fatalf("paths = %s", got)
	}
}

func contextWithRawToken(token string) context.Context {
	return context.WithValue(context.Background(), "raw_bearer_token", token)
}
