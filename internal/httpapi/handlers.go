package httpapi

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
	"streaming-golang/internal/platform/config"
)

type handlers struct {
	config                config.Config
	logger                *slog.Logger
	streamFlushEvery      int
	transactionalPipeline *transactional.Pipeline
}

func (h handlers) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h handlers) info(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"buildNumber": h.config.Build.Number,
		"stage":       h.config.Build.Stage,
	})
}

func (h handlers) transactional(w http.ResponseWriter, r *http.Request) {
	decodeStart := time.Now()
	var req []transactional.Request
	if err := decodeSchemaJSON(r.Body, transactionalRequestSchema, &req); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", invalidRequestBodyDetail("Invalid or empty transactional data body.", err))
		return
	}
	h.logPhase(r, "request body decoded", decodeStart,
		slog.String("handler", "transactional"),
		slog.Int("request_count", len(req)),
	)

	requestContext := transactionalRequestContext(r.URL.Path, transactional.ModeJSON, h.config.Build.Stage)
	pipelineStart := time.Now()
	result, err := h.transactionalPipeline.Execute(r.Context(), requestContext, req)
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	h.logPhase(r, "handler pipeline completed", pipelineStart,
		slog.String("handler", "transactional"),
		slog.String("mode", string(requestContext.Mode)),
		slog.Int("row_count", len(result.TransactionalData)),
	)

	writeStart := time.Now()
	writeJSON(w, http.StatusOK, result)
	h.logPhase(r, "response written", writeStart,
		slog.String("handler", "transactional"),
		slog.String("format", "json"),
		slog.Int("row_count", len(result.TransactionalData)),
	)
}

func (h handlers) transactionalStream(w http.ResponseWriter, r *http.Request) {
	decodeStart := time.Now()
	var req []transactional.Request
	if err := decodeSchemaJSON(r.Body, transactionalRequestSchema, &req); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", invalidRequestBodyDetail("Invalid or empty transactional stream body.", err))
		return
	}
	h.logPhase(r, "request body decoded", decodeStart,
		slog.String("handler", "transactional_stream"),
		slog.Int("request_count", len(req)),
	)

	mode := transactional.ModeJSONStream
	if acceptsNDJSON(r) {
		mode = transactional.ModeNDJSONStream
	}
	requestContext := transactionalRequestContext(r.URL.Path, mode, h.config.Build.Stage)
	pipelineStart := time.Now()
	stream, err := h.transactionalPipeline.Stream(r.Context(), requestContext, req)
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	defer stream.Close()
	h.logPhase(r, "handler pipeline completed", pipelineStart,
		slog.String("handler", "transactional_stream"),
		slog.String("mode", string(requestContext.Mode)),
	)

	writeStart := time.Now()
	if acceptsNDJSON(r) {
		if err := writeTransactionalNDJSONStream(r.Context(), w, r, stream, h.streamFlushEvery); err != nil && !errors.Is(err, r.Context().Err()) {
			return
		}
		h.logPhase(r, "response written", writeStart,
			slog.String("handler", "transactional_stream"),
			slog.String("format", "ndjson"),
		)
		return
	}

	if err := writeTransactionalJSONStream(r.Context(), w, r, stream, h.streamFlushEvery); err != nil && !errors.Is(err, r.Context().Err()) {
		return
	}
	h.logPhase(r, "response written", writeStart,
		slog.String("handler", "transactional_stream"),
		slog.String("format", "json_stream"),
	)
}

func (h handlers) genericCSV(w http.ResponseWriter, r *http.Request) {
	decodeStart := time.Now()
	request, ok := decodeGenericRequest(w, r, "Invalid or empty generic body.")
	if !ok {
		return
	}
	h.logPhase(r, "request body decoded", decodeStart,
		slog.String("handler", "generic_csv"),
		slog.Int("id_count", len(request.IDs)),
	)

	requestContext := genericRequestContext(r.URL.Path, transactional.ModeCSV, h.config.Build.Stage)
	pipelineStart := time.Now()
	plan, result, err := h.transactionalPipeline.ExecuteWithPlan(r.Context(), requestContext, []transactional.Request{request})
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	h.logPhase(r, "handler pipeline completed", pipelineStart,
		slog.String("handler", "generic_csv"),
		slog.String("mode", string(requestContext.Mode)),
		slog.Int("plan_steps", len(plan.Steps)),
		slog.Int("row_count", len(result.TransactionalData)),
	)

	writeStart := time.Now()
	if err := writeTransactionalCSV(w, result, csvColumnsFromPlan(plan), csvIncludeOffset(plan), true); err != nil {
		return
	}
	h.logPhase(r, "response written", writeStart,
		slog.String("handler", "generic_csv"),
		slog.String("format", "csv"),
		slog.Int("row_count", len(result.TransactionalData)),
	)
}

