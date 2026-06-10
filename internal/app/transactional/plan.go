package transactional

import (
	"context"
	"fmt"

	"streaming-golang/internal/app/apperr"
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

		commands := commandsForRequest(requestContext, request, mappings)
		for _, command := range commands {
			if err := validateAgainstMappings(requestContext, request, command, command.Mappings); err != nil {
				return Plan{}, err
			}
			command.Source = sourceFromMappings(command.Mappings)
			if command.Source == domain.SourceCassandra {
				command.QuoteIndices, err = CassandraQuoteIndexPlanner{}.PlanQuoteIndices(ctx, command)
				if err != nil {
					return Plan{}, err
				}
			} else {
				command.QuoteIndices, err = p.quoteIndices.PlanQuoteIndices(ctx, command)
				if err != nil {
					return Plan{}, err
				}
			}

			splitCommands := p.strategy.Plan(command)
			for _, splitCommand := range splitCommands {
				built, err := p.queryBuilder.BuildQueries(ctx, splitCommand)
				if err != nil {
					return Plan{}, err
				}
				if len(built) == 0 {
					return Plan{}, apperr.New(apperr.Unavailable, fmt.Sprintf("no query builder produced a query for source %q", splitCommand.Source))
				}
				steps = append(steps, PlanStep{
					Command: splitCommand,
					Queries: built,
				})
			}
		}
	}

	return Plan{Steps: steps}, nil
}

func commandsForRequest(requestContext RequestContext, request Request, mappings []Mapping) []Command {
	if endpointKind(requestContext) != EndpointGeneric || requestContext.DataCategory != "" || len(mappings) == 0 {
		return []Command{newCommand(requestContext, request, requestContext.DataCategory, request.IDs, mappings)}
	}

	grouped := make(map[domain.DataCategory][]Mapping)
	order := make([]domain.DataCategory, 0)
	for _, mapping := range mappings {
		category := mapping.DataCategory
		if _, exists := grouped[category]; !exists {
			order = append(order, category)
		}
		grouped[category] = append(grouped[category], mapping)
	}

	commands := make([]Command, 0, len(order))
	for _, category := range order {
		group := grouped[category]
		ids := make([]domain.Identifier, 0, len(group))
		for _, mapping := range group {
			ids = append(ids, mapping.ID)
		}
		commands = append(commands, newCommand(requestContext, request, category, ids, group))
	}
	return commands
}

func newCommand(requestContext RequestContext, request Request, category domain.DataCategory, ids []domain.Identifier, mappings []Mapping) Command {
	command := Command{
		IDs:               append([]domain.Identifier(nil), ids...),
		DataCategory:      category,
		Columns:           append([]string(nil), request.Columns...),
		VersionAsOf:       request.VersionAsOf,
		IncludeDeleted:    includeDeleted(request),
		IncludeIdentifier: includeIdentifier(requestContext),
		IncludeOffset:     includeOffset(requestContext, request),
		FilterTimeZone:    filterTimeZone(request),
		TargetTimeZone:    targetTimeZone(request),
		HasAggregations:   hasAggregations(request),
		HasShape:          request.Filters != nil && len(request.Filters.Shape) > 0,
		Mappings:          append([]Mapping(nil), mappings...),
	}
	if request.Filters != nil {
		command.Filters = request.Filters.Parsed
		if len(command.Filters.Expressions) == 0 {
			command.Filters.Expressions = request.Filters.Expressions
		}
	}
	return command
}

func filterTimeZone(request Request) string {
	if request.Filters == nil {
		return ""
	}
	return request.Filters.FilterTimeZone
}

func includeDeleted(request Request) bool {
	return request.IncludeDeleted != nil && *request.IncludeDeleted
}

func includeIdentifier(requestContext RequestContext) bool {
	return requestContext.Mode == ModeCSV || requestContext.Mode == ModeCSVStream
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
