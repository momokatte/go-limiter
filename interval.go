package limiter

import (
	"sync"
	"time"
)

type IntervalLimiter struct {
	mu       sync.Mutex
	interval time.Duration
	recheck  time.Duration
	last     time.Time
}

func NewIntervalLimiter(interval time.Duration) *IntervalLimiter {
	return &IntervalLimiter{
		interval: interval,
		recheck:  time.Second / 4,
	}
}

func (l *IntervalLimiter) SetInterval(interval time.Duration) {
	l.interval = interval
}

func (l *IntervalLimiter) CheckWait() {
	l.mu.Lock()
	defer l.mu.Unlock()
	var t time.Time
	for {
		next := l.last.Add(l.interval)
		t = time.Now()
		if !t.Before(next) {
			break
		}
		sleepMin(l.recheck, next.Sub(t))
	}
	l.last = t
}

func sleepMin(a, b time.Duration) {
	if a <= b {
		time.Sleep(a)
	} else {
		time.Sleep(b)
	}
}
