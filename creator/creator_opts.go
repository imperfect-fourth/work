package creator

import "time"

type CreatorOpt func(Creator)

func WithCooldown(interval time.Duration) CreatorOpt {
	return func(c Creator) {
		c.withCooldown(interval)
	}
}

func WithQueueSize(s int) CreatorOpt {
	return func(c Creator) {
		c.withQueueSize(s)
	}
}
