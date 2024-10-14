package consumer

import (
	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Consumer interface {
	ConsumeOnce()
	Consume()
	Work()

	setErrorChan(err chan error)
	setWorkerPoolSize(int)
	setSpanName(string)
}

func New[In any](name string, in job.Queue[In], fn func(In) error, opts ...Option) (Consumer, chan error) {
	c := &consumer[In]{
		name:           name,
		spanName:       "consume job",
		fn:             fn,
		in:             in,
		err:            make(chan error),
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

	fn  func(In) error
	in  job.Queue[In]
	err chan error

	workerPoolSize int
}

func (c consumer[In]) ConsumeOnce() {
	j := <-c.in
	defer j.End()

	ctx, span := otel.Tracer(c.name).Start(j.Context(), c.name)
	defer span.End()

	_, jobspan := otel.Tracer(c.name).Start(ctx, c.spanName)
	err := c.fn(j.Input())
	jobspan.End()
	if err != nil {
		c.err <- err
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

func (c *consumer[In]) setErrorChan(err chan error) {
	c.err = err
}

func (c *consumer[In]) setWorkerPoolSize(n int) {
	c.workerPoolSize = n
}

func (c *consumer[In]) setSpanName(name string) {
	c.spanName = name
}
