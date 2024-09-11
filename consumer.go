package work

type Consumer[T any, U chan T] interface {
	Worker
	ConsumeOnce() error
	Consume() error
	InChannel() U
}

type consumer[T any, U chan T] struct {
	fn func(T) error
	in U
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

func (c consumer[T, U]) InChannel() U {
	return c.in
}

func NewConsumer[T any, U chan T](in U, fn func(T) error) Consumer[T, U] {
	return consumer[T, U]{fn, in}
}
