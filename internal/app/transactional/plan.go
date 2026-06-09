package transactional

import (
	"context"

	"streaming-golang/internal/domain"
)

type Plan struct {
	Steps []PlanStep
}

type PlanStep struct {
	Command Command
	Queries []ExecutableQuery
}

type requestPlanner struct {
	mappings     MappingResolver
	quoteIndices QuoteIndexPlanner
	strategy     QueryStrategy
	queryBuilder QueryBuilder
}

type ExecutableQuery = domain.ExecutableQuery

type Mapping = domain.Mapping

type SourceKind = domain.SourceKind

func NewPlanner(options ...PlannerOption) Planner {
	planner := requestPlanner{
		mappings:     StaticMappingResolver{},
		quoteIndices: FilterQuoteIndexPlanner{},
		strategy:     SingleQueryStrategy{},
		queryBuilder: PlaceholderQueryBuilder{},
	}
	for _, option := range options {
		option(&planner)
	}
	return planner
}

type PlannerOption func(*requestPlanner)

func WithQueryStrategy(strategy QueryStrategy) PlannerOption {
	return func(p *requestPlanner) {
		if strategy != nil {
			p.strategy = strategy
		}
	}
}

func WithMappingResolver(resolver MappingResolver) PlannerOption {
	return func(p *requestPlanner) {
		if resolver != nil {
			p.mappings = resolver
		}
	}
}

func WithQuoteIndexPlanner(planner QuoteIndexPlanner) PlannerOption {
	return func(p *requestPlanner) {
		if planner != nil {
			p.quoteIndices = planner
		}
	}
}

func WithQueryBuilder(builder QueryBuilder) PlannerOption {
	return func(p *requestPlanner) {
		if builder != nil {
			p.queryBuilder = builder
		}
	}
}

func NewPlannerWithStrategy(strategies ...QueryStrategy) Planner {
	strategy := QueryStrategy(SingleQueryStrategy{})
	if len(strategies) > 0 && strategies[0] != nil {
		strategy = strategies[0]
	}
	return NewPlanner(WithQueryStrategy(strategy))
}

func (p requestPlanner) BuildPlan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, error) {
	steps := make([]PlanStep, 0, len(requests))

	for _, request := range requests {
		mappings, err := p.mappings.ResolveMappings(ctx, request.IDs, requestContext.DataCategory, requestContext.Stage)
		if err != nil {
			return Plan{}, err
		}

		command := Command{
			IDs:             request.IDs,
			DataCategory:    requestContext.DataCategory,
			Columns:         append([]string(nil), request.Columns...),
			IncludeOffset:   includeOffset(requestContext, request),
			TargetTimeZone:  targetTimeZone(request),
			HasAggregations: hasAggregations(request),
			HasShape:        request.Filters != nil && len(request.Filters.Shape) > 0,
			Mappings:        mappings,
		}
		if request.Filters != nil {
			command.Filters = request.Filters.Parsed
			if len(command.Filters.Expressions) == 0 {
				command.Filters.Expressions = request.Filters.Expressions
			}
		}
		if err := validateAgainstMappings(requestContext, request, command, mappings); err != nil {
			return Plan{}, err
		}
		command.Source = sourceFromMappings(mappings)
		command.QuoteIndices, err = p.quoteIndices.PlanQuoteIndices(ctx, command)
		if err != nil {
			return Plan{}, err
		}

		splitCommands := p.strategy.Plan(command)
		for _, splitCommand := range splitCommands {
			built, err := p.queryBuilder.BuildQueries(ctx, splitCommand)
			if err != nil {
				return Plan{}, err
			}
			steps = append(steps, PlanStep{
				Command: splitCommand,
				Queries: built,
			})
		}
	}

	return Plan{Steps: steps}, nil
}

func includeOffset(requestContext RequestContext, request Request) bool {
	if requestContext.Mode == ModeJSON || requestContext.Mode == ModeJSONStream || requestContext.Mode == ModeNDJSONStream {
		return true
	}
	return request.Transformations != nil && request.Transformations.Offset != nil && *request.Transformations.Offset
}

func targetTimeZone(request Request) string {
	if request.Transformations == nil {
		return ""
	}
	if request.Transformations.TargetTimeZone != "" {
		return request.Transformations.TargetTimeZone
	}
	return request.Transformations.Timezone
}

func hasAggregations(request Request) bool {
	return request.Transformations != nil && (len(request.Transformations.Keys) > 0 || len(request.Transformations.Values) > 0)
}

func sourceFromMappings(mappings []Mapping) SourceKind {
	if len(mappings) == 0 {
		return domain.SourceCMDP
	}
	return mappings[0].Source
}
