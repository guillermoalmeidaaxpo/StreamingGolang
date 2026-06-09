package httpapi

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"sync"

	"github.com/xeipuuv/gojsonschema"
)

const (
	transactionalRequestSchema = "transactional.schema.json"
	genericRequestSchema       = "generic.schema.json"
	liteRequestSchema          = "lite.schema.json"
)

//go:embed schemas/*.json
var schemaFiles embed.FS

var (
	schemaOnce  sync.Once
	schemaCache map[string]*gojsonschema.Schema
	schemaErr   error
)

type requestSchemaError struct {
	messages []string
}

func (e requestSchemaError) Error() string {
	return "request schema validation failed: " + strings.Join(e.messages, "; ")
}

func decodeSchemaJSON(body io.Reader, schemaName string, target any) error {
	raw, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if len(bytes.TrimSpace(raw)) == 0 {
		return io.EOF
	}
	if err := validateJSONBytes(schemaName, raw); err != nil {
		return err
	}
	return decodeStrictJSON(bytes.NewReader(raw), target)
}

func validateLiteQuerySchema(values url.Values) error {
	document := make(map[string]any, len(values))
	for key, rawValues := range values {
		if len(rawValues) > 1 {
			return requestSchemaError{messages: []string{fmt.Sprintf("%s: repeated query parameter", key)}}
		}
		if len(rawValues) == 1 {
			document[key] = rawValues[0]
		}
	}
	return validateJSONDocument(liteRequestSchema, document)
}

func validateJSONBytes(schemaName string, raw []byte) error {
	return validateJSON(schemaName, gojsonschema.NewBytesLoader(raw))
}

func validateJSONDocument(schemaName string, document any) error {
	return validateJSON(schemaName, gojsonschema.NewGoLoader(document))
}

func validateJSON(schemaName string, document gojsonschema.JSONLoader) error {
	schema, err := requestSchema(schemaName)
	if err != nil {
		return err
	}
	result, err := schema.Validate(document)
	if err != nil {
		return err
	}
	if result.Valid() {
		return nil
	}

	messages := make([]string, 0, len(result.Errors()))
	for _, validationError := range result.Errors() {
		messages = append(messages, validationError.String())
	}
	sort.Strings(messages)
	return requestSchemaError{messages: messages}
}

func requestSchema(name string) (*gojsonschema.Schema, error) {
	schemaOnce.Do(func() {
		schemaCache = make(map[string]*gojsonschema.Schema)
		for _, schemaName := range []string{transactionalRequestSchema, genericRequestSchema, liteRequestSchema} {
			raw, err := schemaFiles.ReadFile("schemas/" + schemaName)
			if err != nil {
				schemaErr = err
				return
			}
			compiled, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(raw))
			if err != nil {
				schemaErr = err
				return
			}
			schemaCache[schemaName] = compiled
		}
	})
	if schemaErr != nil {
		return nil, schemaErr
	}
	schema, ok := schemaCache[name]
	if !ok {
		return nil, fmt.Errorf("unknown request schema %q", name)
	}
	return schema, nil
}
