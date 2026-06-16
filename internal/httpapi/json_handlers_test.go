package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestTransactionalJSONStreamBatchesColumnWiseByFlushSize(t *testing.T) {
	zurich := time.FixedZone("CEST", 2*60*60)
	stream := &csvTestStream{items: []transactional.DataItem{
		{ID: 536013751, Fields: map[string]any{"ReferenceTime": time.Date(2024, 4, 26, 0, 0, 0, 0, zurich), "Value": 1.1}},
		{ID: 536013751, Fields: map[string]any{"ReferenceTime": time.Date(2024, 4, 26, 0, 0, 0, 0, zurich), "Value": 2.2}},
		{ID: 536013751, Fields: map[string]any{"ReferenceTime": time.Date(2024, 4, 26, 0, 0, 0, 0, zurich), "Value": 3.3}},
	}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves/streaming", nil)
	rec := httptest.NewRecorder()

	stats, err := writeTransactionalJSONStream(context.Background(), rec, req, stream, 2)
	if err != nil {
		t.Fatalf("write stream: %v", err)
	}
	if stats.Rows != 3 || stats.Batches != 2 {
		t.Fatalf("stats = %+v, want 3 rows and 2 batches", stats)
	}

	var batches []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &batches); err != nil {
		t.Fatalf("decode JSON stream body: %v; body=%s", err, rec.Body.String())
	}
	if len(batches) != 2 {
		t.Fatalf("batches = %d, want 2; body=%s", len(batches), rec.Body.String())
	}
	if got := batches[0]["Identifier"]; got != float64(536013751) {
		t.Fatalf("Identifier = %#v, want 536013751", got)
	}
	referenceTimes, ok := batches[0]["ReferenceTime"].([]any)
	if !ok || len(referenceTimes) != 2 || referenceTimes[0] != "2024-04-26T00:00:00.000+02:00" {
		t.Fatalf("first batch ReferenceTime = %#v, want C# DateTimeOffset with milliseconds", batches[0]["ReferenceTime"])
	}
	values, ok := batches[0]["Value"].([]any)
	if !ok || len(values) != 2 {
		t.Fatalf("first batch Value = %#v, want array length 2", batches[0]["Value"])
	}
	values, ok = batches[1]["Value"].([]any)
	if !ok || len(values) != 1 {
		t.Fatalf("second batch Value = %#v, want array length 1", batches[1]["Value"])
	}
}

func TestTransactionalNDJSONStreamBatchesColumnWiseByFlushSize(t *testing.T) {
	stream := &csvTestStream{items: []transactional.DataItem{
		{ID: 536013751, Fields: map[string]any{"ReferenceTime": "2024-04-26T00:00:00Z", "Value": 1.1}},
		{ID: 536013751, Fields: map[string]any{"ReferenceTime": "2024-04-26T00:00:00Z", "Value": 2.2}},
		{ID: 536013751, Fields: map[string]any{"ReferenceTime": "2024-04-26T00:00:00Z", "Value": 3.3}},
	}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves/streaming", nil)
	rec := httptest.NewRecorder()

	stats, err := writeTransactionalNDJSONStream(context.Background(), rec, req, stream, 2)
	if err != nil {
		t.Fatalf("write stream: %v", err)
	}
	if stats.Rows != 3 || stats.Batches != 2 {
		t.Fatalf("stats = %+v, want 3 rows and 2 batches", stats)
	}

	lines := strings.Split(strings.TrimSpace(rec.Body.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("lines = %d, want 2; body=%s", len(lines), rec.Body.String())
	}
	for _, line := range lines {
		var batch map[string]any
		if err := json.Unmarshal([]byte(line), &batch); err != nil {
			t.Fatalf("decode ndjson line: %v; line=%s", err, line)
		}
		if _, ok := batch["Value"].([]any); !ok {
			t.Fatalf("Value = %#v, want array", batch["Value"])
		}
	}
	var firstBatch map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &firstBatch); err != nil {
		t.Fatalf("decode first ndjson line: %v; line=%s", err, lines[0])
	}
	values := firstBatch["Value"].([]any)
	if len(values) != 2 {
		t.Fatalf("first ndjson batch Value length = %d, want 2; batch=%#v", len(values), firstBatch)
	}
}

func TestTransactionalJSONStreamStartsNewBatchWhenIdentifierChanges(t *testing.T) {
	stream := &csvTestStream{items: []transactional.DataItem{
		{ID: 10, Fields: map[string]any{"Value": 1}},
		{ID: 20, Fields: map[string]any{"Value": 2}},
	}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/curves/streaming", nil)
	rec := httptest.NewRecorder()

	stats, err := writeTransactionalJSONStream(context.Background(), rec, req, stream, 1000)
	if err != nil {
		t.Fatalf("write stream: %v", err)
	}
	if stats.Rows != 2 || stats.Batches != 2 {
		t.Fatalf("stats = %+v, want 2 rows and 2 batches", stats)
	}

	var batches []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &batches); err != nil {
		t.Fatalf("decode JSON stream body: %v; body=%s", err, rec.Body.String())
	}
	if len(batches) != 2 {
		t.Fatalf("batches = %d, want 2; body=%s", len(batches), rec.Body.String())
	}
	if batches[0]["Identifier"] != float64(10) || batches[1]["Identifier"] != float64(20) {
		t.Fatalf("identifiers = %#v, %#v; want 10, 20", batches[0]["Identifier"], batches[1]["Identifier"])
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

func TestTransactionalStreamingAcceptsNullableNestedTransformation(t *testing.T) {
	router := newCSVTestRouter()

	body := `[{"ids":[504078501],"filters":{"filterTimeZone":"CET","expressions":["ReferenceTime = LatestGlobal()"]},"transformations":{"targetTimeZone":"CET","nested":null},"columns":["CreatedOn"]}]`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/timeseries/streaming", strings.NewReader(body))
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
