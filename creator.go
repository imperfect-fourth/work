package work

type Creator[T any, U chan T] interface {
	Worker
	CreateOnce() error
	Create() error
	OutChannel() U
}

type creator[T any, U chan T] struct {
	fn  func() ([]T, error)
	out U
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
	}
}

func (c creator[T, U]) Work() error {
	return c.Create()
}

func (c creator[T, U]) OutChannel() U {
	return c.out
}

func NewCreator[T any, U chan T](fn func() ([]T, error)) (Creator[T, U], U) {
	out := make(U)
	return creator[T, U]{fn, out}, out
}
