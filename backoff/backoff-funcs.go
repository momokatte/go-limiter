package backoff

import (
	"math"
	"math/rand"
)

/*
HalfJitter creates an exponential backoff function with half jitter, using the provided minimum and maximum values.
*/
func HalfJitter(min uint, max uint) (f func(uint) uint) {
	f = func(failCount uint) uint {
		if failCount < 1 {
			return 0
		}
		spread := pow2Uint(failCount) / 2
		return jitter(min+spread, spread, max)
	}
	return
}

/*
FullJitter creates an exponential backoff function with full jitter, using the provided minimum and maximum values.
*/
func FullJitter(min uint, max uint) (f func(uint) uint) {
	f = func(failCount uint) uint {
		if failCount < 1 {
			return 0
		}
		return jitter(min, pow2Uint(failCount), max)
	}
	return
}

/*
Pow2 is a basic exponential backoff function which returns n^x where x is the failCount.
*/
func Pow2(failCount uint) uint {
	if failCount < 1 {
		return 0
	}
	return pow2Uint(failCount)
}

/*
None is a backoff function which always returns 0.

This should only be used for testing.
*/
func None(failCount uint) uint {
	return 0
}

/*
Jitter returns a random number between base and a safe maximum.
*/
func jitter(base uint, spread uint, max uint) uint {
	if spread == 0 {
		return base
	}
	if base > max {
		return base
	}
	// make sure we stay positive when converting to signed int
	if spread > math.MaxInt64 {
		spread = math.MaxInt64
	}
	if j := base + uint(rand.Intn(int(spread))); j < max {
		return j
	}
	return max
}

/*
Calculate power of 2, but don't go over 2^63.
*/
func pow2Uint(exponent uint) uint {
	if exponent > 63 {
		exponent = 63
	}
	return 1 << exponent
}
