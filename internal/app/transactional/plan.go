package transactional

import (
	"context"
	"log/slog"
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
	logger       *slog.Logger
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
		logger:       slog.Default(),
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

func WithLogger(logger *slog.Logger) PlannerOption {
	return func(p *requestPlanner) {
		if logger != nil {
			p.logger = logger
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
			command = selectFetchingStrategy(command)
			command, err = p.normalizeDefaultFilters(ctx, command)
			if err != nil {
				return Plan{}, err
			}
			if err := validateAgainstMappings(requestContext, request, command, command.Mappings); err != nil {
				return Plan{}, err
			}

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
					p.logQuoteIndices(ctx, "cassandra quote indices generated", hCommand)
					if len(hCommand.QuoteIndices) == 0 && hasReferenceTimeFilter(hCommand.Filters.Nodes) {
						p.logger.WarnContext(ctx, "skipping Cassandra command because quote indices are empty",
							slog.Any("identifiers", hCommand.IDs),
							slog.String("data_category", string(hCommand.DataCategory)),
							slog.String("filter_time_zone", hCommand.FilterTimeZone),
						)
						continue
					}
				} else {
					commandQuoteIndices, err := p.quoteIndices.PlanQuoteIndices(ctx, hCommand)
					if err != nil {
						return Plan{}, err
					}
					hCommand.QuoteIndices = commandQuoteIndices
					p.logQuoteIndices(ctx, "CMDP quote indices generated", hCommand)
				}

				splitCommands := p.strategy.Plan(hCommand)
				for _, splitCommand := range splitCommands {
					built, err := p.queryBuilder.BuildQueries(ctx, splitCommand)
					if err != nil {
						return Plan{}, err
					}
					if len(built) == 0 {
						p.logger.WarnContext(ctx, "query builder produced no queries",
							slog.Any("identifiers", splitCommand.IDs),
							slog.String("source", string(splitCommand.Source)),
							slog.String("data_category", string(splitCommand.DataCategory)),
							slog.Int("mapping_count", len(splitCommand.Mappings)),
							slog.Int("quote_index_count", len(splitCommand.QuoteIndices)),
						)
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

func (p requestPlanner) normalizeDefaultFilters(ctx context.Context, command Command) (Command, error) {
	if command.Source == domain.SourceHyperscale && len(command.Filters.Nodes) == 0 {
		command.LatestReferenceTime = true
	}

	if !hasLatestGlobalFilter(command.Filters.Nodes) {
		return command, nil
	}

	if command.Source == domain.SourceHyperscale {
		command.Filters.Nodes = removeLatestGlobalFilters(command.Filters.Nodes)
		command.LatestReferenceTime = true
		p.logger.InfoContext(ctx, "latestGlobal filter resolved by hyperscale latest view",
			slog.Any("identifiers", command.IDs),
			slog.String("source", string(command.Source)),
			slog.String("data_category", string(command.DataCategory)),
		)
		return command, nil
	}

	watermark, err := p.mappings.GetWatermark(ctx, command.Mappings)
	if err != nil {
		return Command{}, err
	}
	command.Filters.Nodes = replaceLatestGlobalFilters(command.Filters.Nodes, watermark)
	p.logger.InfoContext(ctx, "latestGlobal filter resolved",
		slog.Any("identifiers", command.IDs),
		slog.String("source", string(command.Source)),
		slog.String("data_category", string(command.DataCategory)),
		slog.Time("reference_time", watermark),
	)
	return command, nil
}

func (p requestPlanner) splitHybridCommand(ctx context.Context, command Command) ([]Command, error) {
	// Only split if the command is eligible
	if command.Source != domain.SourceCassandra || command.HasAggregations || command.HasShape || !isEligibleForHybridSplit(command.Mappings) || !hasReferenceTimeFilter(command.Filters.Nodes) {
		return []Command{command}, nil
	}

	filterLocation, _ := loadLocation(command.FilterTimeZone)
	watermark, err := p.mappings.GetWatermark(ctx, command.Mappings)
	if err != nil {
		return nil, err
	}

	if point, ok, err := referenceTimeEqualityPoint(command.Filters.Nodes, filterLocation); err != nil {
		return nil, err
	} else if ok {
		selectedSource := domain.SourceCMDP
		if point.Before(watermark) {
			selectedSource = domain.SourceCassandra
		}
		p.logger.InfoContext(ctx, "hybrid equality route selected",
			slog.Any("identifiers", command.IDs),
			slog.Time("reference_time", point),
			slog.Time("watermark", watermark),
			slog.String("source", string(selectedSource)),
			slog.String("filter_time_zone", command.FilterTimeZone),
		)
		return []Command{withCommandSource(command, selectedSource)}, nil
	}

	// Analyze ReferenceTime filters to decide how to split
	location, _ := loadCassandraLocation(cassandraTimeZone(command.Mappings))
	limits, err := cassandraReferenceTimeRange(command.Filters.Nodes, location, time.Now().UTC(), filterLocation)
	if err != nil {
		return nil, err
	}

	// 1. Entirely Cassandra: UpperLimit < watermark
	if limits.end != nil && limits.end.Before(watermark) {
		p.logger.InfoContext(ctx, "hybrid range route selected",
			slog.Any("identifiers", command.IDs),
			slog.Time("range_end", *limits.end),
			slog.Time("watermark", watermark),
			slog.String("source", string(domain.SourceCassandra)),
		)
		return []Command{withCommandSource(command, domain.SourceCassandra)}, nil
	}

	// 2. Entirely CMDP: LowerLimit >= watermark
	if limits.start != nil && (limits.start.After(watermark) || limits.start.Equal(watermark)) {
		p.logger.InfoContext(ctx, "hybrid range route selected",
			slog.Any("identifiers", command.IDs),
			slog.Time("range_start", *limits.start),
			slog.Time("watermark", watermark),
			slog.String("source", string(domain.SourceCMDP)),
		)
		return []Command{withCommandSource(command, domain.SourceCMDP)}, nil
	}

	// 3. Hybrid: Crossing the watermark
	// Create a Cassandra command for [min, watermark)
	cassandraCmd := withCommandSource(command, domain.SourceCassandra)
	cassandraCmd.Filters = command.Filters.Clone()
	cassandraCmd.Filters.Nodes = append(cassandraCmd.Filters.Nodes, domain.ComparisonFilter{
		Field:    referenceTimeField,
		Operator: "<",
		Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: watermark.Format(time.RFC3339Nano)},
	})

	// Create a CMDP command for [watermark, max]
	cmdpCmd := withCommandSource(command, domain.SourceCMDP)
	cmdpCmd.Filters = command.Filters.Clone()
	cmdpCmd.Filters.Nodes = append(cmdpCmd.Filters.Nodes, domain.ComparisonFilter{
		Field:    referenceTimeField,
		Operator: ">=",
		Value:    domain.FilterValue{Kind: domain.FilterValuePointInTime, Raw: watermark.Format(time.RFC3339Nano)},
	})

	p.logger.InfoContext(ctx, "hybrid range split selected",
		slog.Any("identifiers", command.IDs),
		slog.Time("watermark", watermark),
	)
	return []Command{cassandraCmd, cmdpCmd}, nil
}

func hasLatestGlobalFilter(nodes []domain.FilterNode) bool {
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if ok && filter.Value.Kind == domain.FilterValueLatestGlobal {
			return true
		}
	}
	return false
}

func removeLatestGlobalFilters(nodes []domain.FilterNode) []domain.FilterNode {
	result := make([]domain.FilterNode, 0, len(nodes))
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if ok && filter.Value.Kind == domain.FilterValueLatestGlobal {
			continue
		}
		result = append(result, node)
	}
	return result
}

func replaceLatestGlobalFilters(nodes []domain.FilterNode, referenceTime time.Time) []domain.FilterNode {
	result := make([]domain.FilterNode, 0, len(nodes))
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if !ok || filter.Value.Kind != domain.FilterValueLatestGlobal {
			result = append(result, node)
			continue
		}
		filter.Value = domain.FilterValue{
			Kind: domain.FilterValuePointInTime,
			Raw:  referenceTime.UTC().Format(time.RFC3339Nano),
		}
		result = append(result, filter)
	}
	return result
}

func withCommandSource(command Command, source domain.SourceKind) Command {
	command.Source = source
	command.Mappings = append([]domain.Mapping(nil), command.Mappings...)
	for i := range command.Mappings {
		command.Mappings[i].Source = source
	}
	return command
}

func (p requestPlanner) logQuoteIndices(ctx context.Context, message string, command Command) {
	count := len(command.QuoteIndices)
	attrs := []any{
		slog.Any("identifiers", command.IDs),
		slog.String("source", string(command.Source)),
		slog.String("data_category", string(command.DataCategory)),
		slog.Int("quote_index_count", count),
	}
	if count > 0 {
		attrs = append(attrs,
			slog.Int("first_quote_index", command.QuoteIndices[0]),
			slog.Int("last_quote_index", command.QuoteIndices[count-1]),
		)
		if count <= 10 {
			attrs = append(attrs, slog.Any("quote_indices", command.QuoteIndices))
		}
	}
	p.logger.InfoContext(ctx, message, attrs...)
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

func referenceTimeEqualityPoint(nodes []domain.FilterNode, loc *time.Location) (time.Time, bool, error) {
	for _, node := range nodes {
		filter, ok := node.(domain.ComparisonFilter)
		if !ok || !strings.EqualFold(filter.Field, referenceTimeField) || filter.Operator != "=" {
			continue
		}
		point, ok, err := pointTime(filter.Value, loc)
		return point, ok, err
	}
	return time.Time{}, false, nil
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
		Columns:           projectionColumns(request.Columns, mappings, includeIdentifier(requestContext)),
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

func projectionColumns(requested []string, mappings []Mapping, isCSVEndpoint bool) []string {
	if len(mappings) == 0 {
		return append([]string(nil), requested...)
	}

	if len(requested) == 0 {
		return mappedColumnNames(mappings, isCSVEndpoint, false)
	}

	if len(requested) == 1 && strings.EqualFold(strings.TrimSpace(requested[0]), "CreatedOn") {
		return mappedColumnNames(mappings, isCSVEndpoint, true)
	}

	available := make(map[string]string)
	for _, mapping := range mappings {
		for _, column := range mapping.Columns {
			if strings.TrimSpace(column.MDSName) != "" {
				available[strings.ToLower(column.MDSName)] = column.MDSName
			}
		}
	}
	available["createdon"] = "CreatedOn"

	seen := make(map[string]struct{})
	columns := make([]string, 0)
	for _, mapping := range mappings {
		for _, column := range mapping.Columns {
			if column.IsProjectable {
				continue
			}
			appendProjectionColumn(&columns, seen, column.MDSName, isCSVEndpoint)
		}
	}
	for _, requestedColumn := range requested {
		if mapped, ok := available[strings.ToLower(strings.TrimSpace(requestedColumn))]; ok {
			appendProjectionColumn(&columns, seen, mapped, isCSVEndpoint)
		}
	}
	return columns
}

func mappedColumnNames(mappings []Mapping, isCSVEndpoint bool, includeCreatedOn bool) []string {
	seen := make(map[string]struct{})
	columns := make([]string, 0)
	for _, mapping := range mappings {
		for _, column := range mapping.Columns {
			appendProjectionColumn(&columns, seen, column.MDSName, isCSVEndpoint)
		}
	}
	if includeCreatedOn {
		appendProjectionColumn(&columns, seen, "CreatedOn", isCSVEndpoint)
	}
	return columns
}

func appendProjectionColumn(columns *[]string, seen map[string]struct{}, name string, isCSVEndpoint bool) {
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	if !isCSVEndpoint && (strings.EqualFold(name, "Identifier") || strings.EqualFold(name, "MdoId")) {
		return
	}
	key := strings.ToLower(name)
	if _, exists := seen[key]; exists {
		return
	}
	seen[key] = struct{}{}
	*columns = append(*columns, name)
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

func selectFetchingStrategy(command Command) Command {
	return withCommandSource(command, fetchingSource(command))
}

func fetchingSource(command Command) SourceKind {
	if len(command.Mappings) == 0 {
		return domain.SourceCMDP
	}

	mapping := command.Mappings[0]
	if mapping.Source == domain.SourceMesap {
		return domain.SourceMesap
	}
	if mapping.HyperscaleID != nil || mapping.Source == domain.SourceHyperscale {
		return domain.SourceHyperscale
	}
	if shouldUseCMDPStrategy(command) {
		return domain.SourceCMDP
	}
	return domain.SourceCassandra
}

func shouldUseCMDPStrategy(command Command) bool {
	if command.HasShape || command.HasAggregations {
		return true
	}
	if len(command.Mappings) == 0 {
		return true
	}

	mapping := command.Mappings[0]
	if mapping.Source == domain.SourceCMDP {
		return true
	}
	if _, ok := hpfcIDsForcedToCMDP[mapping.ID]; ok {
		return true
	}
	if strings.TrimSpace(mapping.CassandraID) == "" {
		return true
	}
	return !isEuropeZurichCassandraTimezone(cassandraTimeZone(command.Mappings))
}

func isEuropeZurichCassandraTimezone(timezone string) bool {
	switch strings.ToLower(strings.TrimSpace(timezone)) {
	case "", "cet", "europe/zurich":
		return true
	default:
		return false
	}
}

var hpfcIDsForcedToCMDP = map[domain.Identifier]struct{}{
	536000751: {},
	536214287: {},
	536346251: {},
}
