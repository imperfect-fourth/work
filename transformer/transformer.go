package transformer

import (
	"context"
	"fmt"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Transformer[In, Out any] interface {
	TransformOnce()
	Transform()
	Work()

	Input() job.Queue[In]
	Output() job.Queue[Out]
	Error() job.Queue[error]

	WithInput(job.Queue[In]) Transformer[In, Out]
	WithErrorQueue(job.Queue[error]) Transformer[In, Out]
	WithWorkerPoolSize(int) Transformer[In, Out]
}

func New[In any, Out any](name string, fn func(context.Context, In) (Out, error)) Transformer[In, Out] {
	t := &transformer[In, Out]{
		name:           name,
		fn:             fn,
		in:             make(job.Queue[In]),
		out:            make(job.Queue[Out]),
		err:            make(job.Queue[error]),
		workerPoolSize: 1,
	}
	return t
}

type transformer[In any, Out any] struct {
	name string

	fn  func(context.Context, In) (Out, error)
	in  job.Queue[In]
	out job.Queue[Out]
	err job.Queue[error]

	workerPoolSize int
}

func (t transformer[In, Out]) TransformOnce() {
	j := <-t.in
	ctx, span := otel.Tracer(t.name).Start(j.Context(), t.name)
	defer span.End()

	fnctx, fnspan := otel.Tracer(t.name).Start(ctx, fmt.Sprintf("%s - run function", t.name))
	o, err := t.fn(fnctx, j.Input())
	fnspan.End()
	if err != nil {
		_, errJob := job.New(ctx, err)
		t.err <- errJob
		return
	}

	_, jobspan := otel.Tracer(t.name).Start(ctx, fmt.Sprintf("%s - output queue wait", t.name))
	_, newJob := job.New(j.Context(), o)
	t.out <- newJob
	jobspan.End()
}

func (t transformer[In, Out]) Transform() {
	throttle := make(chan bool, t.workerPoolSize)
	for {
		throttle <- true
		go func() {
			t.TransformOnce()
			<-throttle
		}()
	}
}

func (t transformer[In, Out]) Work() {
	t.Transform()
}

func (t *transformer[In, Out]) Input() job.Queue[In] {
	return t.in
}
func (t *transformer[In, Out]) Output() job.Queue[Out] {
	return t.out
}
func (t *transformer[In, Out]) Error() job.Queue[error] {
	return t.err
}

func (t *transformer[In, Out]) WithInput(in job.Queue[In]) Transformer[In, Out] {
	t.in = in
	return t
}
func (t *transformer[In, Out]) WithErrorQueue(err job.Queue[error]) Transformer[In, Out] {
	t.err = err
	return t
}
func (t *transformer[In, Out]) WithWorkerPoolSize(n int) Transformer[In, Out] {
	t.workerPoolSize = n
	return t
}
