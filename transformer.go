package work

type Transformer[T any, D any, U chan T, V chan D] interface {
	Worker

	Transform() error
	InChannel() U
	OutChannel() V
}

type transformer[T any, D any, U chan T, V chan D] struct {
	fn  func(T) (D, error)
	in  U
	out V
}

func (t transformer[T, D, U, V]) Transform() error {
	i := <-t.in
	o, err := t.fn(i)
	if err != nil {
		return err
	}
	t.out <- o
	return nil
}

func (c transformer[T, D, U, V]) Work() error {
	return c.Transform()
}

func (t transformer[T, D, U, V]) InChannel() U {
	return t.in
}

func (t transformer[T, D, U, V]) OutChannel() V {
	return t.out
}

func NewTransformer[T any, D any, U chan T, V chan D](in U, fn func(T) (D, error)) (Transformer[T, D, U, V], V) {
	out := make(V)
	return transformer[T, D, U, V]{fn, in, out}, out
}
