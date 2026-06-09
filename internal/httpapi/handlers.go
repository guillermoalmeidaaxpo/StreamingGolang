package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
	"streaming-golang/internal/platform/config"
)

type handlers struct {
	config                config.Config
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
	var req []transactional.Request
	if err := decodeSchemaJSON(r.Body, transactionalRequestSchema, &req); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", "Invalid or empty transactional data body.")
		return
	}

	requestContext := transactionalRequestContext(r.URL.Path, transactional.ModeJSON, h.config.Build.Stage)
	result, err := h.transactionalPipeline.Execute(r.Context(), requestContext, req)
	if err != nil {
		writeAppError(w, r, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h handlers) transactionalStream(w http.ResponseWriter, r *http.Request) {
	var req []transactional.Request
	if err := decodeSchemaJSON(r.Body, transactionalRequestSchema, &req); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", "Invalid or empty transactional stream body.")
		return
	}

	mode := transactional.ModeJSONStream
	if acceptsNDJSON(r) {
		mode = transactional.ModeNDJSONStream
	}
	requestContext := transactionalRequestContext(r.URL.Path, mode, h.config.Build.Stage)
	stream, err := h.transactionalPipeline.Stream(r.Context(), requestContext, req)
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	defer stream.Close()

	if acceptsNDJSON(r) {
		if err := writeTransactionalNDJSONStream(r.Context(), w, r, stream); err != nil && !errors.Is(err, r.Context().Err()) {
			return
		}
		return
	}

	if err := writeTransactionalJSONStream(r.Context(), w, r, stream); err != nil && !errors.Is(err, r.Context().Err()) {
		return
	}
}

func (h handlers) genericCSV(w http.ResponseWriter, r *http.Request) {
	request, ok := decodeGenericRequest(w, r, "Invalid or empty generic body.")
	if !ok {
		return
	}

	requestContext := genericRequestContext(r.URL.Path, transactional.ModeCSV, h.config.Build.Stage)
	plan, result, err := h.transactionalPipeline.ExecuteWithPlan(r.Context(), requestContext, []transactional.Request{request})
	if err != nil {
		writeAppError(w, r, err)
		return
	}

	if err := writeTransactionalCSV(w, result, csvColumnsFromPlan(plan), csvIncludeOffset(plan), true); err != nil {
		return
	}
}

func (h handlers) genericCSVStream(w http.ResponseWriter, r *http.Request) {
	request, ok := decodeGenericRequest(w, r, "Invalid or empty generic stream body.")
	if !ok {
		return
	}

	requestContext := genericRequestContext(r.URL.Path, transactional.ModeCSVStream, h.config.Build.Stage)
	plan, stream, err := h.transactionalPipeline.StreamWithPlan(r.Context(), requestContext, []transactional.Request{request})
	if err != nil {
		writeAppError(w, r, err)
		return
	}
	defer stream.Close()

	if err := writeTransactionalCSVStream(r.Context(), w, stream, csvColumnsFromPlan(plan), csvIncludeOffset(plan), false); err != nil && !errors.Is(err, r.Context().Err()) {
		return
	}
}

func (h handlers) liteCSV(w http.ResponseWriter, r *http.Request) {
	if err := validateLiteQuerySchema(r.URL.Query()); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-query", "Invalid or empty lite query.")
		return
	}

	request, err := transactionalRequestFromLiteQuery(r)
	if err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-query", "Invalid or empty lite query.")
		return
	}

	requestContext := transactionalRequestContext(r.URL.Path, transactional.ModeCSV, h.config.Build.Stage)
	plan, result, err := h.transactionalPipeline.ExecuteWithPlan(r.Context(), requestContext, []transactional.Request{request})
	if err != nil {
		writeAppError(w, r, err)
		return
	}

	if err := writeTransactionalCSV(w, result, csvColumnsFromPlan(plan), csvIncludeOffset(plan), true); err != nil {
		return
	}
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
	requestContext.DataCategory = domain.Curves
	requestContext.EndpointKind = transactional.EndpointGeneric
	return requestContext
}

func decodeGenericRequest(w http.ResponseWriter, r *http.Request, detail string) (transactional.Request, bool) {
	var request genericRequest
	if err := decodeSchemaJSON(r.Body, genericRequestSchema, &request); err != nil {
		writeProblem(w, r, http.StatusBadRequest, "invalid-request-body", detail)
		return transactional.Request{}, false
	}
	return request.toTransactionalRequest(), true
}
