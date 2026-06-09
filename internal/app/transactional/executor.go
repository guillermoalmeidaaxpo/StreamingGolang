package transactional

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"

	"streaming-golang/internal/domain"
)

type Stream interface {
	Next(context.Context) bool
	Item() DataItem
	Err() error
	Close() error
}

type executor struct {
	repositories map[domain.SourceKind]Repository
	maxParallel  int
	transformer  *TransformationProcessor
}

func NewExecutor(repositories map[domain.SourceKind]Repository, maxParallel int) Executor {
	if maxParallel <= 0 {
		maxParallel = 1
	}
	return &executor{
		repositories: repositories,
		maxParallel:  maxParallel,
		transformer:  NewTransformationProcessor(),
	}
}

func (e *executor) Execute(ctx context.Context, plan Plan) (Response, error) {
	if len(plan.Steps) == 0 {
		return Response{ReferenceData: ReferenceData(requestedIDs(plan))}, nil
	}

	allItems := make([]DataItem, 0)
	for _, step := range plan.Steps {
		items, err := e.executeStep(ctx, step)
		if err != nil {
			return Response{}, err
		}
		allItems = append(allItems, items...)
	}

	return Response{
		TransactionalData: allItems,
		ReferenceData:     ReferenceData(requestedIDs(plan)),
	}, nil
}

func (e *executor) executeStep(ctx context.Context, step PlanStep) ([]DataItem, error) {
	results := make(chan []DataItem, len(step.Queries))
	group, gCtx := errgroup.WithContext(ctx)
	semaphore := make(chan struct{}, e.maxParallel)

	for _, query := range step.Queries {
		q := query
		group.Go(func() error {
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-gCtx.Done():
				return gCtx.Err()
			}

			repo, ok := e.repositories[q.Source]
			if !ok {
				return nil
			}

			items, err := repo.Execute(gCtx, q)
			if err != nil {
				return err
			}
			results <- items
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return nil, err
	}
	close(results)

	stepItems := make([]DataItem, 0)
	for items := range results {
		stepItems = append(stepItems, items...)
	}

	return e.transformer.Process(ctx, stepItems, step.Command), nil
}

func requestedIDs(plan Plan) []domain.Identifier {
	ids := make([]domain.Identifier, 0)
	for _, step := range plan.Steps {
		for _, id := range step.Command.IDs {
			ids = append(ids, id)
		}
	}
	return ids
}

func (e *executor) Stream(ctx context.Context, plan Plan) (Stream, error) {
	if len(plan.Steps) == 0 {
		return &sliceStream{}, nil
	}

	// For simple one-step plans, we can stream directly
	if len(plan.Steps) == 1 && len(plan.Steps[0].Queries) == 1 {
		step := plan.Steps[0]
		query := step.Queries[0]
		repo, ok := e.repositories[query.Source]
		if !ok {
			return &sliceStream{}, nil
		}
		
		// If no transformation is needed, use direct stream
		if step.Command.TargetTimeZone == "" && !step.Command.HasAggregations {
			return repo.Stream(ctx, query)
		}
	}

	// For complex plans, use multiplexed transformed stream
	return newTransformedStream(ctx, e, plan)
}

type transformedStream struct {
	ctx     context.Context
	cancel  context.CancelFunc
	results chan DataItem
	err     error
	once    sync.Once
	item    DataItem
}

func newTransformedStream(ctx context.Context, e *executor, plan Plan) (*transformedStream, error) {
	mCtx, cancel := context.WithCancel(ctx)
	s := &transformedStream{
		ctx:     mCtx,
		cancel:  cancel,
		results: make(chan DataItem, 100),
	}

	go func() {
		defer close(s.results)
		
		for _, step := range plan.Steps {
			group, gCtx := errgroup.WithContext(mCtx)
			semaphore := make(chan struct{}, e.maxParallel)
			stepResults := make(chan DataItem, 100)

			group.Go(func() error {
				defer close(stepResults)
				innerGroup, innerGCtx := errgroup.WithContext(gCtx)
				
				for _, query := range step.Queries {
					q := query
					innerGroup.Go(func() error {
						select {
						case semaphore <- struct{}{}:
							defer func() { <-semaphore }()
						case <-innerGCtx.Done():
							return innerGCtx.Err()
						}

						repo, ok := e.repositories[q.Source]
						if !ok {
							return nil
						}

						stream, err := repo.Stream(innerGCtx, q)
						if err != nil {
							return err
						}
						defer stream.Close()

						for stream.Next(innerGCtx) {
							select {
							case stepResults <- stream.Item():
							case <-innerGCtx.Done():
								return innerGCtx.Err()
							}
						}
						return stream.Err()
					})
				}
				return innerGroup.Wait()
			})

			// Apply transformations to step results
			// For streaming, we apply per-item if possible, but aggregations might need buffering
			// For now, assume per-item (like TimeZone)
			go func() {
				for item := range stepResults {
					transformed := e.transformer.processItem(item, step.Command, nil) // Location will be re-loaded inside if needed, or we can optimize
					select {
					case s.results <- transformed:
					case <-mCtx.Done():
						return
					}
				}
			}()

			if err := group.Wait(); err != nil {
				s.err = err
				return
			}
		}
	}()

	return s, nil
}

func (s *transformedStream) Next(ctx context.Context) bool {
	select {
	case item, ok := <-s.results:
		if !ok {
			return false
		}
		s.item = item
		return true
	case <-ctx.Done():
		s.err = ctx.Err()
		return false
	case <-s.ctx.Done():
		return false
	}
}

func (s *transformedStream) Item() DataItem {
	return s.item
}

func (s *transformedStream) Err() error {
	return s.err
}

func (s *transformedStream) Close() error {
	s.once.Do(s.cancel)
	return nil
}

type sliceStream struct {
	items []DataItem
	index int
	item  DataItem
}

func (s *sliceStream) Next(ctx context.Context) bool {
	if ctx.Err() != nil || s.index >= len(s.items) {
		return false
	}
	s.item = s.items[s.index]
	s.index++
	return true
}

func (s *sliceStream) Item() DataItem {
	return s.item
}

func (s *sliceStream) Err() error {
	return nil
}

func (s *sliceStream) Close() error {
	return nil
}
