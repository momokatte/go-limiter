package limiter

import (
	"github.com/momokatte/go-backoff"
)

/*
FailRateLimiter combines a FailLimiter and a RateLimiter to act as a single FailLimiter.
*/
type FailRateLimiter struct {
	failLimiter FailLimiter
	rateLimiter RateLimiter
}

/*
NewFailRateLimiter instantiates a new FailRateLimiter with the provided maximum rate and backoff function.
*/
func NewFailRateLimiter(maxRate Rate, backOff func(uint) uint) (l *FailRateLimiter) {
	l = &FailRateLimiter{}
	l.SetBackOffFunc(backOff)
	l.SetMaxRate(maxRate)
	return
}

/*
NewHalfJitterFailRateLimiter instantiates a new FailRateLimiter with the provided maximum rate, and provdes a half-jitter backoff function with the provided maximum delay.
*/
func NewHalfJitterFailRateLimiter(maxRate Rate, maxBackOff uint) (l *FailRateLimiter) {
	return NewFailRateLimiter(maxRate, backoff.HalfJitter(1, maxBackOff))
}

/*
CheckWait should be called at the beginning of the caller's action.

It blocks if the limiter needs to restrict execution, otherwise it returns immediately. Restriction is typically based on the last received status, but may also be controlled by other factors.
*/

func (l *FailRateLimiter) CheckWait() {
	l.failLimiter.CheckWait()
	l.rateLimiter.CheckWait()
	return
}

/*
Report should be called at the end of the caller's action, providing the limiter with the success/fail status of the action.

Failure statuses should be expected to incur rate throttling on subsequent calls to CheckWait.
*/
func (l *FailRateLimiter) Report(success bool) {
	l.failLimiter.Report(success)
}

/*
SetBackOffFunc sets a new backoff function for this limiter.
*/
func (l *FailRateLimiter) SetBackOffFunc(f func(uint) uint) {
	l.failLimiter = NewFailBackOffLimiter(f)
}

/*
SetRateLimit sets the maximum rate for this limiter.
*/
func (l *FailRateLimiter) SetMaxRate(rate Rate) {
	l.rateLimiter = NewBurstRateLimiter(rate)
}

/*
Invoke enforces this limiter's limits before the invocation of the provided function and uses the function's return value to adjust the backoff rate for subsequent invocations.
*/
func (l *FailRateLimiter) Invoke(f func() error) (err error) {
	l.CheckWait()
	err = f()
	l.Report(err == nil)
	return
}
