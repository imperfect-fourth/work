package consumer

type ConsumerOpt func(Consumer)

func WithErrorChan(err chan error) ConsumerOpt {
	return func(c Consumer) {
		c.setErrorChan(err)
	}
}

func WithParallelism(p int) ConsumerOpt {
	return func(c Consumer) {
		c.setParallelism(p)
	}
}
