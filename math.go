package main

import "math"

func InterpolatePosition(cap1 CaptureInfo, cap2 CaptureInfo) (float32, float32, float32) {
	panic("not yet implemented")
}

// fov is the angle covered by the widest axis, usually width
func ScreenToYawPitch(screenX float64, screenY float64, screenWidth float64, screenHeight float64, fov float64) (float64, float64) {
	halfWidth, halfHeight := screenWidth/2, screenHeight/2
	if screenWidth > screenHeight {
		ratio := screenHeight / screenWidth
		innerFov := fov * ratio
		return RelScreenToYawPitch((screenX-halfWidth)/halfWidth / ratio, (screenY-halfHeight)/halfHeight, innerFov)
	} else {
		ratio := screenWidth / screenHeight
		innerFov := fov * ratio
		return RelScreenToYawPitch((screenX-halfWidth)/halfWidth, (screenY-halfHeight)/halfHeight / ratio, innerFov)
	}
}

// screenX and screenY are coordinates relative to the "inner square" of the viewport
// this means that one of them can be bigger than 1
func RelScreenToYawPitch(screenX float64, screenY float64, fov float64) (float64, float64) {
	return screenX * fov, screenY * fov
}

// Transform yaw and pitch to unit vector, where Y is the up direction
func UnitVectorFromAngles(yaw float64, pitch float64) (float64, float64, float64) {
	yaw = ToRadians(yaw)
	pitch = ToRadians(pitch)
	return math.Cos(yaw) * math.Cos(pitch),
		math.Sin(pitch),
		math.Sin(yaw) * math.Cos(pitch)
}

// TODO: probably not needed, can remove this
func CrossProduct(x1 float64, y1 float64, z1 float64, x2 float64, y2 float64, z2 float64) (float64, float64, float64) {
	return y1*z2 - z1*y2,
		z1*x2 - x1*z2,
		x1*y2 - y1*x2
}

func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func ToDegrees(radians float64) float64 {
	return radians / math.Pi * 180
}
