/*
Package limiter provides interfaces for various execution-limiting needs, along with implementations of those interfaces.
*/
package limiter

/*
TokenLimiter is the interface that wraps the AcquireToken and ReleaseToken methods, representing the use of a token mechanism to enforce concurrency limits.

AcquireToken blocks until a token can be acquired from the limiter's supply. The token must be held for the duration of the activity which needs to be limited, and then it must be passed to the ReleaseToken method without modification.

ReleaseToken notifies the limiter that the provided token (pointer and value) can be used by another goroutine. The caller must not modify the value of the token at any time, but if the token implementation is known by the caller then unmarshaling of its value is not discouraged.

The token type is a pointer to a 16-byte array (128 bits) to give limiter implementations many options with a fixed type, like 128-bit binary time.Time or UUID values.

- time.Time objects can be converted to/from []byte using the Time::MarshalBinary() and Time::UnmarshalBinary(data []byte) methods.

- Hexadecimal string UUIDs can be converted to/from []byte using the "encoding/hex" package.

- Numeric values can be converted to/from []byte using the "encoding/binary" package.
*/
type TokenLimiter interface {
	AcquireToken() (token *[16]byte)
	ReleaseToken(token *[16]byte)
}

/*
RateLimiter is the interface that wraps the CheckWait method, representing the use of a delay mechanism to enforce a rate limit.

CheckWait should be called at the beginning of the caller's action. It blocks if the limiter needs to restrict execution, otherwise it returns immediately. Restriction is typically based on consumption of a fixed rate budget, but may also be controlled by other factors.
*/
type RateLimiter interface {
	CheckWait()
}

/*
FailLimiter is the interface that wraps the CheckWait and Report methods, representing the use of a delay mechanism to enforce a rate limit and a feedback method to control it.

CheckWait should be called at the beginning of the caller's action. It blocks if the limiter needs to restrict execution, otherwise it returns immediately. Restriction is typically based on the last received status, but may also be controlled by other factors.

Report should be called at the end of the caller's action, providing the limiter with the success/fail status of the action. Failure statuses should be expected to incur rate throttling on subsequent calls to CheckWait.
*/
type FailLimiter interface {
	CheckWait()
	Report(success bool)
}

/*
TokenAndFailLimiter is the interface that simplifies the combination of a TokenLimiter and a FailLimiter, wrapping a ReleaseTokenAndReport method instead of ReleaseToken.

AcquireToken blocks until a token can be acquired from the limiter's supply, and also blocks if the limiter needs to restrict execution. The token must be held for the duration of the action which needs to be limited, and then it must be passed to the ReleaseTokenAndReport method without modification.

ReleaseTokenAndReport should be called at the end of the caller's action, notifying the limiter that the provided token (pointer and value) can be used by another goroutine and providing the limiter with the success/fail status of the action. The caller must not modify the value of the token at any time, but if the token implementation is known by the caller then unmarshaling of its value is not discouraged.

Report can be called outside the context of a rate-limited action to notify the limiter that an error has occurred and that the allowed execution rate should be throttled.
*/
type TokenAndFailLimiter interface {
	AcquireToken() (token *[16]byte)
	ReleaseTokenAndReport(token *[16]byte, success bool)
	Report(success bool)
}

/*
InvocationLimiter is the interface that wraps the Invoke method.

Invoke enforces the limiter's limits around the invocation of the passed function. The error returned by the function invocation is returned to the caller without modification, and its existence may be used by the limiter to delay the current return or subsequent invocations.
*/
type InvocationLimiter interface {
	Invoke(f func() error) error
}
