package httpapi

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"streaming-golang/internal/app/transactional"
	"streaming-golang/internal/domain"
)

type genericRequest struct {
	ID              *domain.Identifier             `json:"id,omitempty"`
	IDs             []domain.Identifier            `json:"ids,omitempty"`
	VersionAsOf     *time.Time                     `json:"versionAsOf,omitempty"`
	Filters         *transactional.Filters         `json:"filters,omitempty"`
	Transformations *transactional.Transformations `json:"transformations,omitempty"`
	Columns         []string                       `json:"columns,omitempty"`
	IncludeDeleted  *bool                          `json:"includeDeleted,omitempty"`
}

func (r genericRequest) toTransactionalRequest() transactional.Request {
	ids := append([]domain.Identifier(nil), r.IDs...)
	if len(ids) == 0 && r.ID != nil {
		ids = []domain.Identifier{*r.ID}
	}
	transformations := normalizeGenericTransformations(r.Transformations)

	return transactional.Request{
		IDs:             ids,
		VersionAsOf:     r.VersionAsOf,
		Filters:         r.Filters,
		Transformations: transformations,
		Columns:         append([]string(nil), r.Columns...),
		IncludeDeleted:  r.IncludeDeleted,
	}
}

func normalizeGenericTransformations(transformations *transactional.Transformations) *transactional.Transformations {
	if transformations == nil {
		return &transactional.Transformations{
			Timezone:       "UTC",
			TargetTimeZone: "UTC",
			Offset:         boolPtr(false),
		}
	}

	normalized := *transformations
	includeOffset := normalized.Offset != nil && *normalized.Offset
	hasAggregations := len(normalized.Keys) > 0 || len(normalized.Values) > 0
	if normalized.Offset == nil {
		normalized.Offset = boolPtr(false)
	}
	if normalized.TargetTimeZone == "" && normalized.Timezone != "" {
		normalized.TargetTimeZone = normalized.Timezone
	}
	if normalized.TargetTimeZone == "" && !includeOffset && !hasAggregations {
		normalized.TargetTimeZone = "UTC"
		normalized.Timezone = "UTC"
	}
	return &normalized
}

func transactionalRequestFromLiteQuery(r *http.Request) (transactional.Request, error) {
	query := r.URL.Query()

	rawID := strings.TrimSpace(query.Get("id"))
	if rawID == "" {
		return transactional.Request{}, fmt.Errorf("missing id")
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return transactional.Request{}, fmt.Errorf("invalid id")
	}

	from, fromTime, err := litePointInTime(query.Get("from"))
	if err != nil {
		return transactional.Request{}, fmt.Errorf("invalid from")
	}

	expressions := []string{"ReferenceTime >= " + from}
	if rawTo := strings.TrimSpace(query.Get("to")); rawTo != "" {
		to, toTime, err := litePointInTime(rawTo)
		if err != nil {
			return transactional.Request{}, fmt.Errorf("invalid to")
		}
		if !toTime.After(fromTime) {
			return transactional.Request{}, fmt.Errorf("to must be greater than from")
		}
		expressions = append(expressions, "ReferenceTime < "+to)
	}

	return transactional.Request{
		IDs: []domain.Identifier{domain.Identifier(id)},
		Filters: &transactional.Filters{
			Expressions: expressions,
		},
		Transformations: &transactional.Transformations{
			Timezone:       "UTC",
			TargetTimeZone: "UTC",
			Offset:         boolPtr(false),
		},
	}, nil
}

func boolPtr(value bool) *bool {
	return &value
}

func litePointInTime(value string) (string, time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", time.Time{}, fmt.Errorf("empty point in time")
	}

	layouts := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err != nil {
			continue
		}
		return value, parsed, nil
	}

	return "", time.Time{}, fmt.Errorf("unsupported point in time")
}
