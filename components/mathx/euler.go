package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type EulerOrder int

const (
	EULER_ORDER_XYZ EulerOrder = iota
	EULER_ORDER_YXZ
	EULER_ORDER_ZXY
	EULER_ORDER_ZYX
	EULER_ORDER_YZX
	EULER_ORDER_XZY
	EULER_ORDER_DEFAULT = EULER_ORDER_XYZ
)

type Euler[T constraints.Float] struct {
	X     T
	Y     T
	Z     T
	Order EulerOrder
}

func NewEuler[T constraints.Float](x, y, z T, order EulerOrder) *Euler[T] {
	return &Euler[T]{
		X:     x,
		Y:     y,
		Z:     z,
		Order: order,
	}
}

func (e *Euler[T]) Set(x, y, z T, order EulerOrder) *Euler[T] {
	e.X = x
	e.Y = y
	e.Z = z
	e.Order = order
	return e
}

func (e *Euler[T]) Clone() *Euler[T] {
	return NewEuler(e.X, e.Y, e.Z, e.Order)
}

func (e *Euler[T]) Copy(euler *Euler[T]) *Euler[T] {
	e.X = euler.X
	e.Y = euler.Y
	e.Z = euler.Z
	e.Order = euler.Order
	return e
}

func (e *Euler[T]) SetFromRotationMatrix(m Matrix4[T], order EulerOrder, update bool) *Euler[T] {
	// assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled)

	m11, m12, m13 := float64(m[0]), float64(m[4]), float64(m[8])
	m21, m22, m23 := float64(m[1]), float64(m[5]), float64(m[9])
	m31, m32, m33 := float64(m[2]), float64(m[6]), float64(m[10])

	switch order {
	case EULER_ORDER_XYZ:
		e.Y = T(math.Asin(Clamp(m13, -1, 1)))

		if math.Abs(m13) < 0.9999999 {
			e.X = T(math.Atan2(-m23, m33))
			e.Z = T(math.Atan2(-m12, m11))
		} else {
			e.X = T(math.Atan2(m32, m22))
			e.Z = 0
		}
	case EULER_ORDER_YXZ:
		e.X = T(math.Asin(-Clamp(m23, -1, 1)))

		if math.Abs(m23) < 0.9999999 {
			e.Y = T(math.Atan2(m13, m33))
			e.Z = T(math.Atan2(m21, m22))
		} else {
			e.Y = T(math.Atan2(-m31, m11))
			e.Z = 0
		}
	case EULER_ORDER_ZXY:
		e.X = T(math.Asin(Clamp(m32, -1, 1)))

		if math.Abs(m32) < 0.9999999 {
			e.Y = T(math.Atan2(-m31, m33))
			e.Z = T(math.Atan2(-m12, m22))
		} else {
			e.Y = 0
			e.Z = T(math.Atan2(m21, m11))
		}
	case EULER_ORDER_ZYX:
		e.Y = T(math.Asin(-Clamp(m31, -1, 1)))

		if math.Abs(m31) < 0.9999999 {
			e.X = T(math.Atan2(m32, m33))
			e.Z = T(math.Atan2(m21, m11))
		} else {
			e.X = 0
			e.Z = T(math.Atan2(-m12, m22))
		}
	case EULER_ORDER_YZX:
		e.Z = T(math.Asin(Clamp(m21, -1, 1)))

		if math.Abs(m21) < 0.9999999 {
			e.X = T(math.Atan2(-m23, m22))
			e.Y = T(math.Atan2(-m31, m11))
		} else {
			e.X = 0
			e.Y = T(math.Atan2(m13, m33))
		}
	case EULER_ORDER_XZY:
		e.Z = T(math.Asin(-Clamp(m12, -1, 1)))

		if math.Abs(m12) < 0.9999999 {
			e.X = T(math.Atan2(m32, m22))
			e.Y = T(math.Atan2(m13, m11))
		} else {
			e.X = T(math.Atan2(-m23, m33))
			e.Y = 0
		}
	default:
		panic("Invalid order")
	}

	e.Order = order

	return e
}

func (e *Euler[T]) SetFromQuaternion(q Quaternion[T], order EulerOrder, update bool) *Euler[T] {
	m := MakeRotationFromQuaternion(q)
	return e.SetFromRotationMatrix(*m, order, update)
}

func (e *Euler[T]) SetFromVector3(v Vector3[T], order EulerOrder) *Euler[T] {
	return e.Set(v.X, v.Y, v.Z, order)
}

func (e *Euler[T]) Reorder(newOrder EulerOrder) *Euler[T] {
	// WARNING: this discards revolution information -bhouston
	q := NewIdentityQuaternion[T]().SetFromEuler(*e)
	return e.SetFromQuaternion(*q, newOrder, false)
}

func (e *Euler[T]) Equals(euler *Euler[T]) bool {
	return euler.X == e.X && euler.Y == e.Y && euler.Z == e.Z && euler.Order == e.Order
}

func (e *Euler[T]) FromArray(array []T) *Euler[T] {
	e.X = array[0]
	e.Y = array[1]
	e.Z = array[2]
	if len(array) > 3 {
		e.Order = EulerOrder(array[3])
	}

	return e
}

func (e *Euler[T]) ToArray(array []T, offset int) []T {
	array[offset] = e.X
	array[offset+1] = e.Y
	array[offset+2] = e.Z
	array[offset+3] = T(e.Order)
	return array
}
