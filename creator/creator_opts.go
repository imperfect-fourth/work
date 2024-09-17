package creator

import "time"

type CreatorOpt func(Creator)

func WithCooldown(interval time.Duration) CreatorOpt {
	return func(c Creator) {
		c.setCooldown(interval)
	}
}

func WithErrorChan(err chan error) CreatorOpt {
	return func(c Creator) {
		c.setErrorChan(err)
	}
}

func WithQueueSize(s int) CreatorOpt {
	return func(c Creator) {
		c.setQueueSize(s)
	}
}
