package consumer

type Option func(Consumer)

func WithErrorChan(err chan error) Option {
	return func(c Consumer) {
		c.setErrorChan(err)
	}
}

func WithWorkerPoolSize(n int) Option {
	return func(c Consumer) {
		c.setWorkerPoolSize(n)
	}
}
