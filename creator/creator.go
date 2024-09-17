package creator

import (
	"time"
)

type Creator interface {
	CreateOnce()
	Create()
	Work()

	setCooldown(time.Duration)
	setErrorChan(chan error)
}

func New[Out any](fn func() ([]Out, error), opts ...Option) (Creator, chan Out, chan error) {
	c := &creator[Out]{
		fn:  fn,
		out: make(chan Out),
		err: make(chan error),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, c.out, c.err
}

type creator[Out any] struct {
	fn  func() ([]Out, error)
	out chan Out
	err chan error

	cooldown time.Duration
}

func (c creator[Out]) CreateOnce() {
	out, err := c.fn()
	if err != nil {
		c.err <- err
	}
	for _, o := range out {
		c.out <- o
	}
}

func (c creator[Out]) Create() {
	for {
		c.CreateOnce()
		time.Sleep(c.cooldown)
	}
}

func (c creator[Out]) Work() {
	c.Create()
}

func (c *creator[Out]) setCooldown(t time.Duration) {
	c.cooldown = t
}

func (c *creator[Out]) setErrorChan(err chan error) {
	c.err = err
}
