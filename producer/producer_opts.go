package producer

import "time"

type Option func(Producer)

func WithCooldown(interval time.Duration) Option {
	return func(p Producer) {
		p.setCooldown(interval)
	}
}

func WithErrorChan(err chan error) Option {
	return func(p Producer) {
		p.setErrorChan(err)
	}
}
