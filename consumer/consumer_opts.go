package consumer

type Option func(Consumer)

func WithErrorChan(err chan error) Option {
	return func(c Consumer) {
		c.setErrorChan(err)
	}
}

func WithParallelism(p int) Option {
	return func(c Consumer) {
		c.setParallelism(p)
	}
}
