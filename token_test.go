package limiter

import (
	"errors"
	"testing"
)

func TestTokenChanLimiter(t *testing.T) {
	l := NewTokenChanLimiter(1)

	token := l.AcquireToken()
	l.ReleaseToken(token)
}

func TestTokenChanLimiter_Invoke(t *testing.T) {
	l := NewTokenChanLimiter(1)

	if err := l.Invoke(func() error { return errors.New("error") }); err == nil {
		t.Error("Expected error, got nil")
	}

	if err := l.Invoke(func() error { return nil }); err != nil {
		t.Errorf("Unexpected error, got: %s", err.Error())
	}
}

func BenchmarkTokenChanLimiterSingle(b *testing.B) {
	l := NewTokenChanLimiter(1)
	for i := 0; i < b.N; i++ {
		token := l.AcquireToken()
		l.ReleaseToken(token)
	}
}

func BenchmarkTokenChanLimiterMulti(b *testing.B) {
	n := b.N * 5
	l := NewTokenChanLimiter(uint(n))
	tokens := make([]*[16]byte, n)
	for j := 0; j < n; j += 1 {
		tokens[j] = l.AcquireToken()
	}
	for j := 0; j < n; j += 1 {
		l.ReleaseToken(tokens[j])
	}
}
