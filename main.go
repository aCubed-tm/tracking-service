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

func (v Vector3) Add(other Vector3) Vector3 {
	return Vector3{
		v.x + other.x,
		v.y + other.y,
		v.z + other.z,
	}
}

func (v Vector3) Sub(other Vector3) Vector3 {
	return Vector3{
		v.x - other.x,
		v.y - other.y,
		v.z - other.z,
	}
}

func (v Vector3) Mul(operand float64) Vector3 {
	return Vector3{
		v.x * operand,
		v.y * operand,
		v.z * operand,
	}
}

func (v Vector3) Div(operand float64) Vector3 {
	return Vector3{
		v.x / operand,
		v.y / operand,
		v.z / operand,
	}
}

func main() {
	log.Println("Hello, world[")
}
