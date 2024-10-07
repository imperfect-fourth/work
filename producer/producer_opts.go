package producer

import "time"

type Option func(Producer)

func WithCooldown(interval time.Duration) Option {
	return func(c Producer) {
		c.setCooldown(interval)
	}
}

func WithErrorChan(err chan error) Option {
	return func(c Producer) {
		c.setErrorChan(err)
	}
}
