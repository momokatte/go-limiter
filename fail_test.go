package limiter

import (
	"errors"
	"go-limiter/backoff"
	"testing"
)

func TestFailBackOffLimiter(t *testing.T) {
	sleeps := 0
	f := func(failCount uint) uint {
		if failCount > 0 {
			sleeps += 1
		}
		return 0
	}
	l := NewFailBackOffLimiter(f)

	for i := 0; i < 1000; i += 1 {
		l.CheckWait()
		l.Report(false)
		l.CheckWait()
		l.Report(true)
	}

	if l.failCount != 0 {
		t.Errorf("Expected 0, got %d", l.failCount)
	}
	if sleeps != 1000 {
		t.Errorf("Expected 1000, got %d", sleeps)
	}
}

func TestFailBackOffLimiter_Invoke(t *testing.T) {
	sleeps := 0
	f := func(failCount uint) uint {
		if failCount > 0 {
			sleeps += 1
		}
		return 0
	}
	l := NewFailBackOffLimiter(f)

	if err := l.Invoke(func() error { return errors.New("error") }); err == nil {
		t.Error("Expected error, got nil")
	}

	if l.failCount != 1 {
		t.Errorf("Expected 1, got %d", l.failCount)
	}

	if err := l.Invoke(func() error { return nil }); err != nil {
		t.Errorf("Unexpected error, got: %s", err.Error())
	}

	if l.failCount != 0 {
		t.Errorf("Expected 0, got %d", l.failCount)
	}
	if sleeps != 1 {
		t.Errorf("Expected 1, got %d", sleeps)
	}
}

func BenchmarkBackOffLimiterSuccess(b *testing.B) {
	l := NewFailBackOffLimiter(backoff.None)
	for i := 0; i < b.N; i++ {
		l.CheckWait()
		l.Report(true)
	}
}

func BenchmarkBackOffLimiterFail(b *testing.B) {
	l := NewFailBackOffLimiter(backoff.None)
	for i := 0; i < b.N; i++ {
		l.CheckWait()
		l.Report(false)
		l.Report(true)
	}
}
