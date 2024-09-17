package consumer

import "github.com/imperfect-fourth/work"

type Consumer interface {
	work.Worker

	ConsumeOnce() error
	Consume() error

	withParallelism(int)
}

func NewConsumer[T any, U chan T](in U, fn func(T) error) Consumer {
	return &consumer[T, U]{
		fn: fn,
		in: in,
	}
}

type consumer[T any, U chan T] struct {
	fn func(T) error
	in U

	parallelism int
}

func (c consumer[T, U]) ConsumeOnce() error {
	i := <-c.in
	return c.fn(i)
}

func (c consumer[T, U]) Consume() error {
	for {
		if err := c.ConsumeOnce(); err != nil {
			return err
		}
	}
}

func (c consumer[T, U]) Work() error {
	return c.Consume()
}

func (c *consumer[T, U]) withParallelism(p int) {
	c.parallelism = p
}
