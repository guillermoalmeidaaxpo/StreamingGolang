package httpapi

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"streaming-golang/internal/app/transactional"
)

const transactionalCSVFilename = `attachment; filename="transactional_data.csv"`

func writeTransactionalCSV(w http.ResponseWriter, response transactional.Response, columns []string, includeOffset bool, attachment bool) error {
	setCSVHeaders(w, attachment)
	sortCSVItems(response.TransactionalData)
	if len(columns) == 0 {
		columns = csvColumns(response.TransactionalData)
	}
	if len(columns) == 0 {
		return nil
	}

	writer := csv.NewWriter(w)
	if err := writer.Write(columns); err != nil {
		return err
	}
	for _, item := range response.TransactionalData {
		if err := writer.Write(csvRow(columns, item, includeOffset)); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func writeTransactionalCSVStream(ctx context.Context, w http.ResponseWriter, stream transactional.Stream, columns []string, includeOffset bool, attachment bool) error {
	setCSVHeaders(w, attachment)
	if len(columns) == 0 {
		return nil
	}

	writer := csv.NewWriter(w)
	flusher, _ := w.(http.Flusher)
	if err := writer.Write(columns); err != nil {
		return err
	}

	for stream.Next(ctx) {
		item := stream.Item()
		if err := writer.Write(csvRow(columns, item, includeOffset)); err != nil {
			return err
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			return err
		}
		if flusher != nil {
			flusher.Flush()
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return stream.Err()
}

func setCSVHeaders(w http.ResponseWriter, attachment bool) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	if attachment {
		w.Header().Set("Content-Disposition", transactionalCSVFilename)
	}
}

func csvIncludeOffset(plan transactional.Plan) bool {
	return len(plan.Steps) > 0 && plan.Steps[0].Command.IncludeOffset
}

func csvColumns(items []transactional.DataItem) []string {
	fields := make(map[string]struct{})
	for _, item := range items {
		for key := range item.Fields {
			fields[key] = struct{}{}
		}
	}
	return csvColumnsFromFields(fields)
}

func csvColumnsFromPlan(plan transactional.Plan) []string {
	columns := make([]string, 0)
	seen := make(map[string]struct{})
	for _, step := range plan.Steps {
		command := step.Command
		if command.HasAggregations && len(command.Columns) > 0 {
			for _, column := range command.Columns {
				columns = appendCSVColumn(columns, seen, column)
			}
			continue
		}
		requested := requestedCSVColumns(command.Columns)
		hasHyperscale := false
		beforeCommandColumns := len(columns)

		for _, mapping := range command.Mappings {
			if mapping.HyperscaleID != nil {
				hasHyperscale = true
			}
			for _, column := range mapping.Columns {
				name := strings.TrimSpace(column.MDSName)
				if name == "" {
					name = strings.TrimSpace(column.SourceName)
				}
				if name == "" {
					continue
				}
				if len(requested) > 0 {
					if _, ok := requested[strings.ToLower(name)]; !ok {
						continue
					}
				} else if !column.IsKey && !column.IsProjectable {
					continue
				}
				columns = appendCSVColumn(columns, seen, name)
			}
		}

		if len(command.Mappings) == 0 {
			for _, column := range command.Columns {
				columns = appendCSVColumn(columns, seen, column)
			}
		}
		if len(columns) == beforeCommandColumns && len(command.Columns) > 0 {
			for _, column := range command.Columns {
				columns = appendCSVColumn(columns, seen, column)
			}
		}
		if hasHyperscale {
			if _, ok := requested["createdon"]; ok {
				columns = appendCSVColumn(columns, seen, "CreatedOn")
			}
		}
	}

	if len(columns) > 0 {
		return columns
	}
	return csvColumnsFromQueries(plan)
}

func csvColumnsFromFields(fields map[string]struct{}) []string {
	columns := make([]string, 0, len(fields))
	preferred := []string{"status", "source", "dataCategory", "statement", "parameterCount", "indexStart", "indexEnd"}
	for _, column := range preferred {
		if _, ok := fields[column]; ok {
			columns = append(columns, column)
			delete(fields, column)
		}
	}

	remaining := make([]string, 0, len(fields))
	for column := range fields {
		remaining = append(remaining, column)
	}
	sort.Strings(remaining)
	return append(columns, remaining...)
}

func csvColumnsFromQueries(plan transactional.Plan) []string {
	fields := map[string]struct{}{
		"status":         {},
		"source":         {},
		"dataCategory":   {},
		"statement":      {},
		"parameterCount": {},
	}

	for _, step := range plan.Steps {
		for _, query := range step.Queries {
			if query.IndexRange != nil {
				fields["indexStart"] = struct{}{}
				fields["indexEnd"] = struct{}{}
			}
		}
	}
	return csvColumnsFromFields(fields)
}

func requestedCSVColumns(columns []string) map[string]struct{} {
	if len(columns) == 0 {
		return nil
	}
	requested := make(map[string]struct{}, len(columns))
	for _, column := range columns {
		column = strings.ToLower(strings.TrimSpace(column))
		if column != "" {
			requested[column] = struct{}{}
		}
	}
	return requested
}

func appendCSVColumn(columns []string, seen map[string]struct{}, column string) []string {
	column = strings.TrimSpace(column)
	if column == "" {
		return columns
	}
	key := strings.ToLower(column)
	if _, exists := seen[key]; exists {
		return columns
	}
	seen[key] = struct{}{}
	return append(columns, column)
}

func sortCSVItems(items []transactional.DataItem) {
	sort.SliceStable(items, func(i, j int) bool {
		leftReference, leftHasReference := csvTimeField(items[i].Fields, "ReferenceTime")
		rightReference, rightHasReference := csvTimeField(items[j].Fields, "ReferenceTime")
		if !leftHasReference || !rightHasReference {
			return false
		}
		if !leftReference.Equal(rightReference) {
			return leftReference.Before(rightReference)
		}

		leftDelivery, leftHasDelivery := csvTimeField(items[i].Fields, "DeliveryStart")
		rightDelivery, rightHasDelivery := csvTimeField(items[j].Fields, "DeliveryStart")
		if !leftHasDelivery || !rightHasDelivery {
			return false
		}
		return leftDelivery.Before(rightDelivery)
	})
}

func csvTimeField(fields map[string]any, name string) (time.Time, bool) {
	value, ok := lookupCSVField(fields, name)
	if !ok {
		return time.Time{}, false
	}
	switch typed := value.(type) {
	case time.Time:
		return typed, true
	case string:
		parsed, err := time.Parse(time.RFC3339Nano, typed)
		if err == nil {
			return parsed, true
		}
		for _, layout := range []string{"2006-01-02T15:04:05.000", "2006-01-02T15:04:05"} {
			parsed, err := time.ParseInLocation(layout, typed, time.UTC)
			if err == nil {
				return parsed, true
			}
		}
	}
	return time.Time{}, false
}

func lookupCSVField(fields map[string]any, column string) (any, bool) {
	if fields == nil {
		return nil, false
	}
	if value, ok := fields[column]; ok {
		return value, true
	}
	for key, value := range fields {
		if strings.EqualFold(key, column) {
			return value, true
		}
	}
	return nil, false
}

func formatCSVValue(value any, includeOffset bool) string {
	if value == nil {
		return "null"
	}

	switch typed := value.(type) {
	case time.Time:
		if includeOffset {
			return typed.Format("2006-01-02T15:04:05.000-07:00")
		}
		return typed.Format("2006-01-02T15:04:05.000")
	case float32:
		return strconv.FormatFloat(float64(typed), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int8:
		return strconv.FormatInt(int64(typed), 10)
	case int16:
		return strconv.FormatInt(int64(typed), 10)
	case int32:
		return strconv.FormatInt(int64(typed), 10)
	case int64:
		return strconv.FormatInt(typed, 10)
	case uint:
		return strconv.FormatUint(uint64(typed), 10)
	case uint8:
		return strconv.FormatUint(uint64(typed), 10)
	case uint16:
		return strconv.FormatUint(uint64(typed), 10)
	case uint32:
		return strconv.FormatUint(uint64(typed), 10)
	case uint64:
		return strconv.FormatUint(typed, 10)
	default:
		return fmt.Sprint(value)
	}
}

func csvRow(columns []string, item transactional.DataItem, includeOffset bool) []string {
	row := make([]string, len(columns))
	for index, column := range columns {
		value, ok := lookupCSVField(item.Fields, column)
		if !ok {
			row[index] = "N/A"
			continue
		}
		row[index] = formatCSVValue(value, includeOffset)
	}
	return row
}
