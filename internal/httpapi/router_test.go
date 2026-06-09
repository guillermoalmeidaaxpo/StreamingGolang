package httpapi

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/platform/config"
	antlrparser "streaming-golang/internal/query/parser/antlr"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter(Dependencies{
		Config: config.Config{},
		Logger: slog.Default(),
		TransactionalPipeline: transactional.NewPipeline(
			transactional.NewValidator(),
			antlrparser.New(),
			transactional.NewPlanner(),
			transactional.NewExecutor(nil, 0),
		),

	})

	req := httptest.NewRequest(http.MethodGet, "/health/liveness", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
