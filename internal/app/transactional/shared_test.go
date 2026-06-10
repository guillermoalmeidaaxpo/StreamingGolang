package transactional

import (
	"context"
	"time"

	"streaming-golang/internal/domain"
)

type fixedMappingResolver struct {
	mappings []domain.Mapping
}

func (r fixedMappingResolver) ResolveMappings(context.Context, []domain.Identifier, domain.DataCategory, string) ([]domain.Mapping, error) {
	return r.mappings, nil
}

func (r fixedMappingResolver) GetWatermark(context.Context, []domain.Mapping) (time.Time, error) {
	return time.Now().UTC(), nil
}

func mappingWithColumns(id domain.Identifier, category domain.DataCategory) domain.Mapping {
	return domain.Mapping{
		ID:           id,
		DataCategory: category,
		Source:       domain.SourceCMDP,
		ViewName:     "CurveView",
		IndexField:   "QuoteDateIndex_FID",
		Columns: []domain.ColumnMapping{
			{MDSName: "ReferenceTime", SourceName: "ReferenceTime"},
			{MDSName: "Value", SourceName: "Value"},
		},
	}
}

func referenceTimeInterval(start, end string) domain.ComparisonFilter {
	return domain.ComparisonFilter{
		Field:    "ReferenceTime",
		Operator: "in",
		Value: domain.FilterValue{
			Kind:  domain.FilterValueTimeInterval,
			Raw:   "ti(" + start + "," + end + ")",
			Start: start,
			End:   end,
		},
	}
}

func referenceTimePoint(operator, raw string) domain.ComparisonFilter {
	return domain.ComparisonFilter{
		Field:    "ReferenceTime",
		Operator: operator,
		Value: domain.FilterValue{
			Kind: domain.FilterValuePointInTime,
			Raw:  raw,
		},
	}
}
