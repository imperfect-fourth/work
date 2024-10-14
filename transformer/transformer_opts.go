package transformer

import "github.com/imperfect-fourth/work/job"

type Option func(Transformer)

func WithErrorChan(err job.Queue[error]) Option {
	return func(t Transformer) {
		t.setErrorQueue(err)
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
