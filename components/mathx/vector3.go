package mathx

import (
	"math"
	"math/rand"
)

type Vector3 struct {
	X, Y, Z float64
}

func NewVector3(x, y, z float64) *Vector3 {
	return &Vector3{X: x, Y: y, Z: z}
}

func NewZeroVector3() *Vector3 {
	return NewVector3(0, 0, 0)
}

func NewOneVector3() *Vector3 {
	return NewVector3(1, 1, 1)
}

var (
	V3Zero = *NewZeroVector3()
	V3One  = *NewOneVector3()
)

func (v *Vector3) Set(x, y, z float64) *Vector3 {
	v.X = x
	v.Y = y
	v.Z = z
	return v
}

func (v *Vector3) SetScalar(scalar float64) *Vector3 {
	v.X = scalar
	v.Y = scalar
	v.Z = scalar
	return v
}

func (v *Vector3) SetComponent(index int, value float64) *Vector3 {
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

func (v *Vector3) Component(index int) float64 {
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

func (v *Vector3) Clone() *Vector3 {
	return NewVector3(v.X, v.Y, v.Z)
}

func (v *Vector3) Copy(vector Vector3) *Vector3 {
	v.X = vector.X
	v.Y = vector.Y
	v.Z = vector.Z
	return v
}

func (v *Vector3) Add(v2 Vector3) *Vector3 {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

func (v *Vector3) AddScalar(scalar float64) *Vector3 {
	v.X += scalar
	v.Y += scalar
	v.Z += scalar
	return v
}

func AddVector3s(a, b Vector3) *Vector3 {
	return NewVector3(a.X+b.X, a.Y+b.Y, a.Z+b.Z)
}

func (v *Vector3) AddScaledVector(vector Vector3, scalar float64) *Vector3 {
	v.X += vector.X * scalar
	v.Y += vector.Y * scalar
	v.Z += vector.Z * scalar
	return v
}

func (v *Vector3) Sub(v2 Vector3) *Vector3 {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}

func (v *Vector3) SubScalar(scalar float64) *Vector3 {
	v.X -= scalar
	v.Y -= scalar
	v.Z -= scalar
	return v
}

func SubVector3s(a, b Vector3) *Vector3 {
	return NewVector3(a.X-b.X, a.Y-b.Y, a.Z-b.Z)
}

func (v *Vector3) Multiply(v2 Vector3) *Vector3 {
	v.X *= v2.X
	v.Y *= v2.Y
	v.Z *= v2.Z
	return v
}

func (v *Vector3) MultiplyScalar(scalar float64) *Vector3 {
	v.X *= scalar
	v.Y *= scalar
	v.Z *= scalar
	return v
}

func MultiplyVector3s(a, b Vector3) *Vector3 {
	return NewVector3(a.X*b.X, a.Y*b.Y, a.Z*b.Z)
}

func (v *Vector3) ApplyEuler(e Euler) *Vector3 {
	return v.ApplyQuaternion(*NewIdentityQuaternion().SetFromEuler(e))
}

func (v *Vector3) ApplyAxisAngle(axis Vector3, angle float64) *Vector3 {
	return v.ApplyQuaternion(*NewIdentityQuaternion().SetFromAxisAngle(axis, angle))
}

func (v *Vector3) ApplyMatrix3(m Matrix3) *Vector3 {
	x, y, z := v.X, v.Y, v.Z
	v.X = m[0]*x + m[3]*y + m[6]*z
	v.Y = m[1]*x + m[4]*y + m[7]*z
	v.Z = m[2]*x + m[5]*y + m[8]*z
	return v
}

func (v *Vector3) ApplyNormalMatrix(m Matrix3) *Vector3 {
	return v.ApplyMatrix3(m).Normalize()
}

func (v *Vector3) ApplyMatrix4(m Matrix4) *Vector3 {
	x, y, z := v.X, v.Y, v.Z
	w := 1 / (m[3]*x + m[7]*y + m[11]*z + m[15])

	v.X = (m[0]*x + m[4]*y + m[8]*z + m[12]) * w
	v.Y = (m[1]*x + m[5]*y + m[9]*z + m[13]) * w
	v.Z = (m[2]*x + m[6]*y + m[10]*z + m[14]) * w
	return v
}

func (v *Vector3) ApplyQuaternion(q Quaternion) *Vector3 {
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

func (v *Vector3) TransformDirection(m Matrix4) *Vector3 {
	x, y, z := v.X, v.Y, v.Z
	v.X = m[0]*x + m[4]*y + m[8]*z
	v.Y = m[1]*x + m[5]*y + m[9]*z
	v.Z = m[2]*x + m[6]*y + m[10]*z
	return v.Normalize()
}

func (v *Vector3) Divide(v2 Vector3) *Vector3 {
	v.X /= v2.X
	v.Y /= v2.Y
	v.Z /= v2.Z
	return v
}

func (v *Vector3) DivideScalar(scalar float64) *Vector3 {
	return v.MultiplyScalar(1 / scalar)
}

func (v *Vector3) Min(v2 Vector3) *Vector3 {
	v.X = min(v.X, v2.X)
	v.Y = min(v.Y, v2.Y)
	v.Z = min(v.Z, v2.Z)
	return v
}

func (v *Vector3) Max(v2 Vector3) *Vector3 {
	v.X = max(v.X, v2.X)
	v.Y = max(v.Y, v2.Y)
	v.Z = max(v.Z, v2.Z)
	return v
}

func (v *Vector3) Clamp(minVal, maxVal Vector3) *Vector3 {
	v.X = max(minVal.X, min(maxVal.X, v.X))
	v.Y = max(minVal.Y, min(maxVal.Y, v.Y))
	v.Z = max(minVal.Z, min(maxVal.Z, v.Z))
	return v
}

func (v *Vector3) ClampScalar(minVal, maxVal float64) *Vector3 {
	v.X = max(minVal, min(maxVal, v.X))
	v.Y = max(minVal, min(maxVal, v.Y))
	v.Z = max(minVal, min(maxVal, v.Z))
	return v
}

func (v *Vector3) ClampLength(minVal, maxVal float64) *Vector3 {
	length := v.Length()
	return v.DivideScalar(length).MultiplyScalar(max(minVal, min(maxVal, length)))
}

func (v *Vector3) Floor() *Vector3 {
	v.X = math.Floor(v.X)
	v.Y = math.Floor(v.Y)
	v.Z = math.Floor(v.Z)
	return v
}

func (v *Vector3) Ceil() *Vector3 {
	v.X = math.Ceil(v.X)
	v.Y = math.Ceil(v.Y)
	v.Z = math.Ceil(v.Z)
	return v
}

func (v *Vector3) Round() *Vector3 {
	v.X = math.Round(v.X)
	v.Y = math.Round(v.Y)
	v.Z = math.Round(v.Z)
	return v
}

func (v *Vector3) RoundToZero() *Vector3 {
	v.X = math.Trunc(v.X)
	v.Y = math.Trunc(v.Y)
	v.Z = math.Trunc(v.Z)
	return v
}

func (v *Vector3) Negate() *Vector3 {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
	return v
}

func (v *Vector3) Dot(v2 Vector3) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v *Vector3) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v *Vector3) Length() float64 {
	return math.Sqrt(v.LengthSq())
}

func (v *Vector3) Normalize() *Vector3 {
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

func (v *Vector3) SetLength(length float64) *Vector3 {
	return v.Normalize().MultiplyScalar(length)
}

func (v *Vector3) Lerp(v2 Vector3, alpha float64) *Vector3 {
	v.X += (v2.X - v.X) * alpha
	v.Y += (v2.Y - v.Y) * alpha
	v.Z += (v2.Z - v.Z) * alpha
	return v
}

func LerpVector3s(v1, v2 Vector3, alpha float64) *Vector3 {
	return NewVector3(
		v1.X+(v2.X-v1.X)*alpha,
		v1.Y+(v2.Y-v1.Y)*alpha,
		v1.Z+(v2.Z-v1.Z)*alpha,
	)
}

func (v *Vector3) Cross(v2 Vector3) *Vector3 {
	ax, ay, az := v.X, v.Y, v.Z
	bx, by, bz := v2.X, v2.Y, v2.Z

	v.X = ay*bz - az*by
	v.Y = az*bx - ax*bz
	v.Z = ax*by - ay*bx

	return v
}

func CrossVector3s(a, b Vector3) *Vector3 {
	return a.Clone().Cross(b)
}

func (v *Vector3) ProjectOnVector(v2 Vector3) *Vector3 {
	denominator := v2.LengthSq()
	if denominator == 0 {
		return v.Set(0, 0, 0)
	}
	scalar := v.Dot(v2) / denominator
	return v.Copy(v2).MultiplyScalar(scalar)
}

func (v *Vector3) ProjectOnPlane(planeNormal Vector3) *Vector3 {
	return v.Sub(*v.Clone().ProjectOnVector(planeNormal))
}

func (v *Vector3) Reflect(normal Vector3) *Vector3 {
	return v.Sub(*v.Clone().ProjectOnPlane(normal).MultiplyScalar(2))
}

func (v *Vector3) AngleTo(v2 Vector3) float64 {
	denominator := math.Sqrt(v.LengthSq() * v2.LengthSq())
	if denominator == 0 {
		return math.Pi / 2
	}

	theta := v.Dot(v2) / denominator
	return math.Acos(Clamp(theta, -1, 1))
}

func (v *Vector3) DistanceTo(v2 Vector3) float64 {
	return math.Sqrt(v.DistanceToSquared(v2))
}

func (v *Vector3) DistanceToSquared(v2 Vector3) float64 {
	dx := v.X - v2.X
	dy := v.Y - v2.Y
	dz := v.Z - v2.Z
	return dx*dx + dy*dy + dz*dz
}

func (v *Vector3) ManhattanDistanceTo(v2 Vector3) float64 {
	return math.Abs(v.X-v2.X) + math.Abs(v.Y-v2.Y) + math.Abs(v.Z-v2.Z)
}

func (v *Vector3) SetFromSpherical(s Spherical) *Vector3 {
	return v.SetFromSphericalCoords(s.Radius, s.Phi, s.Theta)
}

func (v *Vector3) SetFromSphericalCoords(radius, phi, theta float64) *Vector3 {
	sinPhiRadius := math.Sin(phi) * radius
	v.X = sinPhiRadius * math.Sin(theta)
	v.Y = math.Cos(phi) * radius
	v.Z = sinPhiRadius * math.Cos(theta)
	return v
}

func (v *Vector3) SetFromCylindrical(c Cylindrical) *Vector3 {
	return v.SetFromCylindricalCoords(c.Radius, c.Theta, c.Y)
}

func (v *Vector3) SetFromCylindricalCoords(radius, theta, y float64) *Vector3 {
	v.X = radius * math.Sin(theta)
	v.Y = y
	v.Z = radius * math.Cos(theta)
	return v
}

func (v *Vector3) SetFromMatrixPosition(m Matrix4) *Vector3 {
	return v.Set(m[12], m[13], m[14])
}

func (v *Vector3) SetFromMatrixScale(m Matrix4) *Vector3 {
	_, _, scale := m.Decompose()
	return v.Copy(scale)
}

func (v *Vector3) SetFromMatrixColumn(m Matrix4, index int) *Vector3 {
	return v.FromArray(m[:], index*4)
}

func (v *Vector3) Equals(v2 Vector3) bool {
	return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
}

func (v *Vector3) FromArray(array []float64, offset int) *Vector3 {
	v.X = array[offset]
	v.Y = array[offset+1]
	v.Z = array[offset+2]
	return v
}

func (v *Vector3) ToArray(array []float64, offset int) []float64 {
	array[offset] = v.X
	array[offset+1] = v.Y
	array[offset+2] = v.Z
	return array
}

func (v *Vector3) Random() *Vector3 {
	v.X = rand.Float64()
	v.Y = rand.Float64()
	v.Z = rand.Float64()
	return v
}

func (v *Vector3) RandomDirection() *Vector3 {
	u := (rand.Float64() - 0.5) * 2
	t := rand.Float64() * math.Pi * 2
	f := math.Sqrt(1 - u*u)
	v.X = f * math.Cos(t)
	v.Y = f * math.Sin(t)
	v.Z = u
	return v
}
