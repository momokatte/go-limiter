package backoff

import (
	"math"
	"testing"
)

func TestPow2(t *testing.T) {

	args := []uint{0, 1, 2, 3, 4}
	expectedValues := []uint{0, 2, 4, 8, 16}

	for i, expected := range expectedValues {
		if actual := Pow2(args[i]); expected != actual {
			t.Errorf("Expected '%d', got '%d'", expected, actual)
		}
	}
}

func TestHalfJitter(t *testing.T) {
	f := HalfJitter(0, math.MaxInt64)

	expected := uint(0)
	if actual := f(0); expected != actual {
		t.Errorf("Expected '%d', got '%d'", expected, actual)
	}

	args := []uint{1, 2, 3, 4, 10, 60, 100}
	expectedMins := []uint{1, 2, 4, 8, 512, 180000, 180000}
	expectedCaps := []uint{2, 4, 8, 16, 1024, 1152921504606846976, math.MaxInt64}

	for i, expectedCap := range expectedCaps {
		for j := 0; j < 1000; j += 1 {
			expectedMin := expectedMins[i]
			actual := f(args[i])
			if actual < expectedMin {
				t.Errorf("Expected value above '%d', got '%d'", expectedMin, actual)
				t.FailNow()
			}
			if expectedCap < actual {
				t.Errorf("Expected value below '%d', got '%d'", expectedCap, actual)
				t.FailNow()
			}
		}
	}
}

func TestFullJitter(t *testing.T) {
	f := FullJitter(0, math.MaxInt64)

	expected := uint(0)
	if actual := f(0); expected != actual {
		t.Errorf("Expected '%d', got '%d'", expected, actual)
	}

	args := []uint{1, 2, 3, 4, 10, 60, 100}
	expectedCaps := []uint{20, 40, 160, 320, 10240, 1152921504606846976, math.MaxInt64}

	for i, expectedCap := range expectedCaps {
		for j := 0; j < 1000; j += 1 {
			if actual := f(args[i]); expectedCap < actual {
				t.Errorf("Expected value below '%d', got '%d'", expectedCap, actual)
				t.FailNow()
			}
		}
	}
}

func TestPow2Uint(t *testing.T) {

	args := []uint{0, 1, 2, 3, 4, 10, 60, 100}
	expectedValues := []uint{1, 2, 4, 8, 16, 1024, 1152921504606846976, 9223372036854775808}

	for i, expected := range expectedValues {
		if actual := pow2Uint(args[i]); expected < actual {
			t.Errorf("Expected '%d', got '%d'", expected, actual)
		}
	}
}
