package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"streaming-golang/internal/app/transactional"
)

func TestTransactionalEndpointReturnsReferenceDataAsArray(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves", strings.NewReader(`[{"ids":[10,20]}]`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	referenceData, ok := body["referenceData"].([]any)
	if !ok {
		t.Fatalf("referenceData = %#v, want array", body["referenceData"])
	}
	if len(referenceData) != 2 || referenceData[0].(float64) != 10 || referenceData[1].(float64) != 20 {
		t.Fatalf("referenceData = %#v, want [10,20]", referenceData)
	}
}

func TestTransactionalStreamingDefaultsToJSONArray(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves/streaming", strings.NewReader(`[{"ids":[10,20]}]`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	body := strings.TrimSpace(rec.Body.String())
	if !strings.HasPrefix(body, "[") || !strings.HasSuffix(body, "]") {
		t.Fatalf("body is not JSON array: %s", body)
	}
	var items []map[string]any
	if err := json.Unmarshal([]byte(body), &items); err != nil {
		t.Fatalf("decode JSON stream body: %v; body=%s", err, body)
	}
	if len(items) != 2 {
		t.Fatalf("items = %d, want 2", len(items))
	}
}

func TestTransactionalStreamingNegotiatesNDJSON(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves/streaming", strings.NewReader(`[{"ids":[10,20]}]`))
	req.Header.Set("Accept", "application/x-ndjson")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/x-ndjson" {
		t.Fatalf("content-type = %q", got)
	}
	body := strings.TrimSpace(rec.Body.String())
	if strings.HasPrefix(body, "[") {
		t.Fatalf("ndjson must not start with JSON array bracket: %s", body)
	}
	lines := strings.Split(body, "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2; body=%s", len(lines), body)
	}
	for _, line := range lines {
		var item map[string]any
		if err := json.Unmarshal([]byte(line), &item); err != nil {
			t.Fatalf("decode ndjson line: %v; line=%s", err, line)
		}
	}
}

func TestTransactionalStreamingAcceptsNullableTransformations(t *testing.T) {
	router := newCSVTestRouter()

	body := `[{"ids":[1000000001],"filters":{"filterTimeZone":"CET","expressions":["ReferenceTime >= now()+P1D"]},"transformations":null,"columns":["CreatedOn"]}]`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/surfaces/streaming", strings.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestTransactionalEndpointRejectsUnknownSchemaFields(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves", strings.NewReader(`[{"ids":[10],"unknown":true}]`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestGenericEndpointRejectsUnknownSchemaFields(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic", strings.NewReader(`{"id":10,"unknown":true}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestTransactionalEndpointRejectsBodyOutsideSchema(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves", strings.NewReader(`{"ids":[10]}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestGenericEndpointRejectsInvalidSchemaTypes(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic", strings.NewReader(`{"id":"10"}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestGenericRequestContextDoesNotForceDataCategory(t *testing.T) {
	requestContext := genericRequestContext("/api/v1/generic", transactional.ModeCSV, "development")
	if requestContext.EndpointKind != transactional.EndpointGeneric {
		t.Fatalf("endpoint kind = %q, want generic", requestContext.EndpointKind)
	}
	if requestContext.DataCategory != "" {
		t.Fatalf("data category = %q, want empty for generic", requestContext.DataCategory)
	}
}
