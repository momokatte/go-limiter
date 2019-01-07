package limiter

import (
	"sync"
	"testing"
	"time"
)

func TestFixedIntervalLimiter(t *testing.T) {
	l := NewFixedIntervalLimiter(time.Millisecond * 10)

	var wg sync.WaitGroup
	wg.Add(4)
	start := time.Now()
	for i := 0; i < 4; i += 1 {
		go func() {
			l.CheckWait()
			wg.Done()
		}()
	}
	wg.Wait()
	duration := time.Now().Sub(start)

	expectedMin := time.Duration(30) * time.Millisecond
	if duration < expectedMin {
		t.Fatalf("Expected duration greater than %d, got %d", expectedMin, duration)
	}

	expectedMax := time.Duration(50) * time.Millisecond
	if expectedMax < duration {
		t.Fatalf("Expected duration less than %d, got %d", expectedMax, duration)
	}
}
