package main

import "log"

type CameraInfo struct {
	pos        Vector3
	resX, resY int16
	pitch, yaw float32
	fov        float32
}

type CaptureInfo struct {
	origin    Vector3
	direction Vector3
}

func main() {
	log.Println("Hello, world[")
}
