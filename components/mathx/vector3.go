package mathx

import (
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

type Vector3[T constraints.Float] struct {
	X, Y, Z T
}

func NewVector3[T constraints.Float](x, y, z T) *Vector3[T] {
	return &Vector3[T]{X: x, Y: y, Z: z}
}

func NewZeroVector3[T constraints.Float]() *Vector3[T] {
	return NewVector3[T](0, 0, 0)
}

func NewOneVector3[T constraints.Float]() *Vector3[T] {
	return NewVector3[T](1, 1, 1)
}

var (
	V3Zero32 = *NewZeroVector3[float32]()
	V3One32  = *NewOneVector3[float32]()
	V3Zero64 = *NewZeroVector3[float64]()
	V3One64  = *NewOneVector3[float64]()
)

func (v *Vector3[T]) Set(x, y, z T) *Vector3[T] {
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

func (v *Vector3[T]) SetScalar(scalar T) *Vector3[T] {
	v.X = scalar
	v.Y = scalar
	v.Z = scalar
	return v
}

func (v *Vector3[T]) SetComponent(index int, value T) *Vector3[T] {
	switch index {
	case 0:
		v.X = value
	case 1:
		v.Y = value
	case 2:
		v.Z = value
	default:
		panic("index is out of range")
	}
	return v
}

func (v *Vector3[T]) Component(index int) T {
	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	case 2:
		return v.Z
	default:
		panic("index is out of range")
	}
}

func (v *Vector3[T]) Clone() *Vector3[T] {
	return NewVector3(v.X, v.Y, v.Z)
}

func (v *Vector3[T]) Copy(vector Vector3[T]) *Vector3[T] {
	v.X = vector.X
	v.Y = vector.Y
	v.Z = vector.Z
	return v
}

func (v *Vector3[T]) Add(v2 Vector3[T]) *Vector3[T] {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

func (v *Vector3[T]) AddScalar(scalar T) *Vector3[T] {
	v.X += scalar
	v.Y += scalar
	v.Z += scalar
	return v
}

func AddVector3s[T constraints.Float](a, b Vector3[T]) *Vector3[T] {
	return NewVector3(a.X+b.X, a.Y+b.Y, a.Z+b.Z)
}

func (v *Vector3[T]) AddScaledVector(vector Vector3[T], scalar T) *Vector3[T] {
	v.X += vector.X * scalar
	v.Y += vector.Y * scalar
	v.Z += vector.Z * scalar
	return v
}

func (v *Vector3[T]) Sub(v2 Vector3[T]) *Vector3[T] {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}

func (v *Vector3[T]) SubScalar(scalar T) *Vector3[T] {
	v.X -= scalar
	v.Y -= scalar
	v.Z -= scalar
	return v
}

func SubVector3s[T constraints.Float](a, b Vector3[T]) *Vector3[T] {
	return NewVector3(a.X-b.X, a.Y-b.Y, a.Z-b.Z)
}

func (v *Vector3[T]) Multiply(v2 Vector3[T]) *Vector3[T] {
	v.X *= v2.X
	v.Y *= v2.Y
	v.Z *= v2.Z
	return v
}

func (v *Vector3[T]) MultiplyScalar(scalar T) *Vector3[T] {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	return v
}

func MultiplyVector3s[T constraints.Float](a, b Vector3[T]) *Vector3[T] {
	return NewVector3(a.X*b.X, a.Y*b.Y, a.Z*b.Z)
}

func (v *Vector3[T]) ApplyEuler(e Euler[T]) *Vector3[T] {
	return v.ApplyQuaternion(*NewIdentityQuaternion[T]().SetFromEuler(e))
}

func (v *Vector3[T]) ApplyAxisAngle(axis Vector3[T], angle T) *Vector3[T] {
	return v.ApplyQuaternion(*NewIdentityQuaternion[T]().SetFromAxisAngle(axis, angle))
}

func (v *Vector3[T]) ApplyMatrix3(m Matrix3[T]) *Vector3[T] {
	x, y, z := v.X, v.Y, v.Z
	v.X = m[0]*x + m[3]*y + m[6]*z
	v.Y = m[1]*x + m[4]*y + m[7]*z
	v.Z = m[2]*x + m[5]*y + m[8]*z
	return v
}

func (v *Vector3[T]) ApplyNormalMatrix(m Matrix3[T]) *Vector3[T] {
	return v.ApplyMatrix3(m).Normalize()
}

func (v *Vector3[T]) ApplyMatrix4(m Matrix4[T]) *Vector3[T] {
	x, y, z := v.X, v.Y, v.Z
	w := 1 / (m[3]*x + m[7]*y + m[11]*z + m[15])

	v.X = (m[0]*x + m[4]*y + m[8]*z + m[12]) * w
	v.Y = (m[1]*x + m[5]*y + m[9]*z + m[13]) * w
	v.Z = (m[2]*x + m[6]*y + m[10]*z + m[14]) * w
	return v
}

func (v *Vector3[T]) ApplyQuaternion(q Quaternion[T]) *Vector3[T] {
	vx, vy, vz := v.X, v.Y, v.Z
	qx, qy, qz, qw := q.X, q.Y, q.Z, q.W

	tx := 2 * (qy*vz - qz*vy)
	ty := 2 * (qz*vx - qx*vz)
	tz := 2 * (qx*vy - qy*vx)

	v.X = vx + qw*tx + qy*tz - qz*ty
	v.Y = vy + qw*ty + qz*tx - qx*tz
	v.Z = vz + qw*tz + qx*ty - qy*tx
	return v
}

func (v *Vector3[T]) TransformDirection(m Matrix4[T]) *Vector3[T] {
	x, y, z := v.X, v.Y, v.Z
	v.X = m[0]*x + m[4]*y + m[8]*z
	v.Y = m[1]*x + m[5]*y + m[9]*z
	v.Z = m[2]*x + m[6]*y + m[10]*z
	return v.Normalize()
}

func (v *Vector3[T]) Divide(v2 Vector3[T]) *Vector3[T] {
	v.X /= v2.X
	v.Y /= v2.Y
	v.Z /= v2.Z
	return v
}

func (v *Vector3[T]) DivideScalar(scalar T) *Vector3[T] {
	return v.MultiplyScalar(1 / scalar)
}

func (v *Vector3[T]) Min(v2 Vector3[T]) *Vector3[T] {
	v.X = min(v.X, v2.X)
	v.Y = min(v.Y, v2.Y)
	v.Z = min(v.Z, v2.Z)
	return v
}

func (v *Vector3[T]) Max(v2 Vector3[T]) *Vector3[T] {
	v.X = max(v.X, v2.X)
	v.Y = max(v.Y, v2.Y)
	v.Z = max(v.Z, v2.Z)
	return v
}

func (v *Vector3[T]) Clamp(minVal, maxVal Vector3[T]) *Vector3[T] {
	v.X = max(minVal.X, min(maxVal.X, v.X))
	v.Y = max(minVal.Y, min(maxVal.Y, v.Y))
	v.Z = max(minVal.Z, min(maxVal.Z, v.Z))
	return v
}

func (v *Vector3[T]) ClampScalar(minVal, maxVal T) *Vector3[T] {
	v.X = max(minVal, min(maxVal, v.X))
	v.Y = max(minVal, min(maxVal, v.Y))
	v.Z = max(minVal, min(maxVal, v.Z))
	return v
}

func (v *Vector3[T]) ClampLength(minVal, maxVal T) *Vector3[T] {
	length := v.Length()
	return v.DivideScalar(length).MultiplyScalar(max(minVal, min(maxVal, length)))
}

func (v *Vector3[T]) Floor() *Vector3[T] {
	v.X = T(math.Floor(float64(v.X)))
	v.Y = T(math.Floor(float64(v.Y)))
	v.Z = T(math.Floor(float64(v.Z)))
	return v
}

func (v *Vector3[T]) Ceil() *Vector3[T] {
	v.X = T(math.Ceil(float64(v.X)))
	v.Y = T(math.Ceil(float64(v.Y)))
	v.Z = T(math.Ceil(float64(v.Z)))
	return v
}

func (v *Vector3[T]) Round() *Vector3[T] {
	v.X = T(math.Round(float64(v.X)))
	v.Y = T(math.Round(float64(v.Y)))
	v.Z = T(math.Round(float64(v.Z)))
	return v
}

func (v *Vector3[T]) RoundToZero() *Vector3[T] {
	v.X = T(math.Trunc(float64(v.X)))
	v.Y = T(math.Trunc(float64(v.Y)))
	v.Z = T(math.Trunc(float64(v.Z)))
	return v
}

func (v *Vector3[T]) Negate() *Vector3[T] {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

func (v *Vector3[T]) Dot(v2 Vector3[T]) T {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v *Vector3[T]) LengthSq() T {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vector3[T]) Length() T {
	return T(math.Sqrt(float64(v.LengthSq())))
}

func (v *Vector3[T]) Normalize() *Vector3[T] {
	l := v.Length()
	if l == 0 {
		v.X = 0
		v.Y = 0
		v.Z = 0
	} else {
		l = 1 / l
		v.X *= l
		v.Y *= l
		v.Z *= l
	}
	return v
}

func (v *Vector3[T]) SetLength(length T) *Vector3[T] {
	return v.Normalize().MultiplyScalar(length)
}

func (v *Vector3[T]) Lerp(v2 Vector3[T], alpha T) *Vector3[T] {
	v.X += (v2.X - v.X) * alpha
	v.Y += (v2.Y - v.Y) * alpha
	v.Z += (v2.Z - v.Z) * alpha
	return v
}

func LerpVector3s[T constraints.Float](v1, v2 Vector3[T], alpha T) *Vector3[T] {
	return NewVector3(
		v1.X+(v2.X-v1.X)*alpha,
		v1.Y+(v2.Y-v1.Y)*alpha,
		v1.Z+(v2.Z-v1.Z)*alpha,
	)
}

func (v *Vector3[T]) Cross(v2 Vector3[T]) *Vector3[T] {
	ax, ay, az := v.X, v.Y, v.Z
	bx, by, bz := v2.X, v2.Y, v2.Z

	v.X = ay*bz - az*by
	v.Y = az*bx - ax*bz
	v.Z = ax*by - ay*bx

	return v
}

func CrossVector3s[T constraints.Float](a, b Vector3[T]) *Vector3[T] {
	return a.Clone().Cross(b)
}

func (v *Vector3[T]) ProjectOnVector(v2 Vector3[T]) *Vector3[T] {
	denominator := v2.LengthSq()
	if denominator == 0 {
		return v.Set(0, 0, 0)
	}
	scalar := v.Dot(v2) / denominator
	return v.Copy(v2).MultiplyScalar(scalar)
}

func (v *Vector3[T]) ProjectOnPlane(planeNormal Vector3[T]) *Vector3[T] {
	return v.Sub(*v.Clone().ProjectOnVector(planeNormal))
}

func (v *Vector3[T]) Reflect(normal Vector3[T]) *Vector3[T] {
	return v.Sub(*v.Clone().ProjectOnPlane(normal).MultiplyScalar(2))
}

func (v *Vector3[T]) AngleTo(v2 Vector3[T]) T {
	denominator := math.Sqrt(float64(v.LengthSq() * v2.LengthSq()))
	if denominator == 0 {
		return math.Pi / 2
	}

	theta := v.Dot(v2) / T(denominator)
	return T(math.Acos(float64(Clamp(theta, -1, 1))))
}

func (v *Vector3[T]) DistanceTo(v2 Vector3[T]) T {
	return T(math.Sqrt(float64(v.DistanceToSquared(v2))))
}

func (v *Vector3[T]) DistanceToSquared(v2 Vector3[T]) T {
	dx := v.X - v2.X
	dy := v.Y - v2.Y
	dz := v.Z - v2.Z
	return dx*dx + dy*dy + dz*dz
}

func (v *Vector3[T]) ManhattanDistanceTo(v2 Vector3[T]) T {
	return T(math.Abs(float64(v.X-v2.X)) + math.Abs(float64(v.Y-v2.Y)) + math.Abs(float64(v.Z-v2.Z)))
}

func (v *Vector3[T]) SetFromSpherical(s Spherical[T]) *Vector3[T] {
	return v.SetFromSphericalCoords(s.Radius, s.Phi, s.Theta)
}

func (v *Vector3[T]) SetFromSphericalCoords(radius, phi, theta T) *Vector3[T] {
	rf, pf, tf := float64(radius), float64(phi), float64(theta)
	sinPhiRadius := math.Sin(pf) * rf
	v.X = T(sinPhiRadius * math.Sin(tf))
	v.Y = T(math.Cos(pf) * rf)
	v.Z = T(sinPhiRadius * math.Cos(tf))
	return v
}

func (v *Vector3[T]) SetFromCylindrical(c Cylindrical[T]) *Vector3[T] {
	return v.SetFromCylindricalCoords(c.Radius, c.Theta, c.Y)
}

func (v *Vector3[T]) SetFromCylindricalCoords(radius, theta, y T) *Vector3[T] {
	v.X = radius * T(math.Sin(float64(theta)))
	v.Y = y
	v.Z = radius * T(math.Cos(float64(theta)))
	return v
}

func (v *Vector3[T]) SetFromMatrixPosition(m Matrix4[T]) *Vector3[T] {
	return v.Set(m[12], m[13], m[14])
}

func (v *Vector3[T]) SetFromMatrixScale(m Matrix4[T]) *Vector3[T] {
	_, _, scale := m.Decompose()
	return v.Copy(scale)
}

func (v *Vector3[T]) SetFromMatrixColumn(m Matrix4[T], index int) *Vector3[T] {
	return v.FromArray(m[:], index*4)
}

func (v *Vector3[T]) Equals(v2 Vector3[T]) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}

func (v *Vector3[T]) FromArray(array []T, offset int) *Vector3[T] {
	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
	return v
}

func (v *Vector3[T]) ToArray(array []T, offset int) []T {
	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
	return array
}

func (v *Vector3[T]) Random() *Vector3[T] {
	v.X = T(rand.Float64())
	v.Y = T(rand.Float64())
	v.Z = T(rand.Float64())
	return v
}

func (v *Vector3[T]) RandomDirection() *Vector3[T] {
	u := (rand.Float64() - 0.5) * 2
	t := rand.Float64() * math.Pi * 2
	f := math.Sqrt(1 - u*u)
	v.X = T(f * math.Cos(t))
	v.Y = T(f * math.Sin(t))
	v.Z = T(u)
	return v
}
