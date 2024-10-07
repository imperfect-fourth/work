package producer

import (
	"time"
)

type Producer interface {
	ProduceOnce()
	Produce()
	Work()

	setCooldown(time.Duration)
	setErrorChan(chan error)
}

func New[Out any](fn func() ([]Out, error), opts ...Option) (Producer, chan Out, chan error) {
	c := &producer[Out]{
		fn:  fn,
		out: make(chan Out),
		err: make(chan error),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, c.out, c.err
}

type producer[Out any] struct {
	fn  func() ([]Out, error)
	out chan Out
	err chan error

	cooldown time.Duration
}

func (c producer[Out]) ProduceOnce() {
	out, err := c.fn()
	if err != nil {
		c.err <- err
	}
	for _, o := range out {
		c.out <- o
	}
}

func (c producer[Out]) Produce() {
	for {
		c.ProduceOnce()
		time.Sleep(c.cooldown)
	}
}

func (c producer[Out]) Work() {
	c.Produce()
}

func (c *producer[Out]) setCooldown(t time.Duration) {
	c.cooldown = t
}

func (c *producer[Out]) setErrorChan(err chan error) {
	c.err = err
}
