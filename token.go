package limiter

import (
	"sync"
)

/*
TokenChanLimiter enforces a concurrency limit using tokens and satisfies the TokenLimiter and InvocationLimiter interfaces.
*/
type TokenChanLimiter struct {
	mu     sync.Mutex
	tokens chan *[16]byte
}

/*
NewTokenChanLimiter instantiates a new TokenChanLimiter with the provided number of tokens.
*/
func NewTokenChanLimiter(initialTokens uint) (l *TokenChanLimiter) {
	l = &TokenChanLimiter{
		tokens: make(chan *[16]byte, initialTokens),
	}
	fillTokenChan(l.tokens)
	return
}

/*
AcquireToken blocks until a token can be acquired from the limiter's supply. The token must be held for the duration of the activity which needs to be limited, and then it must be passed to the ReleaseToken method without modification.
*/
func (l *TokenChanLimiter) AcquireToken() (token *[16]byte) {
	select {
	case token = <-l.tokens:
		return
	}
}

/*
ReleaseToken notifies the limiter that the provided token (pointer and value) can be used by another goroutine. The caller must not modify the value of the token at any time, but if the token implementation is known by the caller then unmarshaling of its value is not discouraged.
*/
func (l *TokenChanLimiter) ReleaseToken(token *[16]byte) {
	l.tokens <- token
}

/*
Invoke enforces the limiter's limits around the invocation of the passed function. The error returned by the function invocation is returned to the caller without modification, and its existence may be used by the limiter to delay the current return or subsequent invocations.
*/
func (l *TokenChanLimiter) Invoke(f func() error) (err error) {
	token := l.AcquireToken()
	err = f()
	l.ReleaseToken(token)
	return
}

func fillTokenChan(c chan *[16]byte) {
	capacity := cap(c)
	for i := 0; i < capacity; i += 1 {
		c <- new([16]byte)
	}
	return
}
