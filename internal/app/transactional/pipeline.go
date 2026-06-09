package transactional

import "context"

type Validator interface {
	Validate(context.Context, []Request) error
}

type FilterParser interface {
	Parse(context.Context, []string) (FilterSet, error)
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
}

func NewPipeline(validator Validator, parser FilterParser, planner Planner, executor Executor) *Pipeline {
	return &Pipeline{
		validator: validator,
		parser:    parser,
		planner:   planner,
		executor:  executor,
	}
}

func (p *Pipeline) Execute(ctx context.Context, requestContext RequestContext, requests []Request) (Response, error) {
	_, response, err := p.ExecuteWithPlan(ctx, requestContext, requests)
	return response, err
}

func (p *Pipeline) ExecuteWithPlan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, Response, error) {
	plan, err := p.Plan(ctx, requestContext, requests)
	if err != nil {
		return Plan{}, Response{}, err
	}

	response, err := p.executor.Execute(ctx, plan)
	if err != nil {
		return Plan{}, Response{}, err
	}
	return plan, response, nil
}

func (p *Pipeline) Stream(ctx context.Context, requestContext RequestContext, requests []Request) (Stream, error) {
	_, stream, err := p.StreamWithPlan(ctx, requestContext, requests)
	return stream, err
}

func (p *Pipeline) StreamWithPlan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, Stream, error) {
	plan, err := p.Plan(ctx, requestContext, requests)
	if err != nil {
		return Plan{}, nil, err
	}

	stream, err := p.executor.Stream(ctx, plan)
	if err != nil {
		return Plan{}, nil, err
	}
	return plan, stream, nil
}

func (p *Pipeline) Plan(ctx context.Context, requestContext RequestContext, requests []Request) (Plan, error) {
	if err := p.prepare(ctx, requests); err != nil {
		return Plan{}, err
	}

	plan, err := p.planner.BuildPlan(ctx, requestContext, requests)
	if err != nil {
		return Plan{}, err
	}

	return plan, nil
}

func (p *Pipeline) prepare(ctx context.Context, requests []Request) error {
	if err := p.validator.Validate(ctx, requests); err != nil {
		return err
	}

	for i := range requests {
		request := requests[i]
		if request.Filters == nil {
			continue
		}
		parsed, err := p.parser.Parse(ctx, request.Filters.Expressions)
		if err != nil {
			return err
		}
		requests[i].Filters.Parsed = parsed
	}

	return nil
}
