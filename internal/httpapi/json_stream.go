package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"streaming-golang/internal/app/transactional"
)

const (
	ndjsonContentType       = "application/x-ndjson"
	defaultStreamFlushEvery = 1000
)

func acceptsNDJSON(r *http.Request) bool {
	return strings.Contains(strings.ToLower(r.Header.Get("Accept")), ndjsonContentType)
}

func writeTransactionalJSONStream(ctx context.Context, w http.ResponseWriter, r *http.Request, stream transactional.Stream, flushEvery int) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	flushEvery = normalizeStreamFlushEvery(flushEvery)

	flusher, _ := w.(http.Flusher)
	encoder := json.NewEncoder(w)

	if _, err := w.Write([]byte("[")); err != nil {
		return err
	}
	if flusher != nil {
		flusher.Flush()
	}

	index := 0
	for stream.Next(ctx) {
		if index > 0 {
			if _, err := w.Write([]byte(",")); err != nil {
				return err
			}
		}
		if err := encoder.Encode(stream.Item()); err != nil {
			return err
		}
		index++
		if shouldFlushStream(index, flushEvery, flusher) {
			flusher.Flush()
		}
	}

	if err := stream.Err(); err != nil {
		if index > 0 {
			if _, writeErr := w.Write([]byte(",")); writeErr != nil {
				return writeErr
			}
		}
		if encodeErr := encoder.Encode(streamErrorObject(err.Error(), index+1, r, "JSON")); encodeErr != nil {
			return encodeErr
		}
	}

	_, err := w.Write([]byte("]"))
	if flusher != nil {
		flusher.Flush()
	}
	return err
}

func writeTransactionalNDJSONStream(ctx context.Context, w http.ResponseWriter, r *http.Request, stream transactional.Stream, flushEvery int) error {
	w.Header().Set("Content-Type", ndjsonContentType)
	flushEvery = normalizeStreamFlushEvery(flushEvery)

	encoder := json.NewEncoder(w)
	flusher, _ := w.(http.Flusher)

	index := 0
	for stream.Next(ctx) {
		if err := encoder.Encode(stream.Item()); err != nil {
			return err
		}
		index++
		if shouldFlushStream(index, flushEvery, flusher) {
			flusher.Flush()
		}
	}

	if err := stream.Err(); err != nil {
		if encodeErr := encoder.Encode(streamErrorObject(err.Error(), index+1, r, "NDJSON")); encodeErr != nil {
			return encodeErr
		}
		if flusher != nil {
			flusher.Flush()
		}
	}

	return nil
}

func shouldFlushStream(index, flushEvery int, flusher http.Flusher) bool {
	return flusher != nil && index > 0 && index%flushEvery == 0
}

func normalizeStreamFlushEvery(value int) int {
	if value <= 0 {
		return defaultStreamFlushEvery
	}
	return value
}

func streamErrorObject(message string, itemIndex int, r *http.Request, streamFormat string) map[string]any {
	return map[string]any{
		"_error":       true,
		"_type":        "stream_termination",
		"_message":     message,
		"_itemIndex":   itemIndex,
		"_timestamp":   time.Now().UTC().Format(time.RFC3339Nano),
		"_requestPath": r.URL.Path,
		"_details": map[string]any{
			"method":       r.Method,
			"contentType":  streamFormat,
			"streamFormat": streamFormat,
		},
	}
}
