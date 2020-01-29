package main

import (
	"math"
	"testing"
)

func checkVector3(t *testing.T, actual, expected Vector3) {
	if !closeTo(actual.x, expected.x) || !closeTo(actual.y, expected.y) || !closeTo(actual.z, expected.z) {
		t.Errorf("Expected (%v,%v,%v), got (%v,%v,%v)", expected.x, expected.y, expected.z, actual.x, actual.y, actual.z)
	}
}

func closeTo(val, to float64) bool {
	return math.Abs(val-to) < 0.000001
}
