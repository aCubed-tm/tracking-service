package main

import "log"

type CameraInfo struct {
	x float32
	y float32
	z float32
	pitch float32
	yaw float32
	fov float32
}

type CaptureInfo struct {
	x float32
	y float32
	time int
	camera CameraInfo
}

func main() {
	log.Println("Hello, world[")
}
