package consumer

import (
	"context"
	"fmt"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Consumer[In any] interface {
	ConsumeOnce()
	Consume()
	Work()

	Input() job.Queue[In]
	Error() job.Queue[error]

	WithInput(job.Queue[In]) Consumer[In]
	WithErrorQueue(job.Queue[error]) Consumer[In]
	WithWorkerPoolSize(int) Consumer[In]
}

func New[In any](name string, fn func(context.Context, In) error) Consumer[In] {
	c := &consumer[In]{
		name:           name,
		fn:             fn,
		in:             make(job.Queue[In]),
		err:            make(job.Queue[error]),
		workerPoolSize: 1,
	}
	return c
}

type consumer[In any] struct {
	name     string

	fn  func(context.Context, In) error
	in  job.Queue[In]
	err job.Queue[error]

	workerPoolSize int
}

func (c consumer[In]) ConsumeOnce() {
	j := <-c.in
	defer j.End()

	ctx, span := otel.Tracer(c.name).Start(j.Context(), c.name)
	defer span.End()

	fnctx, fnspan := otel.Tracer(c.name).Start(ctx, fmt.Sprintf("%s - run function", c.name))
	err := c.fn(fnctx, j.Input())
	fnspan.End()
	if err != nil {
		_, errJob := job.New(ctx, err)
		c.err <- errJob
	}
}

func (c consumer[In]) Consume() {
	throttle := make(chan bool, c.workerPoolSize)
	for {
		throttle <- true
		go func() {
			c.ConsumeOnce()
			<-throttle
		}()
	}
}

func (c consumer[In]) Work() {
	c.Consume()
}

func (c *consumer[In]) Input() job.Queue[In] {
	return c.in
}
func (c *consumer[In]) Error() job.Queue[error] {
	return c.err
}

func (c *consumer[In]) WithInput(in job.Queue[In]) Consumer[In] {
	c.in = in
	return c
}
func (c *consumer[In]) WithErrorQueue(err job.Queue[error]) Consumer[In] {
	c.err = err
	return c
}
func (c *consumer[In]) WithWorkerPoolSize(n int) Consumer[In] {
	c.workerPoolSize = n
	return c
}
