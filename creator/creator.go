package creator

import (
	"time"

	"github.com/imperfect-fourth/work"
)

type Creator interface {
	work.Worker
	CreateOnce()
	Create()

	setCooldown(time.Duration)
	setQueueSize(int)
}

func NewCreator[Out any](fn func() ([]Out, error), opts ...CreatorOpt) (Creator, chan Out, chan error) {
	c := creator[Out]{
		fn:  fn,
		out: make(chan Out),
		err: make(chan error),
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c, c.out, c.err
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

func (c *creator[Out]) setQueueSize(s int) {
	c.out = make(chan Out, s)
}
