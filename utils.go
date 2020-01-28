package main

import (
	"math"
	"testing"
)

func checkFloatTriplets(t *testing.T, actualX float64, actualY float64, actualZ float64,
	expectedX float64, expectedY float64, expectedZ float64) {
	if !closeTo(actualX, expectedX) || !closeTo(actualY, expectedY) || !closeTo(actualZ, expectedZ) {
		t.Errorf("Expected (%v,%v,%v), got (%v,%v,%v)", expectedX, expectedY, expectedZ, actualX, actualY, actualZ)
	}
}

func closeTo(val float64, to float64) bool {
	return math.Abs(val-to) < 0.000001
}
