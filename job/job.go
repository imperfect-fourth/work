package job

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type contextKey struct{}

var jobGlobalsKey = contextKey{}

type jobGlobals struct {
	id       uuid.UUID
	rootSpan trace.Span
}

type Job[T any] struct {
	*jobGlobals

	ctx   context.Context
	input T
}

func New[T any](ctx context.Context, in T) (context.Context, Job[T]) {
	if ctx == nil {
		ctx = context.Background()
	}
	if jg := ctx.Value(jobGlobalsKey); jg != nil {
		return ctx, Job[T]{
			jobGlobals: jg.(*jobGlobals),
			ctx:        ctx,
			input:      in,
		}
	}
	id := uuid.New()
	ctx, rootSpan := otel.Tracer("root").Start(ctx, fmt.Sprintf("job %s", id.String()))
	jg := &jobGlobals{
		id:       id,
		rootSpan: rootSpan,
	}
	rootCtx := context.WithValue(ctx, jobGlobalsKey, jg)
	return rootCtx, Job[T]{
		ctx:        rootCtx,
		jobGlobals: jg,
		input:      in,
	}
}

func (j Job[T]) Context() context.Context {
	return j.ctx
}

func (j Job[T]) ID() uuid.UUID {
	return j.id
}

func (j Job[T]) Input() T {
	return j.input
}

func (j Job[T]) End() {
	j.rootSpan.End()
}

type Queue[T any] chan Job[T]
