package mathx

import (
	"math"
	"math/rand"

	"golang.org/x/exp/constraints"
)

type Vector2[T constraints.Float] struct {
	X T
	Y T
}

func NewVector2[T constraints.Float](x, y T) *Vector2[T] {
	return &Vector2[T]{X: x, Y: y}
}

func NewZeroVector2[T constraints.Float]() *Vector2[T] {
	return NewVector2[T](0, 0)
}

func NewOneVector2[T constraints.Float]() *Vector2[T] {
	return NewVector2[T](1, 1)
}

func (v *Vector2[T]) Width() T {
	return v.X
}

func (v *Vector2[T]) SetWidth(value T) *Vector2[T] {
	v.X = value
	return v
}

func (v *Vector2[T]) Height() T {
	return v.Y
}

func (v *Vector2[T]) SetHeight(value T) *Vector2[T] {
	v.Y = value
	return v
}

func (v *Vector2[T]) Set(x, y T) *Vector2[T] {
	v.X = x
	v.Y = y
	return v
}

func (v *Vector2[T]) SetScalar(scalar T) *Vector2[T] {
	v.X = scalar
	v.Y = scalar
	return v
}

func (v *Vector2[T]) SetX(x T) *Vector2[T] {
	v.X = x
	return v
}

func (v *Vector2[T]) SetY(y T) *Vector2[T] {
	v.Y = y
	return v
}

