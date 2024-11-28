package transformer

import (
	"context"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Transformer interface {
	TransformOnce()
	Transform()
	Work()

	setErrorQueue(job.Queue[error])
	setWorkerPoolSize(int)
	setSpanName(string)
}

func New[In any, Out any](name string, in job.Queue[In], fn func(context.Context, In) (Out, error), opts ...Option) (Transformer, job.Queue[Out], job.Queue[error]) {
	t := &transformer[In, Out]{
		name:           name,
		spanName:       "transform job",
		fn:             fn,
		in:             in,
		out:            make(job.Queue[Out]),
		err:            make(job.Queue[error]),
		workerPoolSize: 1,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t, t.out, t.err
}

type transformer[In any, Out any] struct {
	name     string
	spanName string

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

	fnctx, fnspan := otel.Tracer(t.name).Start(ctx, t.spanName)
	o, err := t.fn(fnctx, j.Input())
	fnspan.End()
	if err != nil {
		_, errJob := job.New(ctx, err)
		t.err <- errJob
		return
	}

	_, jobspan := otel.Tracer(t.name).Start(ctx, "output queue wait")
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

func (t *transformer[In, Out]) setErrorQueue(err job.Queue[error]) {
	t.err = err
}

func (t *transformer[In, Out]) setWorkerPoolSize(n int) {
	t.workerPoolSize = n
}

func (t *transformer[In, Out]) setSpanName(name string) {
	t.spanName = name
}
