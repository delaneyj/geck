package mathx

import (
	"math"
	"math/rand"
)

const (
	DEG2RAD   = math.Pi / 180
	RAD2DEG   = 180 / math.Pi
	EPSILON64 = float64(7.)/3 - float64(4.)/3 - float64(1.)
	EPSILON32 = float32(7.)/3 - float32(4.)/3 - float32(1.)
)

type CoordinateSystem int

const (
	CoordinateSystemWebGL CoordinateSystem = iota
	CoordinateSystemWebGPU
)

var (
	coordinateSystem = CoordinateSystemWebGL
)

func Clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(max, value))
}

// compute euclidean modulo of m % n
// https://en.wikipedia.org/wiki/Modulo_operation
func EuclideanModulo(n, m float64) float64 {
	return math.Mod(math.Mod(n, m)+m, m)
}

// Linear mapping from range <a1, a2> to range <b1, b2>
func MapLinear(x, a1, a2, b1, b2 float64) float64 {
	return b1 + (x-a1)*(b2-b1)/(a2-a1)
}

// https://www.gamedev.net/tutorials/programming/general-and-gameplay-programming/inverse-lerp-a-super-useful-yet-often-overlooked-function-r5230/
func InverseLerp(x, y, value float64) float64 {
	if x != y {
		return (value - x) / (y - x)
	} else {
		return 0
	}
}

// https://en.wikipedia.org/wiki/Linear_interpolation
func Lerp(x, y, t float64) float64 {
	return (1-t)*x + t*y
}

// http://www.rorydriscoll.com/2016/03/07/frame-rate-independent-damping-using-lerp/
func Damp(x, y, lambda, dt float64) float64 {
	return Lerp(x, y, 1-math.Exp(-lambda*dt))
}

// https://www.desmos.com/calculator/vcsjnyz7x4
func PingPong(x, length float64) float64 {
	return length - math.Abs(EuclideanModulo(x, length*2)-length)
}

// http://en.wikipedia.org/wiki/Smoothstep
func Smoothstep(x, min, max float64) float64 {
	if x <= min {
		return 0
	}
	if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return x * x * (3 - 2*x)
}

func Smootherstep(x, min, max float64) float64 {
	if x <= min {
		return 0
	}
	if x >= max {
		return 1
	}
	x = (x - min) / (max - min)
	return x * x * x * (x*(x*6-15) + 10)
}

// Random float from <-range/2, range/2> interval
func RandFloatSpread(rng float64) float64 {
	return rng * (0.5 - rand.Float64())
}

func DegToRad(degrees float64) float64 {
	return degrees * DEG2RAD
}

func RadToDeg(radians float64) float64 {
	return radians * RAD2DEG
}

func IsPowerOfTwo(value int) bool {
	return value != 0 && (value&(value-1)) == 0
}

func CeilPowerOfTwo(value int) int {
	return int(math.Pow(2, math.Ceil(math.Log(float64(value))/math.Ln2)))
}

func FloorPowerOfTwo(value int) int {
	return int(math.Pow(2, math.Floor(math.Log(float64(value))/math.Ln2)))
}

func SetQuaternionFromProperEuler(q *Quaternion, a, b, c float64, order string) {
	// Intrinsic Proper Euler Angles - see https://en.wikipedia.org/wiki/Euler_angles

	// rotations are applied to the axes in the order specified by 'order'
	// rotation by angle 'a' is applied first, then by angle 'b', then by angle 'c'
	// angles are in radians

	cos := math.Cos
	sin := math.Sin

	c2 := cos(b / 2)
	s2 := sin(b / 2)

	c13 := cos((a + c) / 2)
	s13 := sin((a + c) / 2)

	c1_3 := cos((a - c) / 2)
	s1_3 := sin((a - c) / 2)

	c3_1 := cos((c - a) / 2)
	s3_1 := sin((c - a) / 2)

	switch order {
	case "XYX":
		q.Set(c2*s13, s2*c1_3, s2*s1_3, c2*c13)
	case "YZY":
		q.Set(s2*s1_3, c2*s13, s2*c1_3, c2*c13)
	case "ZXZ":
		q.Set(s2*c1_3, s2*s1_3, c2*s13, c2*c13)
	case "XZX":
		q.Set(c2*s13, s2*s3_1, s2*c3_1, c2*c13)
	case "YXY":
		q.Set(s2*c3_1, c2*s13, s2*s3_1, c2*c13)
	case "ZYZ":
		q.Set(s2*s3_1, s2*c3_1, c2*s13, c2*c13)
	default:
		panic("Invalid order")
	}
}
