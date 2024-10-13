package job

import "context"

type Job[T any] struct {
	ctx context.Context

	input T
}

func New[T any](ctx context.Context, in T) Job[T] {
	return Job[T]{
		ctx:   ctx,
		input: in,
	}
}

func (j Job[T]) Context() context.Context {
	return j.ctx
}

func (j Job[T]) Input() T {
	return j.input
}

type Queue[T any] chan Job[T]
