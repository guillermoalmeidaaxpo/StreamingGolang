package transactional

import (
	"context"
	"strings"
	"time"

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

			hybridCommands, err := p.splitHybridCommand(ctx, command)
			if err != nil {
				return Plan{}, err
			}

			for _, hCommand := range hybridCommands {
				if hCommand.Source == domain.SourceCassandra {
					hCommand.QuoteIndices, err = CassandraQuoteIndexPlanner{}.PlanQuoteIndices(ctx, hCommand)
					if err != nil {
						return Plan{}, err
					}
				} else {
					commandQuoteIndices, err := p.quoteIndices.PlanQuoteIndices(ctx, hCommand)
					if err != nil {
						return Plan{}, err
					}
					hCommand.QuoteIndices = commandQuoteIndices
				}

				splitCommands := p.strategy.Plan(hCommand)
				for _, splitCommand := range splitCommands {
					built, err := p.queryBuilder.BuildQueries(ctx, splitCommand)
					if err != nil {
						return Plan{}, err
					}
					if len(built) == 0 {
						continue
					}
					steps = append(steps, PlanStep{
						Command: splitCommand,
						Queries: built,
					})
				}
			}
		}
	}

	return Plan{Steps: steps}, nil
}

func (p requestPlanner) splitHybridCommand(ctx context.Context, command Command) ([]Command, error) {
	// Only split if the command is eligible
	if command.HasAggregations || !isEligibleForHybridSplit(command.Mappings) || !hasReferenceTimeFilter(command.Filters.Nodes) {
		return []Command{command}, nil
	}

	watermark, err := p.mappings.GetWatermark(ctx, command.Mappings)
	if err != nil {
		return nil, err
	}

	// Analyze ReferenceTime filters to decide how to split
	location, _ := loadCassandraLocation(cassandraTimeZone(command.Mappings))
	limits, err := cassandraReferenceTimeRange(command.Filters.Nodes, location, time.Now().UTC(), location)
	if err != nil {
		return nil, err
	}

	// 1. Entirely Cassandra: UpperLimit < watermark
	if limits.end != nil && limits.end.Before(watermark) {
		command.Source = domain.SourceCassandra
		return []Command{command}, nil
	}

	// 2. Entirely CMDP: LowerLimit >= watermark
	if limits.start != nil && (limits.start.After(watermark) || limits.start.Equal(watermark)) {
		command.Source = domain.SourceCMDP
		return []Command{command}, nil
	}

	// 3. Hybrid: Crossing the watermark
	// Create a Cassandra command for [min, watermark)
	cassandraCmd := command
	cassandraCmd.Source = domain.SourceCassandra
	cassandraCmd.Filters = command.Filters.Clone()
	cassandraCmd.Filters.Nodes = append(cassandraCmd.Filters.Nodes, domain.ComparisonFilter{
		Field:    referenceTimeField,
		Operator: "<",
		Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: watermark.Format(time.RFC3339Nano)},
	})

	// Create a CMDP command for [watermark, max]
	cmdpCmd := command
	cmdpCmd.Source = domain.SourceCMDP
	cmdpCmd.Filters = command.Filters.Clone()
	cmdpCmd.Filters.Nodes = append(cmdpCmd.Filters.Nodes, domain.ComparisonFilter{
		Field:    referenceTimeField,
		Operator: ">=",
		Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: watermark.Format(time.RFC3339Nano)},
	})

	return []Command{cassandraCmd, cmdpCmd}, nil
}

func isEligibleForHybridSplit(mappings []Mapping) bool {
	if len(mappings) == 0 {
		return false
	}
	for _, m := range mappings {
		// Only split if explicitly allowed AND it has a Cassandra ID
		// (if it has no Cassandra ID, there is nothing to split to)
		if !m.SplitQuery || strings.TrimSpace(m.CassandraID) == "" {
			return false
		}
	}
	return true
}

func hasReferenceTimeFilter(nodes []domain.FilterNode) bool {
	for _, node := range nodes {
		if filter, ok := node.(domain.ComparisonFilter); ok && strings.EqualFold(filter.Field, referenceTimeField) {
			return true
		}
	}
	return false
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
