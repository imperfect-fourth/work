package consumer

import "github.com/imperfect-fourth/work/job"

type Option func(Consumer)

func WithErrorQueue(err job.Queue[error]) Option {
	return func(c Consumer) {
		c.setErrorQueue(err)
	}
}

func WithWorkerPoolSize(n int) Option {
	return func(c Consumer) {
		c.setWorkerPoolSize(n)
	}
}

func WithSpanName(name string) Option {
	return func(c Consumer) {
		c.setSpanName(name)
	}
}
