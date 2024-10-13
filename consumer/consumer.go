package consumer

import (
	"fmt"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
)

type Consumer interface {
	ConsumeOnce()
	Consume()
	Work()

	setErrorChan(err chan error)
	setWorkerPoolSize(int)
}

func New[In any](name string, in job.Queue[In], fn func(In) error, opts ...Option) (Consumer, chan error) {
	c := &consumer[In]{
		name:           name,
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
	name string

	fn  func(In) error
	in  job.Queue[In]
	err chan error

	workerPoolSize int
}

func (c consumer[In]) ConsumeOnce() {
	i := <-c.in

	_, span := otel.Tracer(c.name).Start(i.Context(), "consume jobs")
	defer func() {
		fmt.Printf("consumer span %+v\n", span)
		span.End()
	}()
	if err := c.fn(i.Input()); err != nil {
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
