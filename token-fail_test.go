package limiter

import (
	"errors"
	"testing"

	"github.com/momokatte/go-backoff"
)

func TestTokenFailLimiter(t *testing.T) {
	fl := NewFailBackOffLimiter(backoff.None)
	l := NewTokenFailLimiter(NewTokenChanLimiter(1), fl)

	token := l.AcquireToken()
	l.ReleaseTokenAndReport(token, true)

	if expected := fl.failCount; expected != 0 {
		t.Errorf("Expected '%d', got '%d'", expected, 0)
	}

	token = l.AcquireToken()
	l.ReleaseTokenAndReport(token, false)

	if expected := fl.failCount; expected != 1 {
		t.Errorf("Expected '%d', got '%d'", expected, 1)
	}

}

func TestTokenFailLimiter_Invoke(t *testing.T) {
	fl := NewFailBackOffLimiter(backoff.None)
	l := NewTokenFailLimiter(NewTokenChanLimiter(1), fl)

	if err := l.Invoke(func() error { return errors.New("error") }); err == nil {
		t.Error("Expected error, got nil")
	}

	if err := l.Invoke(func() error { return nil }); err != nil {
		t.Errorf("Unexpected error, got: %s", err.Error())
	}
}

func BenchmarkTokenFailRateLimiter(b *testing.B) {
	n := b.N * 5
	tl := NewTokenChanLimiter(uint(n))
	l := NewTokenFailLimiter(tl, NewFailBackOffLimiter(backoff.None))

	tokens := make([]*[16]byte, n)
	for j := 0; j < n; j += 1 {
		tokens[j] = l.AcquireToken()
	}
	for j := 0; j < n; j += 1 {
		l.ReleaseTokenAndReport(tokens[j], true)
	}
}
