package consumer

type Consumer interface {
	ConsumeOnce()
	Consume()
	Work()

	setErrorChan(err chan error)
	setParallelism(int)
}

func New[In any](in chan In, fn func(In) error, opts ...Option) (Consumer, chan error) {
	c := &consumer[In]{
		fn:          fn,
		in:          in,
		err:         make(chan error),
		parallelism: 1,
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

	parallelism int
}

func (c consumer[In]) ConsumeOnce() {
	if err := c.fn(<-c.in); err != nil {
		c.err <- err
	}
}

func (c consumer[In]) Consume() {
	throttle := make(chan bool, c.parallelism)
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

func (c *consumer[In]) setParallelism(p int) {
	c.parallelism = p
}
