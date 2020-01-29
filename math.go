package main

import "math"

func InterpolatePosition(cap1, cap2 CaptureInfo) Vector3 {
	panic("not yet implemented")
}

// fov is the angle covered by the widest axis, usually width
func ScreenToYawPitch(screenX, screenY, screenWidth, screenHeight, fov float64) (float64, float64) {
	halfWidth, halfHeight := screenWidth/2, screenHeight/2
	if screenWidth > screenHeight {
		ratio := screenHeight / screenWidth
		innerFov := fov * ratio
		return RelScreenToYawPitch((screenX-halfWidth)/halfWidth/ratio, (screenY-halfHeight)/halfHeight, innerFov)
	} else {
		ratio := screenWidth / screenHeight
		innerFov := fov * ratio
		return RelScreenToYawPitch((screenX-halfWidth)/halfWidth, (screenY-halfHeight)/halfHeight/ratio, innerFov)
	}
}

// screenX and screenY are coordinates relative to the "inner square" of the viewport
// this means that one of them can be bigger than 1
func RelScreenToYawPitch(screenX, screenY, fov float64) (float64, float64) {
	return screenX * fov, screenY * fov
}

// Transform yaw and pitch to unit vector, where Y is the up direction
func UnitVectorFromAngles(yaw, pitch float64) Vector3 {
	yaw = ToRadians(yaw)
	pitch = ToRadians(pitch)
	return Vec3(math.Cos(yaw)*math.Cos(pitch),
		math.Sin(pitch),
		math.Sin(yaw)*math.Cos(pitch))
}

func CrossProduct(v1, v2 Vector3) Vector3 {
	return Vec3(v1.y*v2.z-v1.z*v2.y,
		v1.z*v2.x-v1.x*v2.z,
		v1.x*v2.y-v1.y*v2.x)
}

func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func ToDegrees(radians float64) float64 {
	return radians / math.Pi * 180
}
