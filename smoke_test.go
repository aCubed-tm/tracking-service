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
	// we wouldn't use these in stage 1, we'd just store the yaw and pitch
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

	cramerLHS := [9]float64{
		d1.x, -d2.x, n1.x,
		d1.y, -d2.y, n1.y,
		d1.z, -d2.z, n1.z,
	}
	rhs := Sub(p2, p1)
	cramerRHS := [3]float64{
		rhs.x,
		rhs.y,
		rhs.z,
	}

	tArr := Cramer3(cramerLHS, cramerRHS)
	c1 := Add(p1, Mul(d1, tArr[0]))
	c2 := Add(p2, Mul(d2, tArr[1]))
	checkVector3(t, c1, Vec3(1, 0, 0))
	checkVector3(t, c2, Vec3(1, 1, 0))

	c := Div(Add(c1, c2), 2)

	checkVector3(t, c, Vec3(1, 0.5, 0))
}
