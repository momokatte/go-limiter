
go-limiter
==========

A Go library for limiting execution. It provides interfaces for various limiter styles, along with implementations of those interfaces.

Limiter styles include:
- Limit concurrency via token pool
- Limit concurrency via wrapped invocation
- Enforce maximum action rate
- Throttle rate on error count


Roadmap
-------

January 2019:

- Add backoff limiter which has a maximum delay and always reduces delay upon success
- Add capacity limiter which reduces capacity on error and restores capacity on success


Online GoDoc
------------

https://godoc.org/github.com/momokatte/go-limiter
