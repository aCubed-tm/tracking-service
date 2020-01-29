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
	vec := UnitVectorFromAngles(0, 0)
	checkVector3(t, vec, Vec3(1, 0, 0))

	vec = UnitVectorFromAngles(90, 0)
	checkVector3(t, vec, Vec3(0, 0, 1))

	vec = UnitVectorFromAngles(0, -90)
	checkVector3(t, vec, Vec3(0, -1, 0))

	// haven't tested, should be right
	val := math.Pow(2, .5) / 2
	vec = UnitVectorFromAngles(45, -45)
	checkVector3(t, vec, Vec3(.5, -val, 0.5))
}

func TestCrossProduct(t *testing.T) {
	// numbers taken from: https://docs.microsoft.com/en-us/dotnet/api/system.windows.media.media3d.vector3d.crossproduct
	vec := CrossProduct(Vec3(20, 30, 40), Vec3(45, 70, 80))
	checkVector3(t, vec, Vec3(-400, 200, 50))

	// numbers taken from: https://www.mathsisfun.com/algebra/vectors-cross-product.html
	vec = CrossProduct(Vec3(2, 3, 4), Vec3(5, 6, 7))
	checkVector3(t, vec, Vec3(-3, 6, -3))
}

func TestDotProduct(t *testing.T) {
	// numbers taken from https://www.mathsisfun.com/algebra/vectors-dot-product.html under "right angles"
	res := DotProduct(Vec3(9, 2, 7), Vec3(4, 8, 10))
	if res != 122 {
		t.Errorf("Expected 122, got %v", res)
	}
}