func (h handlers) genericCSVStream(w http.ResponseWriter, r *http.Request) {
	decodeStart := time.Now()
	request, ok := decodeGenericRequest(w, r, "Invalid or empty generic stream body.")
	if !ok {
		return
	}
	h.logPhase(r, "request body decoded", decodeStart,
		slog.String("handler", "generic_csv_stream"),
		slog.Int("id_count", len(request.IDs)),
	)

	requestContext := genericRequestContext(r.URL.Path, transactional.ModeCSVStream, h.config.Build.Stage)
	pipelineStart := time.Now()
	plan, stream, err := h.transactionalPipeline.StreamWithPlan(r.Context(), requestContext, []transactional.Request{request})
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	defer stream.Close()
	h.logPhase(r, "handler pipeline completed", pipelineStart,
		slog.String("handler", "generic_csv_stream"),
		slog.String("mode", string(requestContext.Mode)),
		slog.Int("plan_steps", len(plan.Steps)),
	)

	writeStart := time.Now()
	if err := writeTransactionalCSVStream(r.Context(), w, stream, csvColumnsFromPlan(plan), csvIncludeOffset(plan), false, h.streamFlushEvery); err != nil && !errors.Is(err, r.Context().Err()) {
		return
	}
	h.logPhase(r, "response written", writeStart,
		slog.String("handler", "generic_csv_stream"),
		slog.String("format", "csv_stream"),
	)
}

func (h handlers) liteCSV(w http.ResponseWriter, r *http.Request) {
	decodeStart := time.Now()
	if err := validateLiteQuerySchema(r.URL.Query()); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-query", "Invalid or empty lite query.")
		return
	}

	request, err := transactionalRequestFromLiteQuery(r)
	if err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-query", "Invalid or empty lite query.")
		return
	}
	h.logPhase(r, "request query decoded", decodeStart,
		slog.String("handler", "lite_csv"),
		slog.Int("id_count", len(request.IDs)),
	)

	requestContext := transactionalRequestContext(r.URL.Path, transactional.ModeCSV, h.config.Build.Stage)
	pipelineStart := time.Now()
	plan, result, err := h.transactionalPipeline.ExecuteWithPlan(r.Context(), requestContext, []transactional.Request{request})
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	h.logPhase(r, "handler pipeline completed", pipelineStart,
		slog.String("handler", "lite_csv"),
		slog.String("mode", string(requestContext.Mode)),
		slog.Int("plan_steps", len(plan.Steps)),
		slog.Int("row_count", len(result.TransactionalData)),
	)

	writeStart := time.Now()
	if err := writeTransactionalCSV(w, result, csvColumnsFromPlan(plan), csvIncludeOffset(plan), true); err != nil {
		return
	}
	h.logPhase(r, "response written", writeStart,
		slog.String("handler", "lite_csv"),
		slog.String("format", "csv"),
		slog.Int("row_count", len(result.TransactionalData)),
	)
}

func (h handlers) notImplemented(feature string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeProblem(w, r, http.StatusNotImplemented, "not-implemented", feature+" endpoint is registered but not implemented yet.")
	}
}

func transactionalRequestContext(path string, mode transactional.ResponseMode, fallbackStage string) transactional.RequestContext {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	stage := fallbackStage
	category := domain.TimeSeries
	endpointKind := transactional.EndpointTransactional

	for _, part := range parts {
		switch part {
		case "design", "validation", "productive", "migration":
			stage = part
		case "curves":
			category = domain.Curves
		case "surfaces":
			category = domain.Surfaces
		case "timeseries":
			category = domain.TimeSeries
		case "generic":
			endpointKind = transactional.EndpointGeneric
		case "lite":
			endpointKind = transactional.EndpointLite
		}
	}

	return transactional.RequestContext{
		DataCategory: category,
		EndpointKind: endpointKind,
		Stage:        stage,
		Mode:         mode,
	}
}

func genericRequestContext(path string, mode transactional.ResponseMode, fallbackStage string) transactional.RequestContext {
	requestContext := transactionalRequestContext(path, mode, fallbackStage)
	requestContext.DataCategory = ""
	requestContext.EndpointKind = transactional.EndpointGeneric
	return requestContext
}

func decodeGenericRequest(w http.ResponseWriter, r *http.Request, detail string) (transactional.Request, bool) {
	var request genericRequest
	if err := decodeSchemaJSON(r.Body, genericRequestSchema, &request); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", invalidRequestBodyDetail(detail, err))
		return transactional.Request{}, false
	}
	return request.toTransactionalRequest(), true
}

func invalidRequestBodyDetail(prefix string, err error) string {
	if err == nil {
		return prefix
	}
	return prefix + " " + err.Error()
}

func (h handlers) logPhase(r *http.Request, message string, start time.Time, attrs ...slog.Attr) {
	if h.logger == nil {
		return
	}
	duration := time.Since(start)
	logAttrs := make([]slog.Attr, 0, len(attrs)+4)
	logAttrs = append(logAttrs, attrs...)
	logAttrs = append(logAttrs,
		slog.String("path", r.URL.Path),
		slog.String("correlation_id", correlationIDFromContext(r.Context())),
		slog.Duration("duration", duration),
		slog.Int64("duration_ms", duration.Milliseconds()),
	)
	h.logger.LogAttrs(r.Context(), slog.LevelInfo, message, logAttrs...)
}
