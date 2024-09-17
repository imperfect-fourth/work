package transformer

import "github.com/imperfect-fourth/work"

type Transformer interface {
	work.Worker

	TransformOnce() error
	Transform() error

	withQueueSize(int)
	withParallelism(int)
}

func NewTransformer[T any, D any, U chan T, V chan D](in U, fn func(T) (D, error), opts ...TransformerOpt) (Transformer, V) {
	out := make(V)
	t := &transformer[T, D, U, V]{
		fn:          fn,
		in:          in,
		out:         out,
		parallelism: 1,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t, t.out
}

type transformer[T any, D any, U chan T, V chan D] struct {
	fn  func(T) (D, error)
	in  U
	out V

	parallelism int
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
	throttle := make(chan bool, t.parallelism)
	for {
		throttle <- true
		var err error
		go func() {
			err = t.TransformOnce()
		}()
		if err != nil {
			return err
		}
		<-throttle
	}
}

func (t transformer[T, D, U, V]) Work() error {
	return t.Transform()
}

func (t *transformer[T, D, U, V]) withParallelism(p int) {
	t.parallelism = p
}

func (t *transformer[T, D, U, V]) withQueueSize(s int) {
	t.out = make(V, s)
}
