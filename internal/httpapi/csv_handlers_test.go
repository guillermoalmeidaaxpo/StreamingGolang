package httpapi

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/platform/config"
	antlrparser "streaming-golang/internal/query/parser/antlr"
)

func TestGenericCSVEndpointExecutesTransactionalFlow(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic", strings.NewReader(`{"Id":10}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	if got := rec.Header().Get("Content-Disposition"); got != `attachment; filename="transactional_data.csv"` {
		t.Fatalf("content-disposition = %q", got)
	}
	body := rec.Body.String()
	for _, want := range []string{"status,source,dataCategory,statement,parameterCount", "planned,cmdp,curves,pending_query_generation,1"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
}

func TestGenericCSVStreamingEndpointExecutesTransactionalFlow(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic/streaming", strings.NewReader(`{"Ids":[10,20],"Columns":["ReferenceTime","Value"]}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	if got := rec.Header().Get("Content-Disposition"); got != "" {
		t.Fatalf("content-disposition = %q", got)
	}
	body := rec.Body.String()
	for _, want := range []string{"ReferenceTime,Value", "N/A,N/A"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
}

func TestLiteCSVEndpointBuildsReferenceTimeFilters(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lite?id=10&from=2023-01-01T00:00:00&to=2023-01-02T00:00:00", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	body := rec.Body.String()
	for _, want := range []string{"status,source,dataCategory,statement,parameterCount", "planned,cmdp,timeseries,pending_query_generation,1"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
}

func TestLiteCSVEndpointRejectsMissingRequiredQuery(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lite?id=10", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestLiteCSVEndpointRejectsRFC3339Date(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lite?id=10&from=2023-01-01T00:00:00Z", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestLiteCSVEndpointRejectsInvalidRange(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lite?id=10&from=2023-01-02T00:00:00&to=2023-01-01T00:00:00", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestLiteCSVEndpointRejectsUnknownQueryParameter(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/lite?id=10&from=2023-01-01T00:00:00&extra=true", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func newCSVTestRouter() http.Handler {
	return NewRouter(Dependencies{
		Config: config.Config{},
		Logger: slog.Default(),
		TransactionalPipeline: transactional.NewPipeline(
			transactional.NewValidator(),
			antlrparser.New(),
			transactional.NewPlanner(),
			transactional.NewExecutor(nil, 0),
		),

	})
}
