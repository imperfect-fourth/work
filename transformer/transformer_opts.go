package transformer

type Option func(Transformer)

func WithErrorChan(err chan error) Option {
	return func(t Transformer) {
		t.setErrorChan(err)
	}
}

func WithParallelism(p int) Option {
	return func(t Transformer) {
		t.setParallelism(p)
	}
}

func WithQueueSize(s int) Option {
	return func(t Transformer) {
		t.setQueueSize(s)
	}
}
