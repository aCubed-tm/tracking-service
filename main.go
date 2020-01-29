package main

import "log"

type CameraInfo struct {
	x     float32
	y     float32
	z     float32
	pitch float32
	yaw   float32
	fov   float32
}

type CaptureInfo struct {
	x      float32
	y      float32
	time   int
	camera CameraInfo
}

type Vector3 struct {
	x float64
	y float64
	z float64
}

func Vec3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

func main() {
	log.Println("Hello, world[")
}
