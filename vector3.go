package main

func Vec3(x, y, z float64) Vector3 {
	return Vector3{x, y, z}
}

func Add(v1, v2 Vector3) Vector3 {
	return Vector3{
		v1.x + v2.x,
		v1.y + v2.y,
		v1.z + v2.z,
	}
}

func Sub(v1, v2 Vector3) Vector3 {
	return Vector3{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
	}
}

func Mul(v Vector3, operand float64) Vector3 {
	return Vector3{
		v.x * operand,
		v.y * operand,
		v.z * operand,
	}
}

func Div(v Vector3, operand float64) Vector3 {
	return Vector3{
		v.x / operand,
		v.y / operand,
		v.z / operand,
	}
}

func CrossProduct(v1, v2 Vector3) Vector3 {
	return Vec3(v1.y*v2.z-v1.z*v2.y,
		v1.z*v2.x-v1.x*v2.z,
		v1.x*v2.y-v1.y*v2.x)
}

func DotProduct(v1, v2 Vector3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}
