package transactional

import (
	"context"
	"time"
)

type TransformationProcessor struct{}

func NewTransformationProcessor() *TransformationProcessor {
	return &TransformationProcessor{}
}

func (p *TransformationProcessor) Process(ctx context.Context, items []DataItem, command Command) []DataItem {
	if command.TargetTimeZone == "" && !command.HasAggregations {
		return items
	}

	location, err := time.LoadLocation(command.TargetTimeZone)
	if err != nil {
		// If timezone is invalid, we don't convert, but it should have been validated already
		location = time.UTC
	}

	for i := range items {
		items[i] = p.processItem(items[i], command, location)
	}

	return items
}

func (p *TransformationProcessor) processItem(item DataItem, command Command, location *time.Location) DataItem {
	if command.TargetTimeZone != "" {
		if location == nil {
			loaded, err := time.LoadLocation(command.TargetTimeZone)
			if err != nil {
				loaded = time.UTC
			}
			location = loaded
		}
		for key, value := range item.Fields {
			if t, ok := value.(time.Time); ok {
				item.Fields[key] = t.In(location)
			}
		}
	}

	// TODO: Implement aggregations, nesting, etc.

	return item
}
