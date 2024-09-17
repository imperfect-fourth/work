package creator

import (
	"time"

	"github.com/imperfect-fourth/work"
)

type Creator interface {
	work.Worker
	CreateOnce() error
	Create() error

	withCooldown(time.Duration)
	withQueueSize(int)
}

func NewCreator[T any, U chan T](fn func() ([]T, error), opts ...CreatorOpt) (Creator, U) {
	out := make(U)
	c := creator[T, U]{
		fn:  fn,
		out: out,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c, c.out
}

type creator[T any, U chan T] struct {
	fn  func() ([]T, error)
	out U

	cooldown time.Duration
}

func (c creator[T, U]) CreateOnce() error {
	out, err := c.fn()
	if err != nil {
		return err
	}
	for _, o := range out {
		c.out <- o
	}
	return nil
}

func (c creator[T, U]) Create() error {
	for {
		if err := c.CreateOnce(); err != nil {
			return err
		}
		time.Sleep(c.cooldown)
	}
}

func (c creator[T, U]) Work() error {
	return c.Create()
}

func (c *creator[T, U]) withCooldown(t time.Duration) {
	c.cooldown = t
}

func (c *creator[T, U]) withQueueSize(s int) {
	c.out = make(U, s)
}
