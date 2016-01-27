package limiter

import (
	"errors"
	"testing"
	"time"
)

func TestBurstRateLimiter(t *testing.T) {
	l := NewBurstRateLimiter(NewRate(1, time.Millisecond))

	start := time.Now()
	for i := 0; i < 40; i += 1 {
		l.CheckWait()
	}
	duration := time.Now().Sub(start)

	expectedMin := time.Duration(30) * time.Millisecond
	if duration < expectedMin {
		t.Errorf("Expected duration greater than %d, got %d", expectedMin, duration)
		t.FailNow()
	}

	expectedMax := time.Duration(50) * time.Millisecond
	if expectedMax < duration {
		t.Errorf("Expected duration less than %d, got %d", expectedMax, duration)
		t.FailNow()
	}
}

func TestBurstRateLimiter_Invoke(t *testing.T) {
	l := NewBurstRateLimiter(NewRate(1, time.Millisecond))

	if err := l.Invoke(func() error { return errors.New("error") }); err == nil {
		t.Error("Expected error, got nil")
	}

	if err := l.Invoke(func() error { return nil }); err != nil {
		t.Errorf("Unexpected error, got: %s", err.Error())
	}

	start := time.Now()
	for i := 0; i < 40; i += 1 {
		l.Invoke(func() error { return nil })
	}
	duration := time.Now().Sub(start)

	expectedMin := time.Duration(30) * time.Millisecond
	if duration < expectedMin {
		t.Errorf("Expected duration greater than %d, got %d", expectedMin, duration)
		t.FailNow()
	}

	expectedMax := time.Duration(50) * time.Millisecond
	if expectedMax < duration {
		t.Errorf("Expected duration less than %d, got %d", expectedMax, duration)
		t.FailNow()
	}
}

func BenchmarkBurstRateLimiter(b *testing.B) {
	l := NewBurstRateLimiter(NewRate(2000000, time.Millisecond))

	for i := 0; i < b.N; i++ {
		l.CheckWait()
	}
}

func BenchmarkBurstRateLimiter_Invoke(b *testing.B) {
	l := NewBurstRateLimiter(NewRate(2000000, time.Millisecond))

	for i := 0; i < b.N; i++ {
		l.Invoke(func() error { return nil })
	}
}
