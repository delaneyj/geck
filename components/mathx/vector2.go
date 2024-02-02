package mathx

import (
	"math"
	"math/rand"
)

type Vector2 struct {
	X float64
	Y float64
}

func NewVector2(x, y float64) *Vector2 {
	return &Vector2{X: x, Y: y}
}

func NewZeroVector2() *Vector2 {
	return NewVector2(0, 0)
}

func NewOneVector2() *Vector2 {
	return NewVector2(1, 1)
}

func (v *Vector2) Width() float64 {
	return v.X
}

func (v *Vector2) SetWidth(value float64) *Vector2 {
	v.X = value
	return v
}

func (v *Vector2) Height() float64 {
	return v.Y
}

func (v *Vector2) SetHeight(value float64) *Vector2 {
	v.Y = value
	return v
}

func (v *Vector2) Set(x, y float64) *Vector2 {
	v.X = x
	v.Y = y
	return v
}

func (v *Vector2) SetScalar(scalar float64) *Vector2 {
	v.X = scalar
	v.Y = scalar
	return v
}

func (v *Vector2) SetX(x float64) *Vector2 {
	v.X = x
	return v
}

func (v *Vector2) SetY(y float64) *Vector2 {
	v.Y = y
	return v
}

func (v *Vector2) SetComponent(index int, value float64) *Vector2 {
	switch index {
	case 0:
		v.X = value
	case 1:
		v.Y = value
	default:
		panic("index is out of range")
	}
	return v
}

func (v *Vector2) Component(index int) float64 {
	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		panic("index is out of range")
	}
}

func (v *Vector2) Clone() *Vector2 {
	return NewVector2(v.X, v.Y)
}

func (v *Vector2) Copy(vector Vector2) *Vector2 {
	v.X = vector.X
	v.Y = vector.Y
	return v
}

func (v *Vector2) Add(vector Vector2) *Vector2 {
	v.X += vector.X
	v.Y += vector.Y
	return v
}

func (v *Vector2) AddScalar(scalar float64) *Vector2 {
	v.X += scalar
	v.Y += scalar
	return v
}

func AddVector2s(a, b Vector2) *Vector2 {
	return NewVector2(a.X+b.X, a.Y+b.Y)
}

func AddScaledVector2(vector Vector2, scalar float64) *Vector2 {
	return NewVector2(vector.X*scalar, vector.Y*scalar)
}

func (v *Vector2) Sub(vector Vector2) *Vector2 {
	v.X -= vector.X
	v.Y -= vector.Y
	return v
}

func (v *Vector2) SubScalar(scalar float64) *Vector2 {
	v.X -= scalar
	v.Y -= scalar
	return v
}

func SubVector2s(a, b Vector2) *Vector2 {
	return NewVector2(a.X-b.X, a.Y-b.Y)
}

func (v *Vector2) Multiply(vector Vector2) *Vector2 {
	v.X *= vector.X
	v.Y *= vector.Y
	return v
}

func (v *Vector2) MultiplyScalar(scalar float64) *Vector2 {
	v.X *= scalar
	v.Y *= scalar
	return v
}

func (v *Vector2) Divide(vector Vector2) *Vector2 {
	v.X /= vector.X
	v.Y /= vector.Y
	return v
}

func (v *Vector2) DivideScalar(scalar float64) *Vector2 {
	return v.MultiplyScalar(1 / scalar)
}

func (v *Vector2) ApplyMatrix3(m Matrix3) *Vector2 {
	x := v.X
	y := v.Y
	v.X = m[0]*x + m[3]*y + m[6]
	v.Y = m[1]*x + m[4]*y + m[7]
	return v
}

func (v *Vector2) Min(vector Vector2) *Vector2 {
	v.X = min(v.X, vector.X)
	v.Y = min(v.Y, vector.Y)
	return v
}

func (v *Vector2) Max(vector Vector2) *Vector2 {
	v.X = max(v.X, vector.X)
	v.Y = max(v.Y, vector.Y)
	return v
}

func (v *Vector2) Clamp(minVal, maxVal Vector2) *Vector2 {
	v.X = max(minVal.X, min(maxVal.X, v.X))
	v.Y = max(minVal.Y, min(maxVal.Y, v.Y))
	return v
}

