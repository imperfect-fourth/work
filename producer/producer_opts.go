package producer

import (
	"time"

	"github.com/imperfect-fourth/work/job"
)

type Option func(Producer)

func WithCooldown(interval time.Duration) Option {
	return func(p Producer) {
		p.setCooldown(interval)
	}
}

func WithErrorQueue(err job.Queue[error]) Option {
	return func(p Producer) {
		p.setErrorQueue(err)
	}
}
