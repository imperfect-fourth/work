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
		fn:          fn,
		in:          in,
		parallelism: 1,
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
	throttle := make(chan bool, c.parallelism)
	for {
		throttle <- true
		var err error
		go func() {
			err = c.ConsumeOnce()
		}()
		if err != nil {
			return err
		}
		<-throttle
	}
}

func (c consumer[T, U]) Work() error {
	return c.Consume()
}

func (c *consumer[T, U]) withParallelism(p int) {
	c.parallelism = p
}
