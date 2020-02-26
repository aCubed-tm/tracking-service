package main

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

func MakeCaptureInfo(posX, posY float64, camera CameraInfo) CaptureInfo {
	// take the absolute screen coordinates and screen resolution+fov to calculate yaw and pitch
	// these yaw and pitch are relative to the camera's yaw and pitch!
	yaw, pitch := ScreenToYawPitch(posX, posY, float64(camera.resX), float64(camera.resY), float64(camera.fov))

	// use camera yaw/pitch to adjust the one we calculated
	yaw = 90.0 - yaw + camera.yaw // subtracting 90 from the yaw because yaw=0 is the Y axis
	pitch += camera.pitch

	// why the fuck does go not have float modulo
	yaw = float64(int64(yaw) % int64(360.0))
	pitch = float64(int64(pitch) % int64(360.0))

	uv := UnitVectorFromAngles(yaw, pitch)

	return CaptureInfo{
		direction: uv,
		origin:    camera.pos,
	}
}

func CalculateIntersection(capture1 CaptureInfo, capture2 CaptureInfo) Vector3 {
	p1, d1 := capture1.origin, capture1.direction
	p2, d2 := capture2.origin, capture2.direction

	// assuming 2 vectors v1 and v2 defined by v=p+t*d, there is one and only one direction vector d3 that is
	// perpendicular to both v1 and v2. this vector also happens to describe the shortest distance between these other
	// vectors
	// calculate that perpendicular vector
	d3 := CrossProduct(d1, d2)

	// create a system which solves for t1, t2 and t3 for the case where v3 and v2 intersect
	cramerLHS := [9]float64{
		d1.x, -d2.x, d3.x,
		d1.y, -d2.y, d3.y,
		d1.z, -d2.z, d3.z,
	}
	rhs := Sub(p2, p1)
	cramerRHS := [3]float64{
		rhs.x,
		rhs.y,
		rhs.z,
	}
	tArr := Cramer3(cramerLHS, cramerRHS)

	// knowing t1, t2 and t3, we can calculate the points on v1 and v2 that are closest to each other
	c1 := Add(p1, Mul(d1, tArr[0]))
	c2 := Add(p2, Mul(d2, tArr[1]))

	// return their average
	return Div(Add(c1, c2), 2)
}

// fov is the angle covered by the widest axis, usually width
// this function is basically an overload for RelScreenToYawPitch that accepts screenspace coordinates
func ScreenToYawPitch(screenX, screenY, screenWidth, screenHeight, fov float64) (float64, float64) {
	halfWidth, halfHeight := screenWidth/2, screenHeight/2
	if screenWidth > screenHeight {
		ratio := screenHeight / screenWidth // < 1
		innerFov := fov * ratio
		return RelScreenToYawPitch((screenX-halfWidth)/halfWidth/ratio, (screenY-halfHeight)/halfHeight, innerFov)
	} else {
		ratio := screenWidth / screenHeight // < 1
		innerFov := fov * ratio
		return RelScreenToYawPitch((screenX-halfWidth)/halfWidth, (screenY-halfHeight)/halfHeight/ratio, innerFov)
	}
}

// screenX and screenY are coordinates relative to the "inner square" of the viewport
// this means that one of them can be bigger than 1
func RelScreenToYawPitch(screenX, screenY, fov float64) (float64, float64) {
	// division by 2 because fov goes in both directions
	return screenX * fov / 2, screenY * fov / 2
}

// Transform yaw and pitch to unit vector, where Z is up
func UnitVectorFromAngles(yaw, pitch float64) Vector3 {
	yaw = ToRadians(yaw)
	pitch = ToRadians(pitch)
	return Vec3(
		math.Cos(yaw)*math.Cos(pitch),
		math.Sin(yaw)*math.Cos(pitch),
		math.Sin(pitch))
}

// Solves a 3x4 system using Cramer's rule
func Cramer3(lhs [9]float64, rhs [3]float64) [3]float64 {
	lhsMat := mat.NewDense(3, 3, lhs[:])
	res := [3]float64{0., 0., 0.}
	buffer := make([]float64, 3)
	det := mat.Det(lhsMat)
	for c := range rhs {
		mat.Col(buffer, c, lhsMat)
		lhsMat.SetCol(c, rhs[:])
		res[c] = mat.Det(lhsMat) / det
		lhsMat.SetCol(c, buffer)
	}
	return res
}

func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
