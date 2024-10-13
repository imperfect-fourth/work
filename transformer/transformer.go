package transformer

import (
	"fmt"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Transformer interface {
	TransformOnce()
	Transform()
	Work()

	setErrorChan(chan error)
	setWorkerPoolSize(int)
}

func New[In any, Out any](name string, in job.Queue[In], fn func(In) (Out, error), opts ...Option) (Transformer, job.Queue[Out], chan error) {
	t := &transformer[In, Out]{
		name:           name,
		fn:             fn,
		in:             in,
		out:            make(job.Queue[Out]),
		err:            make(chan error),
		workerPoolSize: 1,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t, t.out, t.err
}

type transformer[In any, Out any] struct {
	name string
	fn   func(In) (Out, error)
	in   job.Queue[In]
	out  job.Queue[Out]
	err  chan error

	workerPoolSize int
}

func (t transformer[In, Out]) TransformOnce() {
	i := <-t.in
	ctx, span := otel.Tracer(t.name).Start(i.Context(), "transform jobs")
	defer func() {
		fmt.Printf("transformer span %+v\n", span)
		span.End()
	}()
	o, err := t.fn(i.Input())
	if err != nil {
		t.err <- err
		return
	}
	t.out <- job.New(ctx, o)
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

func (t *transformer[In, Out]) setErrorChan(err chan error) {
	t.err = err
}

func (t *transformer[In, Out]) setWorkerPoolSize(n int) {
	t.workerPoolSize = n
}
