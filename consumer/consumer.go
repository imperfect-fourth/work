package consumer

import "github.com/imperfect-fourth/work"

type Consumer interface {
	work.Worker

	ConsumeOnce()
	Consume()

	withParallelism(int)
}

func NewConsumer[In any](in chan In, fn func(In) error, opts ...ConsumerOpt) (Consumer, chan error) {
	c := &consumer[In]{
		fn:          fn,
		in:          in,
		err:         make(chan error),
		parallelism: 1,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, c.err
}

type consumer[In any] struct {
	fn  func(In) error
	in  chan In
	err chan error

	parallelism int
}

func (c consumer[In]) ConsumeOnce() {
	if err := c.fn(<-c.in); err != nil {
		c.err <- err
	}
}

func (c consumer[In]) Consume() {
	throttle := make(chan bool, c.parallelism)
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

func (c *consumer[In]) withParallelism(p int) {
	c.parallelism = p
}
