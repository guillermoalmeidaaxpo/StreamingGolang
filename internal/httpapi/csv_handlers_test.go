package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
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
	for _, want := range []string{"status,source,dataCategory,statement,parameterCount", "planned,cmdp,,pending_query_generation,1"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
}

func TestGenericCSVEndpointIgnoresJSONAccept(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic", strings.NewReader(`{"Id":10}`))
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	body := rec.Body.String()
	for _, want := range []string{"status,source,dataCategory,statement,parameterCount", "planned,cmdp,,pending_query_generation,1"} {
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

func TestGenericCSVStreamingEndpointIgnoresJSONAccept(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic/streaming", strings.NewReader(`{"Ids":[10,20]}`))
	req.Header.Set("Accept", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	body := rec.Body.String()
	for _, want := range []string{"status,source,dataCategory,statement,parameterCount", "planned,cmdp,,pending_query_generation,1"} {
		if !strings.Contains(body, want) {
			t.Fatalf("body missing %q: %s", want, body)
		}
	}
}

func TestGenericCSVStreamingEndpointIgnoresNDJSONAccept(t *testing.T) {
	router := newCSVTestRouter()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/generic/streaming", strings.NewReader(`{"Ids":[10,20]}`))
	req.Header.Set("Accept", "application/x-ndjson")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "text/csv; charset=utf-8" {
		t.Fatalf("content-type = %q", got)
	}
	body := rec.Body.String()
	for _, want := range []string{"status,source,dataCategory,statement,parameterCount", "planned,cmdp,,pending_query_generation,1"} {
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
			transactional.NewExecutor(map[domain.SourceKind]transactional.Repository{
				domain.SourceCMDP: csvTestRepository{},
			}, 0),
		),
	})
}

type csvTestRepository struct{}

func (csvTestRepository) Execute(_ context.Context, query domain.ExecutableQuery) ([]transactional.DataItem, error) {
	return []transactional.DataItem{{
		ID: query.ID,
		Fields: map[string]any{
			"status":         "planned",
			"source":         query.Source,
			"dataCategory":   query.DataCategory,
			"statement":      query.Statement,
			"parameterCount": len(query.Parameters),
		},
	}}, nil
}

func (r csvTestRepository) Stream(ctx context.Context, query domain.ExecutableQuery) (transactional.Stream, error) {
	items, err := r.Execute(ctx, query)
	if err != nil {
		return nil, err
	}
	return &csvTestStream{items: items}, nil
}

type csvTestStream struct {
	items []transactional.DataItem
	index int
	item  transactional.DataItem
}

func (s *csvTestStream) Next(ctx context.Context) bool {
	if ctx.Err() != nil || s.index >= len(s.items) {
		return false
	}
	s.item = s.items[s.index]
	s.index++
	return true
}

func (s *csvTestStream) Item() transactional.DataItem {
	return s.item
}

func (s *csvTestStream) Err() error {
	return nil
}

func (s *csvTestStream) Close() error {
	return nil
}