func (v *Vector2) ClampScalar(minVal, maxVal float64) *Vector2 {
	v.X = max(minVal, min(maxVal, v.X))
	v.Y = max(minVal, min(maxVal, v.Y))
	return v
}

func (v *Vector2) ClampLength(minVal, maxVal float64) *Vector2 {
	length := v.Length()
	return v.DivideScalar(length).MultiplyScalar(max(minVal, min(maxVal, length)))
}

func (v *Vector2) Floor() *Vector2 {
	v.X = math.Floor(v.X)
	v.Y = math.Floor(v.Y)
	return v
}

func (v *Vector2) Ceil() *Vector2 {
	v.X = math.Ceil(v.X)
	v.Y = math.Ceil(v.Y)
	return v
}

func (v *Vector2) Round() *Vector2 {
	v.X = math.Round(v.X)
	v.Y = math.Round(v.Y)
	return v
}

func (v *Vector2) RoundToZero() *Vector2 {
	v.X = math.Trunc(v.X)
	v.Y = math.Trunc(v.Y)
	return v
}

func (v *Vector2) Negate() *Vector2 {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

func (v *Vector2) Dot(vector Vector2) float64 {
	return v.X*vector.X + v.Y*vector.Y
}

func (v *Vector2) Cross(vector Vector2) float64 {
	return v.X*vector.Y - v.Y*vector.X
}

func (v *Vector2) LengthSq() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v *Vector2) ManhattanLength() float64 {
	return math.Abs(v.X) + math.Abs(v.Y)
}

func (v *Vector2) Normalize() *Vector2 {
	return v.DivideScalar(v.Length())
}

func (v *Vector2) Angle() float64 {
	angle := math.Atan2(-v.Y, -v.X) + math.Pi
	return angle
}

func (v *Vector2) AngleTo(vector Vector2) float64 {
	denominator := math.Sqrt(v.LengthSq() * vector.LengthSq())
	if denominator == 0 {
		return math.Pi / 2
	}
	theta := v.Dot(vector) / denominator
	return math.Acos(Clamp(theta, -1, 1))
}

func (v *Vector2) DistanceTo(vector Vector2) float64 {
	return math.Sqrt(v.DistanceToSquared(vector))
}

func (v *Vector2) DistanceToSquared(vector Vector2) float64 {
	dx := v.X - vector.X
	dy := v.Y - vector.Y
	return dx*dx + dy*dy
}

func (v *Vector2) ManhattanDistanceTo(vector Vector2) float64 {
	return math.Abs(v.X-vector.X) + math.Abs(v.Y-vector.Y)
}

func (v *Vector2) SetLength(length float64) *Vector2 {
	return v.Normalize().MultiplyScalar(length)
}

func (v *Vector2) Lerp(vector Vector2, alpha float64) *Vector2 {
	v.X += (vector.X - v.X) * alpha
	v.Y += (vector.Y - v.Y) * alpha
	return v
}

func LerpVectors(v1, v2 Vector2, alpha float64) *Vector2 {
	v := NewVector2(
		v1.X+(v2.X-v1.X)*alpha,
		v1.Y+(v2.Y-v1.Y)*alpha,
	)
	return v
}

func (v *Vector2) Equals(vector Vector2) bool {
	return (vector.X == v.X) && (vector.Y == v.Y)
}

func (v *Vector2) FromArray(array []float64, offset int) *Vector2 {
	v.X = array[offset]
	v.Y = array[offset+1]
	return v
}

func (v *Vector2) ToArray(array []float64, offset int) []float64 {
	array[offset] = v.X
	array[offset+1] = v.Y
	return array
}

func (v *Vector2) RotateAround(center Vector2, angle float64) *Vector2 {
	c := math.Cos(angle)
	s := math.Sin(angle)
	x := v.X - center.X
	y := v.Y - center.Y
	v.X = x*c - y*s + center.X
	v.Y = x*s + y*c + center.Y
	return v
}

func (v *Vector2) Random() *Vector2 {
	v.X = rand.Float64()
	v.Y = rand.Float64()
	return v
}
