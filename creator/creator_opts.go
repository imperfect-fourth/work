package creator

import "time"

type Option func(Creator)

func WithCooldown(interval time.Duration) Option {
	return func(c Creator) {
		c.setCooldown(interval)
	}
}

func WithErrorChan(err chan error) Option {
	return func(c Creator) {
		c.setErrorChan(err)
	}
}
