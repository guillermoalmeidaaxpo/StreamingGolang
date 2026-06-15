package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
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

	rowIndex := 0
	batchIndex := 0
	batch := newColumnBatch(flushEvery)
	for stream.Next(ctx) {
		item := stream.Item()
		if !batch.canAdd(item) {
			if err := writeJSONStreamBatch(w, encoder, flusher, batch, &batchIndex); err != nil {
				return err
			}
			batch.reset()
		}
		batch.add(item)
		rowIndex++
		if batch.full() {
			if err := writeJSONStreamBatch(w, encoder, flusher, batch, &batchIndex); err != nil {
				return err
			}
			batch.reset()
		}
	}

	if err := writeJSONStreamBatch(w, encoder, flusher, batch, &batchIndex); err != nil {
		return err
	}

	if err := stream.Err(); err != nil {
		if batchIndex > 0 {
			if _, writeErr := w.Write([]byte(",")); writeErr != nil {
				return writeErr
			}
		}
		if encodeErr := encoder.Encode(streamErrorObject(err.Error(), rowIndex+1, r, "JSON")); encodeErr != nil {
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

	rowIndex := 0
	batch := newColumnBatch(flushEvery)
	for stream.Next(ctx) {
		item := stream.Item()
		if !batch.canAdd(item) {
			if err := writeNDJSONStreamBatch(encoder, flusher, batch); err != nil {
				return err
			}
			batch.reset()
		}
		batch.add(item)
		rowIndex++
		if batch.full() {
			if err := writeNDJSONStreamBatch(encoder, flusher, batch); err != nil {
				return err
			}
			batch.reset()
		}
	}

	if err := writeNDJSONStreamBatch(encoder, flusher, batch); err != nil {
		return err
	}

	if err := stream.Err(); err != nil {
		if encodeErr := encoder.Encode(streamErrorObject(err.Error(), rowIndex+1, r, "NDJSON")); encodeErr != nil {
			return encodeErr
		}
		if flusher != nil {
			flusher.Flush()
		}
	}

	return nil
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

type columnBatch struct {
	identifier domain.Identifier
	limit      int
	count      int
	fields     map[string][]any
}

func newColumnBatch(limit int) *columnBatch {
	return &columnBatch{
		limit:  normalizeStreamFlushEvery(limit),
		fields: make(map[string][]any),
	}
}

func (b *columnBatch) canAdd(item transactional.DataItem) bool {
	return b.count == 0 || b.identifier == item.ID
}

func (b *columnBatch) add(item transactional.DataItem) {
	if b.count == 0 {
		b.identifier = item.ID
	}

	seen := make(map[string]struct{}, len(item.Fields))
	for key, value := range item.Fields {
		if isIdentifierField(key) {
			continue
		}
		if _, ok := b.fields[key]; !ok {
			b.fields[key] = make([]any, b.count)
		}
		b.fields[key] = append(b.fields[key], value)
		seen[key] = struct{}{}
	}
	for key := range b.fields {
		if _, ok := seen[key]; !ok {
			b.fields[key] = append(b.fields[key], nil)
		}
	}
	b.count++
}

func (b *columnBatch) full() bool {
	return b.count >= b.limit
}

func (b *columnBatch) empty() bool {
	return b.count == 0
}

func (b *columnBatch) reset() {
	b.count = 0
	b.fields = make(map[string][]any)
}

func (b *columnBatch) object() map[string]any {
	value := make(map[string]any, len(b.fields)+1)
	value["Identifier"] = int64(b.identifier)
	for key, values := range b.fields {
		value[key] = values
	}
	return value
}

func isIdentifierField(key string) bool {
	return strings.EqualFold(key, "Identifier") || strings.EqualFold(key, "MdoId")
}

func writeJSONStreamBatch(w http.ResponseWriter, encoder *json.Encoder, flusher http.Flusher, batch *columnBatch, batchIndex *int) error {
	if batch.empty() {
		return nil
	}
	if *batchIndex > 0 {
		if _, err := w.Write([]byte(",")); err != nil {
			return err
		}
	}
	if err := encoder.Encode(batch.object()); err != nil {
		return err
	}
	(*batchIndex)++
	if flusher != nil {
		flusher.Flush()
	}
	return nil
}

func writeNDJSONStreamBatch(encoder *json.Encoder, flusher http.Flusher, batch *columnBatch) error {
	if batch.empty() {
		return nil
	}
	if err := encoder.Encode(batch.object()); err != nil {
		return err
	}
	if flusher != nil {
		flusher.Flush()
	}
	return nil
}
