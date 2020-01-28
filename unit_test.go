package main

import (
	"math"
	"testing"
)

func TestScreenToAngle(t *testing.T) {
	yaw, pitch := ScreenToYawPitch(0, 100, 300, 200, 45)
	if !closeTo(yaw, -45) || !closeTo(pitch, 0) {
		t.Errorf("Expected (-45,0), got (%v,%v)", yaw, pitch)
	}
	yaw, pitch = ScreenToYawPitch(150, 0, 300, 200, 45)
	if !closeTo(yaw, 0) || !closeTo(pitch, -30) {
		t.Errorf("Expected (0,-30), got (%v,%v)", yaw, pitch)
	}
}

func TestUnitVector(t *testing.T) {
	x, y, z := UnitVectorFromAngles(0, 0)
	checkFloatTriplets(t, x, y, z, 1, 0, 0)

	x, y, z = UnitVectorFromAngles(90, 0)
	checkFloatTriplets(t, x, y, z, 0, 0, 1)

	x, y, z = UnitVectorFromAngles(0, -90)
	checkFloatTriplets(t, x, y, z, 0, -1, 0)

	// haven't tested, should be right
	val := math.Pow(2, .5) / 2
	x, y, z = UnitVectorFromAngles(45, -45)
	checkFloatTriplets(t, x, y, z, .5, -val, 0.5)
}

func TestCrossProduct(t *testing.T) {
	// numbers taken from: https://docs.microsoft.com/en-us/dotnet/api/system.windows.media.media3d.vector3d.crossproduct
	x, y, z := CrossProduct(20, 30, 40, 45, 70, 80)
	checkFloatTriplets(t, x, y, z, -400, 200, 50)

	// numbers taken from: https://www.mathsisfun.com/algebra/vectors-cross-product.html
	x, y, z = CrossProduct(2, 3, 4, 5, 6, 7)
	checkFloatTriplets(t, x, y, z, -3, 6, -3)
}