package httpapi

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"streaming-golang/internal/app/apperr"
	"streaming-golang/internal/app/authz"
)

func TestLicenseValidationValidatesIDsAndRestoresBody(t *testing.T) {
	validator := &recordingLicenseValidator{}
	body := `[{"ids":[10,20,10]}]`

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		read, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read restored body: %v", err)
		}
		if string(read) != body {
			t.Fatalf("restored body = %q, want %q", string(read), body)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/validation/curves", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), correlationIDKey, "corr-1"))
	rec := httptest.NewRecorder()

	licenseValidation(validator, "development")(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if validator.request.Stage != "validation" {
		t.Fatalf("stage = %q, want validation", validator.request.Stage)
	}
	if validator.request.InternalCorrelationID != "corr-1" {
		t.Fatalf("correlation = %q, want corr-1", validator.request.InternalCorrelationID)
	}
	if got, want := validator.request.Identifiers, []int64{10, 20}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("identifiers = %#v, want %#v", got, want)
	}
}

func TestLicenseValidationStopsForbiddenRequests(t *testing.T) {
	validator := &recordingLicenseValidator{
		err: apperr.New(apperr.Forbidden, "license denied"),
	}
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves", strings.NewReader(`[{"ids":[10]}]`))
	rec := httptest.NewRecorder()

	licenseValidation(validator, "development")(next).ServeHTTP(rec, req)

	if nextCalled {
		t.Fatal("next handler should not be called")
	}
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestLicenseValidationLetsMalformedBodyReachHandler(t *testing.T) {
	validator := &recordingLicenseValidator{}
	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusTeapot)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves", strings.NewReader(`{`))
	rec := httptest.NewRecorder()

	licenseValidation(validator, "development")(next).ServeHTTP(rec, req)

	if !nextCalled {
		t.Fatal("next handler should be called")
	}
	if validator.called {
		t.Fatal("validator should not be called for malformed body")
	}
	if rec.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}

func TestLicenseValidationSupportsGenericRequestBody(t *testing.T) {
	validator := &recordingLicenseValidator{}
	body := `{"Id":30}`

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		read, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read restored body: %v", err)
		}
		if string(read) != body {
			t.Fatalf("restored body = %q, want %q", string(read), body)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/productive/generic", strings.NewReader(body))
	rec := httptest.NewRecorder()

	licenseValidation(validator, "development")(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if validator.request.Stage != "productive" {
		t.Fatalf("stage = %q, want productive", validator.request.Stage)
	}
	if got, want := validator.request.Identifiers, []int64{30}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("identifiers = %#v, want %#v", got, want)
	}
}

func TestLicenseValidationSupportsLiteQuery(t *testing.T) {
	validator := &recordingLicenseValidator{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/validation/lite?id=40&from=2023-01-01T00:00:00", nil)
	rec := httptest.NewRecorder()

	licenseValidation(validator, "development")(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
	if validator.request.Stage != "validation" {
		t.Fatalf("stage = %q, want validation", validator.request.Stage)
	}
	if got, want := validator.request.Identifiers, []int64{40}; len(got) != len(want) || got[0] != want[0] {
		t.Fatalf("identifiers = %#v, want %#v", got, want)
	}
}

type recordingLicenseValidator struct {
	called  bool
	request authz.LicenseRequest
	err     error
}

func (v *recordingLicenseValidator) ValidateReadAccess(_ context.Context, request authz.LicenseRequest) error {
	v.called = true
	v.request = request
	return v.err
}
