package main

import "testing"

func TestCalculateCaptureInfo(t *testing.T) {
	// note: assuming 0,0 is bottom left
	camera := CameraInfo{fov: 135, resX: 1920, resY: 1080}
	info := MakeCaptureInfo(192, 720, camera)

	checkVector3(t, info.direction, Vec3(-0.27934805430813975, 0.42755509343028203, -0.8597449078861006))
}

func TestCaptureToPosition(t *testing.T) {
	capture1 := CaptureInfo{
		origin:    Vec3(-1, 0, 0),
		direction: Vec3(1, 0, 0),
	}
	capture2 := CaptureInfo{
		origin:    Vec3(1, 1, 1),
		direction: Vec3(0, 0, -1),
	}

	c := CalculateIntersection(capture1, capture2)
	checkVector3(t, c, Vec3(1, 0.5, 0))
}
