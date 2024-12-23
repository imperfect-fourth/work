package consumer

import (
	"context"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Consumer interface {
	ConsumeOnce()
	Consume()
	Work()

	setErrorQueue(err job.Queue[error])
	setWorkerPoolSize(int)
	setSpanName(string)
}

func New[In any](name string, in job.Queue[In], fn func(context.Context, In) error, opts ...Option) (Consumer, job.Queue[error]) {
	c := &consumer[In]{
		name:           name,
		spanName:       "consume job",
		fn:             fn,
		in:             in,
		err:            make(job.Queue[error]),
		workerPoolSize: 1,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, c.err
}

type consumer[In any] struct {
	name     string
	spanName string

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

	fnctx, fnspan := otel.Tracer(c.name).Start(ctx, c.spanName)
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

func (c *consumer[In]) setErrorQueue(err job.Queue[error]) {
	c.err = err
}

func (c *consumer[In]) setWorkerPoolSize(n int) {
	c.workerPoolSize = n
}

func (c *consumer[In]) setSpanName(name string) {
	c.spanName = name
}
