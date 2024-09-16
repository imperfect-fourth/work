package work

import "time"

type Creator interface {
	Worker
	CreateOnce() error
	Create() error
	setCooldown(time.Duration)
}

type creator[T any, U chan T] struct {
	fn       func() ([]T, error)
	out      U
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

func (c creator[T, U]) OutChannel() U {
	return c.out
}

func (c creator[T, U]) setCooldown(t time.Duration) {
	c.cooldown = t
}

func NewCreator[T any, U chan T](fn func() ([]T, error), opts ...CreatorOpt) (Creator, U) {
	out := make(U)
	c := creator[T, U]{
		fn:  fn,
		out: out,
	}, out

	for _, opt := range opts {
		opt(c)
	}
	return c
}

type CreatorOpt func(Creator)

func WithCooldown(interval time.Duration) CreatorOpt {
	return func(c Creator) {
		c.setCooldown(interval)
	}
}
