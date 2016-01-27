
# go-limiter

Package limiter provides interfaces for various execution-limiting needs, along with implementations of those interfaces.

Limiter styles include:
- Limit concurrency via token pool
- Limit concurrency via wrapped invocation
- Enforce maximum action rate
- Throttle rate on error count

This package also provides builders for half-jitter and full-jitter exponential backoff functions, which can smooth out action retries. This AWS blog post demonstrates the benefits of adding jitter to backoff behavior: https://www.awsarchitectureblog.com/2015/03/backoff.html


## Online GoDoc

https://godoc.org/github.com/momokatte/go-limiter
