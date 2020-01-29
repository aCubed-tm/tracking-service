package main

import "testing"

func TestStage1(t *testing.T) {
	// note: assuming 0,0 is bottom left
	resX, resY := 1920., 1080.
	posX, posY := 192., 720.
	fov := 135.

	// compared to manually calculated
	yaw, pitch := ScreenToYawPitch(posX, posY, resX, resY, fov)
	if !closeTo(yaw, -108) || !closeTo(pitch, 25.3125) {
		t.Errorf("Expected (-108,25.3125), got (%v,%v)", yaw, pitch)
	}

	// didn't calculate these, I just hope they're correct
	// we wouldn't use these in stage 1, we's just store the yaw and pitch
	uv := UnitVectorFromAngles(yaw, pitch)
	checkVector3f(t, uv, Vec3(-0.27934805430813975, 0.42755509343028203, -0.8597449078861006))
}
