package transactional

import (
	"context"
	"log/slog"
	"time"
)

type Validator interface {
	Validate(context.Context, []Request) error
}

type FilterParser interface {
	Parse(context.Context, []string, string) (FilterSet, error)
}

type Planner interface {
	BuildPlan(context.Context, RequestContext, []Request) (Plan, error)
}

type Executor interface {
	Execute(context.Context, Plan) (Response, error)
	Stream(context.Context, Plan) (Stream, error)
}

type Pipeline struct {
	validator Validator
	parser    FilterParser
	planner   Planner
	executor  Executor
	logger    *slog.Logger
}

func NewPipeline(validator Validator, parser FilterParser, planner Planner, executor Executor) *Pipeline {
	return &Pipeline{
		validator: validator,
		parser:    parser,
		planner:   planner,
		executor:  executor,
		logger:    slog.Default(),
	}
}

func (p *Pipeline) WithLogger(logger *slog.Logger) *Pipeline {
	if logger != nil {
		p.logger = logger
	}
	return p
}

func (p *Pipeline) Execute(ctx context.Context, requestContext RequestContext, requests []Request) (Response, error) {
	_, response, err := p.ExecuteWithPlan(ctx, requestContext, requests)
	return response, err
}

func (p *Pipeline) ExecuteWithPlan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, Response, error) {
	start := time.Now()
	plan, err := p.Plan(ctx, requestContext, requests)
	if err != nil {
		return Plan{}, Response{}, err
	}

	executeStart := time.Now()
	response, err := p.executor.Execute(ctx, plan)
	if err != nil {
		return Plan{}, Response{}, err
	}
	p.logger.InfoContext(ctx, "pipeline execute completed",
		slog.String("mode", string(requestContext.Mode)),
		slog.String("endpoint_kind", string(requestContext.EndpointKind)),
		slog.String("data_category", string(requestContext.DataCategory)),
		slog.Int("request_count", len(requests)),
		slog.Int("plan_steps", len(plan.Steps)),
		slog.Int("row_count", len(response.TransactionalData)),
		slog.Duration("execute_duration", time.Since(executeStart)),
		slog.Int64("execute_duration_ms", time.Since(executeStart).Milliseconds()),
		slog.Duration("total_duration", time.Since(start)),
		slog.Int64("total_duration_ms", time.Since(start).Milliseconds()),
	)
	return plan, response, nil
}

func (p *Pipeline) Stream(ctx context.Context, requestContext RequestContext, requests []Request) (Stream, error) {
	_, stream, err := p.StreamWithPlan(ctx, requestContext, requests)
	return stream, err
}

func (p *Pipeline) StreamWithPlan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, Stream, error) {
	start := time.Now()
	plan, err := p.Plan(ctx, requestContext, requests)
	if err != nil {
		return Plan{}, nil, err
	}

	streamStart := time.Now()
	stream, err := p.executor.Stream(ctx, plan)
	if err != nil {
		return Plan{}, nil, err
	}
	p.logger.InfoContext(ctx, "pipeline stream opened",
		slog.String("mode", string(requestContext.Mode)),
		slog.String("endpoint_kind", string(requestContext.EndpointKind)),
		slog.String("data_category", string(requestContext.DataCategory)),
		slog.Int("request_count", len(requests)),
		slog.Int("plan_steps", len(plan.Steps)),
		slog.Duration("stream_open_duration", time.Since(streamStart)),
		slog.Int64("stream_open_duration_ms", time.Since(streamStart).Milliseconds()),
		slog.Duration("total_duration", time.Since(start)),
		slog.Int64("total_duration_ms", time.Since(start).Milliseconds()),
	)
	return plan, stream, nil
}

func (p *Pipeline) Plan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, error) {
	start := time.Now()
	if err := p.prepare(ctx, requests); err != nil {
		return Plan{}, err
	}

	buildStart := time.Now()
	plan, err := p.planner.BuildPlan(ctx, requestContext, requests)
	if err != nil {
		return Plan{}, err
	}
	p.logger.InfoContext(ctx, "pipeline plan completed",
		slog.String("mode", string(requestContext.Mode)),
		slog.String("endpoint_kind", string(requestContext.EndpointKind)),
		slog.String("data_category", string(requestContext.DataCategory)),
		slog.Int("request_count", len(requests)),
		slog.Int("plan_steps", len(plan.Steps)),
		slog.Duration("build_plan_duration", time.Since(buildStart)),
		slog.Int64("build_plan_duration_ms", time.Since(buildStart).Milliseconds()),
		slog.Duration("total_duration", time.Since(start)),
		slog.Int64("total_duration_ms", time.Since(start).Milliseconds()),
	)

	return plan, nil
}

func (p *Pipeline) prepare(ctx context.Context, requests []Request) error {
	start := time.Now()
	if err := p.validator.Validate(ctx, requests); err != nil {
		return err
	}
	validateDuration := time.Since(start)

	parseStart := time.Now()
	parsedExpressions := 0
	for i := range requests {
		request := requests[i]
		if request.Filters == nil {
			continue
		}
		parsedExpressions += len(request.Filters.Expressions)
		parsed, err := p.parser.Parse(ctx, request.Filters.Expressions, request.Filters.FilterTimeZone)
		if err != nil {
			return err
		}
		requests[i].Filters.Parsed = parsed
	}
	parseDuration := time.Since(parseStart)
	p.logger.InfoContext(ctx, "pipeline request prepared",
		slog.Int("request_count", len(requests)),
		slog.Int("filter_expression_count", parsedExpressions),
		slog.Duration("validation_duration", validateDuration),
		slog.Int64("validation_duration_ms", validateDuration.Milliseconds()),
		slog.Duration("filter_parse_duration", parseDuration),
		slog.Int64("filter_parse_duration_ms", parseDuration.Milliseconds()),
		slog.Duration("total_duration", time.Since(start)),
		slog.Int64("total_duration_ms", time.Since(start).Milliseconds()),
	)

	return nil
}
