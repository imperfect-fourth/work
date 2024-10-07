package transformer

type Transformer interface {
	TransformOnce()
	Transform()
	Work()

	setErrorChan(chan error)
	setWorkerPoolSize(int)
}

func New[In any, Out any](in chan In, fn func(In) (Out, error), opts ...Option) (Transformer, chan Out, chan error) {
	t := &transformer[In, Out]{
		fn:             fn,
		in:             in,
		out:            make(chan Out),
		err:            make(chan error),
		workerPoolSize: 1,
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

	workerPoolSize int
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
	throttle := make(chan bool, t.workerPoolSize)
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

func (t *transformer[In, Out]) setErrorChan(err chan error) {
	t.err = err
}

func (t *transformer[In, Out]) setWorkerPoolSize(n int) {
	t.workerPoolSize = n
}
