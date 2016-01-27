package limiter

import (
	"sync"
	"time"
)

/*
FailBackOffLimiter delays execution of subsequent invocations when the caller reports failure conditions, and satisfies the FailLimiter and InvocationLimiter interfaces.
*/
type FailBackOffLimiter struct {
	mu          sync.Mutex
	failCount   uint
	backOffFunc func(uint) uint
}

/*
NewFailBackOffLimiter instantiates a new FailBackOffLimiter with the provided backoff function.

This package provides backoff function builders in go-limiter/backoff.
*/
func NewFailBackOffLimiter(backOffFunc func(uint) uint) (l *FailBackOffLimiter) {
	l = &FailBackOffLimiter{
		backOffFunc: backOffFunc,
	}
	return
}

/*
CheckWait should be called at the beginning of the caller's action.

It blocks if the limiter needs to restrict execution, otherwise it returns immediately. Restriction is typically based on the last received status, but may also be controlled by other factors.
*/
func (l *FailBackOffLimiter) CheckWait() {
	if l.failCount == 0 {
		return
	}
	if sleep := l.backOffFunc(l.failCount); sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

/*
Report should be called at the end of the caller's action, providing the limiter with the success/fail status of the action.

Failure statuses should be expected to incur rate throttling on subsequent calls to CheckWait.
*/
func (l *FailBackOffLimiter) Report(success bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if success && l.failCount > 0 {
		l.failCount -= 1
	} else if !success {
		l.failCount += 1
	}
}

/*
Invoke enforces the limiter's limits around the invocation of the passed function.

The error returned by the function invocation is returned to the caller without modification, and its existence may be used by the limiter to delay the current return or subsequent invocations.
*/
func (l *FailBackOffLimiter) Invoke(f func() error) (err error) {
	l.CheckWait()
	err = f()
	l.Report(err == nil)
	return
}
