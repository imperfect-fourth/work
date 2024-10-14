package transformer

type Option func(Transformer)

func WithErrorChan(err chan error) Option {
	return func(t Transformer) {
		t.setErrorChan(err)
	}
}

func WithWorkerPoolSize(n int) Option {
	return func(t Transformer) {
		t.setWorkerPoolSize(n)
	}
}

func WithSpanName(name string) Option {
	return func(t Transformer) {
		t.setSpanName(name)
	}
}
