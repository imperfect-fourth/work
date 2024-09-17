package consumer

type Consumer interface {
	ConsumeOnce()
	Consume()
	Work()

	setErrorChan(err chan error)
	setWorkerPoolSize(int)
}

func New[In any](in chan In, fn func(In) error, opts ...Option) (Consumer, chan error) {
	c := &consumer[In]{
		fn:             fn,
		in:             in,
		err:            make(chan error),
		workerPoolSize: 1,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, c.err
}

type consumer[In any] struct {
	fn  func(In) error
	in  chan In
	err chan error

	workerPoolSize int
}

func (c consumer[In]) ConsumeOnce() {
	if err := c.fn(<-c.in); err != nil {
		c.err <- err
	}
}

func (c consumer[In]) Consume() {
	throttle := make(chan bool, c.workerPoolSize)
	for {
		throttle <- true
		go func() {
			c.ConsumeOnce()
			<-throttle
		}()
	}
}

func (c consumer[In]) Work() {
	c.Consume()
}

func (c *consumer[In]) setErrorChan(err chan error) {
	c.err = err
}

func (c *consumer[In]) setWorkerPoolSize(n int) {
	c.workerPoolSize = n
}
