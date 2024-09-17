package transformer

import "github.com/imperfect-fourth/work"

type Transformer interface {
	work.Worker

	TransformOnce()
	Transform()

	setQueueSize(int)
	setParallelism(int)
}

func NewTransformer[In any, Out any](in chan In, fn func(In) (Out, error), opts ...TransformerOpt) (Transformer, chan Out, chan error) {
	t := &transformer[In, Out]{
		fn:          fn,
		in:          in,
		out:         make(chan Out),
		err:         make(chan error),
		parallelism: 1,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t, t.out, t.err
}

type transformer[In any, Out any] struct {
	fn  func(In) (Out, error)
	in  chan In
	out chan Out
	err chan error

	parallelism int
}

func (t transformer[In, Out]) TransformOnce() {
	i := <-t.in
	o, err := t.fn(i)
	if err != nil {
		t.err <- err
		return
	}
	t.out <- o
}

func (t transformer[In, Out]) Transform() {
	throttle := make(chan bool, t.parallelism)
	for {
		throttle <- true
		go func() {
			t.TransformOnce()
			<-throttle
		}()
	}
}

func (t transformer[In, Out]) Work() {
	t.Transform()
}

func (t *transformer[In, Out]) setParallelism(p int) {
	t.parallelism = p
}

func (t *transformer[In, Out]) setQueueSize(s int) {
	t.out = make(chan Out, s)
}
