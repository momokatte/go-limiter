package limiter

/*
TokenFailLimiter combines a TokenLimiter and a FailLimiter to satisfy the TokenAndFailLimiter interface.
*/
type TokenFailLimiter struct {
	tokenLimiter TokenLimiter
	failLimiter  FailLimiter
}

/*
NewTokenFailLimiter instantiates a new TokenFailLimiter with the provided TokenLimiter and FailLimiter.
*/
func NewTokenFailLimiter(tl TokenLimiter, fl FailLimiter) (l *TokenFailLimiter) {
	l = &TokenFailLimiter{
		tokenLimiter: tl,
		failLimiter:  fl,
	}
	return
}

/*
AcquireToken blocks until a token can be acquired from the limiter's supply, and also blocks if the limiter needs to restrict execution. The token must be held for the duration of the action which needs to be limited, and then it must be passed to the ReleaseTokenAndReport method without modification.
*/
func (l *TokenFailLimiter) AcquireToken() (token *[16]byte) {
	token = l.tokenLimiter.AcquireToken()
	l.failLimiter.CheckWait()
	return
}

/*
ReleaseTokenAndReport should be called at the end of the caller's action, notifying the limiter that the provided token (pointer and value) can be used by another goroutine and providing the limiter with the success/fail status of the action. The caller must not modify the value of the token at any time, but if the token implementation is known by the caller then unmarshaling of its value is not discouraged.
*/
func (l *TokenFailLimiter) ReleaseTokenAndReport(token *[16]byte, success bool) {
	l.failLimiter.Report(success)
	l.tokenLimiter.ReleaseToken(token)
}

/*
Report can be called outside the context of a rate-limited action to notify the limiter that an error has occurred and that the allowed execution rate should be throttled.
*/
func (l *TokenFailLimiter) Report(success bool) {
	l.failLimiter.Report(success)
}

/*
Invoke enforces the limiter's limits before the invocation of the passed function. The error returned by the function invocation is returned to the caller without modification, and its existence may be used by the limiter to delay the current return or subsequent invocations.
*/
func (l *TokenFailLimiter) Invoke(f func() error) (err error) {
	token := l.AcquireToken()
	err = f()
	l.ReleaseTokenAndReport(token, err == nil)
	return
}
