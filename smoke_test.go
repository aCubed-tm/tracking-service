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
	checkVector3(t, uv, Vec3(-0.27934805430813975, 0.42755509343028203, -0.8597449078861006))
}

func TestStage2(t *testing.T) {
	p1, d1 := Vec3(-1, 0, 0), Vec3(1, 0, 0) // from -x, points to 0,0,0
	p2, d2 := Vec3(1, 1, 1), Vec3(0, 0, -1) // from +x+y+z, points to -z
	// expected point should be at 1,0.5,0

	// calculate the perpendicular line
	n1 := CrossProduct(d1, d2)
	checkVector3(t, n1, Vec3(0, 1, 0))
	n2 := CrossProduct(d2, n1) // 1,0,0 not sure what exactly this is. causes problems later on

	c1 := p1.Add(d1.Mul(DotProduct(p2.Sub(p1), n2) / DotProduct(d1, n2)))
	c2 := p2.Add(d2.Mul(DotProduct(p1.Sub(p2), n1) / DotProduct(d2, n1)))

	checkVector3(t, c1, Vec3(1, 0, 0))
	checkVector3(t, c2, Vec3(1, 1, 0))

	c := c1.Add(c2).Div(2)

	checkVector3(t, c, Vec3(1, 0.5, 0))
}