func (v *Vector2[T]) SetComponent(index int, value T) *Vector2[T] {
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

func (v *Vector2[T]) Component(index int) T {
	switch index {
	case 0:
		return v.X
	case 1:
		return v.Y
	default:
		panic("index is out of range")
	}
}

func (v *Vector2[T]) Clone() *Vector2[T] {
	return NewVector2(v.X, v.Y)
}

func (v *Vector2[T]) Copy(vector Vector2[T]) *Vector2[T] {
	v.X = vector.X
	v.Y = vector.Y
	return v
}

func (v *Vector2[T]) Add(vector Vector2[T]) *Vector2[T] {
	v.X += vector.X
	v.Y += vector.Y
	return v
}

func (v *Vector2[T]) AddScalar(scalar T) *Vector2[T] {
	v.X += scalar
	v.Y += scalar
	return v
}

func AddVector2s[T constraints.Float](a, b Vector2[T]) *Vector2[T] {
	return NewVector2(a.X+b.X, a.Y+b.Y)
}

func AddScaledVector2[T constraints.Float](vector Vector2[T], scalar T) *Vector2[T] {
	return NewVector2(vector.X*scalar, vector.Y*scalar)
}

func (v *Vector2[T]) Sub(vector Vector2[T]) *Vector2[T] {
	v.X -= vector.X
	v.Y -= vector.Y
	return v
}

func (v *Vector2[T]) SubScalar(scalar T) *Vector2[T] {
	v.X -= scalar
	v.Y -= scalar
	return v
}

func SubVector2s[T constraints.Float](a, b Vector2[T]) *Vector2[T] {
	return NewVector2(a.X-b.X, a.Y-b.Y)
}

func (v *Vector2[T]) Multiply(vector Vector2[T]) *Vector2[T] {
	v.X *= vector.X
	v.Y *= vector.Y
	return v
}

func (v *Vector2[T]) MultiplyScalar(scalar T) *Vector2[T] {
	v.X *= scalar
	v.Y *= scalar
	return v
}

func (v *Vector2[T]) Divide(vector Vector2[T]) *Vector2[T] {
	v.X /= vector.X
	v.Y /= vector.Y
	return v
}

func (v *Vector2[T]) DivideScalar(scalar T) *Vector2[T] {
	return v.MultiplyScalar(1 / scalar)
}

func (v *Vector2[T]) ApplyMatrix3(m Matrix3[T]) *Vector2[T] {
	x := v.X
	y := v.Y
	v.X = m[0]*x + m[3]*y + m[6]
	v.Y = m[1]*x + m[4]*y + m[7]
	return v
}

func (v *Vector2[T]) Min(vector Vector2[T]) *Vector2[T] {
	v.X = min(v.X, vector.X)
	v.Y = min(v.Y, vector.Y)
	return v
}

func (v *Vector2[T]) Max(vector Vector2[T]) *Vector2[T] {
	v.X = max(v.X, vector.X)
	v.Y = max(v.Y, vector.Y)
	return v
}

func (v *Vector2[T]) Clamp(minVal, maxVal Vector2[T]) *Vector2[T] {
	v.X = max(minVal.X, min(maxVal.X, v.X))
	v.Y = max(minVal.Y, min(maxVal.Y, v.Y))
	return v
}

func (v *Vector2[T]) ClampScalar(minVal, maxVal T) *Vector2[T] {
	v.X = max(minVal, min(maxVal, v.X))
	v.Y = max(minVal, min(maxVal, v.Y))
	return v
}

func (v *Vector2[T]) ClampLength(minVal, maxVal T) *Vector2[T] {
	length := v.Length()
	return v.DivideScalar(length).MultiplyScalar(max(minVal, min(maxVal, length)))
}

func (v *Vector2[T]) Floor() *Vector2[T] {
	v.X = T(math.Floor(float64(v.X)))
	v.Y = T(math.Floor(float64(v.Y)))
	return v
}

func (v *Vector2[T]) Ceil() *Vector2[T] {
	v.X = T(math.Ceil(float64(v.X)))
	v.Y = T(math.Ceil(float64(v.Y)))
	return v
}

func (v *Vector2[T]) Round() *Vector2[T] {
	v.X = T(math.Round(float64(v.X)))
	v.Y = T(math.Round(float64(v.Y)))
	return v
}

func (v *Vector2[T]) RoundToZero() *Vector2[T] {
	v.X = T(math.Trunc(float64(v.X)))
	v.Y = T(math.Trunc(float64(v.Y)))
	return v
}

func (v *Vector2[T]) Negate() *Vector2[T] {
	v.X = -v.X
	v.Y = -v.Y
	return v
}

func (v *Vector2[T]) Dot(vector Vector2[T]) T {
	return v.X*vector.X + v.Y*vector.Y
}

func (v *Vector2[T]) Cross(vector Vector2[T]) T {
	return v.X*vector.Y - v.Y*vector.X
}

func (v *Vector2[T]) LengthSq() T {
	return v.X*v.X + v.Y*v.Y
}

func (v *Vector2[T]) Length() T {
	return T(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v *Vector2[T]) ManhattanLength() T {
	return T(math.Abs(float64(v.X) + math.Abs(float64(v.Y))))
}

func (v *Vector2[T]) Normalize() *Vector2[T] {
	return v.DivideScalar(v.Length())
}

func (v *Vector2[T]) Angle() T {
	angle := math.Atan2(float64(-v.Y), float64(-v.X)) + math.Pi
	return T(angle)
}

func (v *Vector2[T]) AngleTo(vector Vector2[T]) T {
	denominator := math.Sqrt(float64(v.LengthSq() * vector.LengthSq()))
	if denominator == 0 {
		return math.Pi / 2
	}
	theta := float64(v.Dot(vector)) / denominator
	return T(math.Acos(Clamp(theta, -1, 1)))
}

func (v *Vector2[T]) DistanceTo(vector Vector2[T]) T {
	return T(math.Sqrt(float64(v.DistanceToSquared(vector))))
}

func (v *Vector2[T]) DistanceToSquared(vector Vector2[T]) T {
	dx := v.X - vector.X
	dy := v.Y - vector.Y
	return dx*dx + dy*dy
}

func (v *Vector2[T]) ManhattanDistanceTo(vector Vector2[T]) T {
	return T(math.Abs(float64(v.X-vector.X)) + math.Abs(float64(v.Y-vector.Y)))
}

func (v *Vector2[T]) SetLength(length T) *Vector2[T] {
	return v.Normalize().MultiplyScalar(length)
}

func (v *Vector2[T]) Lerp(vector Vector2[T], alpha T) *Vector2[T] {
	v.X += (vector.X - v.X) * alpha
	v.Y += (vector.Y - v.Y) * alpha
	return v
}

func LerpVectors[T constraints.Float](v1, v2 Vector2[T], alpha T) *Vector2[T] {
	v := NewVector2(
		v1.X+(v2.X-v1.X)*alpha,
		v1.Y+(v2.Y-v1.Y)*alpha,
	)
	return v
}

func (v *Vector2[T]) Equals(vector Vector2[T]) bool {
	return (vector.X == v.X) && (vector.Y == v.Y)
}

func (v *Vector2[T]) FromArray(array []T, offset int) *Vector2[T] {
	v.X = array[offset]
	v.Y = array[offset+1]
	return v
}

func (v *Vector2[T]) ToArray(array []T, offset int) []T {
	array[offset] = v.X
	array[offset+1] = v.Y
	return array
}

func (v *Vector2[T]) RotateAround(center Vector2[T], angle T) *Vector2[T] {
	c := T(math.Cos(float64(angle)))
	s := T(math.Sin(float64(angle)))
	x := v.X - center.X
	y := v.Y - center.Y
	v.X = x*c - y*s + center.X
	v.Y = x*s + y*c + center.Y
	return v
}

func (v *Vector2[T]) Random() *Vector2[T] {
	v.X = T(rand.Float64())
	v.Y = T(rand.Float64())
	return v
}
