package consumer

type ConsumerOpt func(Consumer)

func WithParallelism(p int) ConsumerOpt {
	return func(c Consumer) {
		c.withParallelism(p)
	}
}
