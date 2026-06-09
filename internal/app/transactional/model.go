package transactional

import (
	"encoding/json"
	"time"

	"streaming-golang/internal/domain"
)

type Request struct {
	IDs             []domain.Identifier `json:"ids"`
	VersionAsOf     *time.Time          `json:"versionAsOf,omitempty"`
	Filters         *Filters            `json:"filters,omitempty"`
	Transformations *Transformations    `json:"transformations,omitempty"`
	Columns         []string            `json:"columns,omitempty"`
	IncludeDeleted  *bool               `json:"includeDeleted,omitempty"`
}

type Filters struct {
	Expressions    []string        `json:"expressions,omitempty"`
	FilterTimeZone string          `json:"filterTimeZone,omitempty"`
	Shape          json.RawMessage `json:"shape,omitempty"`
	Parsed         FilterSet       `json:"-"`
}

type Transformations struct {
	Timezone       string     `json:"timezone,omitempty"`
	TargetTimeZone string     `json:"targetTimeZone,omitempty"`
	Offset         *bool      `json:"offset,omitempty"`
	Nested         string     `json:"nested,omitempty"`
	Keys           []string   `json:"keys,omitempty"`
	Values         [][]string `json:"values,omitempty"`
}

type Response struct {
	TransactionalData []DataItem    `json:"transactionalData"`
	ReferenceData     ReferenceData `json:"referenceData"`
}

type ReferenceData []domain.Identifier

type DataItem struct {
	ID     domain.Identifier `json:"id"`
	Fields map[string]any    `json:"fields"`
}

type Command = domain.Command

type FilterSet = domain.FilterSet

type RequestContext struct {
	DataCategory domain.DataCategory
	EndpointKind EndpointKind
	Stage        string
	Mode         ResponseMode
}

type EndpointKind string

const (
	EndpointTransactional EndpointKind = "transactional"
	EndpointGeneric       EndpointKind = "generic"
	EndpointLite          EndpointKind = "lite"
)

type ResponseMode string

const (
	ModeJSON         ResponseMode = "json"
	ModeJSONStream   ResponseMode = "json_stream"
	ModeNDJSONStream ResponseMode = "ndjson_stream"
	ModeCSV          ResponseMode = "csv"
	ModeCSVStream    ResponseMode = "csv_stream"
)
