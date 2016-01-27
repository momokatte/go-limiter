package limiter

import (
	"errors"
	"testing"
)

func TestAdjustableTokenChanLimiter(t *testing.T) {
	l := NewAdjustableTokenChanLimiter(1, 8)

	if actual := l.GetTokenCount(); actual != 1 {
		t.Errorf("%d", actual)
		t.FailNow()
	}

	token := l.AcquireToken()
	l.ReleaseToken(token)

	l.AddTokens(1)

	if actual := l.GetTokenCount(); actual != 2 {
		t.Errorf("%d", actual)
		t.FailNow()
	}

	l.RemoveTokens(2)

	if actual := l.GetTokenCount(); actual != 0 {
		t.Errorf("%d", actual)
		t.FailNow()
	}
}

func TestAdjustableCpuTokenChanLimiter(t *testing.T) {
	l := NewAdjustableCpuTokenChanLimiter(1)

	if actual := l.GetTokenCount(); actual != 1 {
		t.Errorf("%d", actual)
		t.FailNow()
	}

	token := l.AcquireToken()
	l.ReleaseToken(token)
}

func TestAdjustableTokenChanLimiter_Invoke(t *testing.T) {
	l := NewAdjustableTokenChanLimiter(1, 8)

	if err := l.Invoke(func() error { return errors.New("error") }); err == nil {
		t.Error("Expected error, got nil")
	}

	if err := l.Invoke(func() error { return nil }); err != nil {
		t.Errorf("Unexpected error, got: %s", err.Error())
	}
}

func BenchmarkAdjustableTokenChanLimiterSingle(b *testing.B) {
	l := NewAdjustableTokenChanLimiter(1, 1)
	for i := 0; i < b.N; i++ {
		var token *[16]byte = l.AcquireToken()
		l.ReleaseToken(token)
	}
}

func BenchmarkAdjustableTokenChanLimiterMulti(b *testing.B) {
	n := b.N * 5
	l := NewAdjustableTokenChanLimiter(uint(n), uint(n))
	tokens := make([]*[16]byte, n)
	for j := 0; j < n; j += 1 {
		tokens[j] = l.AcquireToken()
	}
	for j := 0; j < n; j += 1 {
		l.ReleaseToken(tokens[j])
	}
}
