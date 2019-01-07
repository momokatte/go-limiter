package limiter

import (
	"sync"
	"time"
)

type FixedIntervalLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     time.Time
}

func NewFixedIntervalLimiter(interval time.Duration) *FixedIntervalLimiter {
	return &FixedIntervalLimiter{
		interval: interval,
	}
}

func (l *FixedIntervalLimiter) CheckWait() {
	l.mu.Lock()
	next := l.last.Add(l.interval)
	t := time.Now()
	if !t.Before(next) {
		l.last = t
	} else {
		time.Sleep(next.Sub(t))
		l.last = time.Now()
	}
	l.mu.Unlock()
}
