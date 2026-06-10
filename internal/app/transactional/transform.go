package transactional

import (
	"context"
	"strings"
	"time"

	"streaming-golang/internal/domain"
)

type TransformationProcessor struct{}

func NewTransformationProcessor() *TransformationProcessor {
	return &TransformationProcessor{}
}

func (p *TransformationProcessor) Process(ctx context.Context, items []DataItem, command Command) []DataItem {
	location := time.UTC
	if command.TargetTimeZone != "" {
		if loc, err := time.LoadLocation(command.TargetTimeZone); err == nil {
			location = loc
		}
	}

	for i := range items {
		items[i] = p.processItem(items[i], command, location)
	}

	return items
}

func (p *TransformationProcessor) processItem(item DataItem, command Command, location *time.Location) DataItem {
	// 1. Timezone conversion
	if command.TargetTimeZone != "" {
		for key, value := range item.Fields {
			if t, ok := value.(time.Time); ok {
				item.Fields[key] = t.In(location)
			}
		}
	}

	// 2. RelativeDeliveryPeriod calculation (for Cassandra data)
	if command.Source == domain.SourceCassandra && containsColumn(command.Columns, "RelativeDeliveryPeriod") {
		p.calculateRDP(item, command)
	}

	return item
}

func (p *TransformationProcessor) calculateRDP(item DataItem, command Command) {
	refObj, ok1 := item.Fields["ReferenceTime"]
	delStartField := resolveDeliveryStartProperty(command.Mappings)
	delStartObj, ok2 := item.Fields[delStartField]

	if ok1 && ok2 {
		refTime, ok1 := refObj.(time.Time)
		delStart, ok2 := delStartObj.(time.Time)
		if ok1 && ok2 && len(command.Mappings) > 0 {
			mapping := command.Mappings[0]
			period := RDPCalculator{}.Calculate(refTime, delStart, mapping.Resolution, "") // Add delivery resolution if needed
			if period != nil {
				item.Fields["RelativeDeliveryPeriod"] = *period
			}
		}
	}
}

func containsColumn(columns []string, name string) bool {
	if len(columns) == 0 {
		return true // If columns are empty, it means all columns
	}
	for _, c := range columns {
		if strings.EqualFold(c, name) {
			return true
		}
	}
	return false
}

func resolveDeliveryStartProperty(mappings []domain.Mapping) string {
	for _, m := range mappings {
		for _, col := range m.Columns {
			switch strings.ToLower(col.MDSName) {
			case "deliverystart", "underlyingstart", "optionstart":
				return col.MDSName
			}
		}
	}
	return "DeliveryStart"
}
