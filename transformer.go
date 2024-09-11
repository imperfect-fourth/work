package work

type Transformer[T any, D any, U chan T, V chan D] interface {
	Worker

	TransformOnce() error
	Transform() error
	InChannel() U
	OutChannel() V
}

type transformer[T any, D any, U chan T, V chan D] struct {
	fn  func(T) (D, error)
	in  U
	out V
}

func (t transformer[T, D, U, V]) TransformOnce() error {
	i := <-t.in
	o, err := t.fn(i)
	if err != nil {
		return err
	}
	t.out <- o
	return nil
}

func (t transformer[T, D, U, V]) Transform() error {
	for {
		if err := t.TransformOnce(); err != nil {
			return err
		}
	}
}

func (t transformer[T, D, U, V]) Work() error {
	return t.Transform()
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
