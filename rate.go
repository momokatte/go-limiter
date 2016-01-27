package limiter

import (
	"time"
)

type Rate struct {
	Count    int
	Duration time.Duration
}

func NewRate(count int, duration time.Duration) (r Rate) {
	r = Rate{
		Count:    count,
		Duration: duration,
	}
	return
}
