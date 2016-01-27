package limiter

import (
	"go-limiter/backoff"
	"time"

	vrl "github.com/youtube/vitess/go/ratelimiter"
)

/*
BurstRateLimiter enforces a rate limit within an interval and satisfies the RateLimiter and InvocationLimiter interfaces.
*/
type BurstRateLimiter struct {
	rateLimiter *vrl.RateLimiter
	backOffFunc func(uint) uint
}

/*
NewBurstRateLimiter instantiates a BurstRateLimiter with the provided rate threshold and a wait-backoff function with full jitter appropriate for high-frequency use (more than 200 actions per second).
*/
func NewBurstRateLimiter(maxRate Rate) (l *BurstRateLimiter) {
	l = &BurstRateLimiter{}
	l.SetMaxRate(maxRate)
	// jitter smooths out retries -- this minimum value is for high-frequency use
	// TODO: calculate minimum based on provided rate
	l.backOffFunc = backoff.FullJitter(uint(time.Millisecond/2), uint(maxRate.Duration))
	return
}

/*
CheckWait should be called at the beginning of the caller's action. It blocks if the limiter needs to restrict execution, otherwise it returns immediately. Restriction is typically based on consumption of a fixed rate budget, but may also be controlled by other factors.

If calls to this method are uniform, the allowed rate will roughly match the rate threshold. Non-uniform use may result in a rate during the end of an interval and the beginning of the subsequent interval which together exceed the specified threshold.
*/
func (l *BurstRateLimiter) CheckWait() {
	// retry with backoff until allowed
	for fails := uint(0); !l.rateLimiter.Allow(); {
		fails += 1
		sleep := l.backOffFunc(fails)
		time.Sleep(time.Duration(sleep) * time.Nanosecond)
	}
	return
}

/*
Invoke enforces the limiter's limits around the invocation of the passed function. The error returned by the function invocation is returned to the caller without modification, and its existence may be used by the limiter to delay the current return or subsequent invocations.
*/
func (l *BurstRateLimiter) Invoke(f func() error) (err error) {
	l.CheckWait()
	err = f()
	return
}

/*
SetRateLimit sets a new rate threshold for this limiter.

It does not make any adjustment based on the previous threshold, so the rate allowed immediately before and after this change may together exceed the specified threshold.
*/
func (l *BurstRateLimiter) SetMaxRate(rate Rate) {
	l.rateLimiter = vrl.NewRateLimiter(rate.Count, rate.Duration)
}
