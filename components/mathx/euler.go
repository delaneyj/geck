package mathx

import "math"

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

type Euler struct {
	X     float64
	Y     float64
	Z     float64
	Order EulerOrder
}

func NewEuler(x, y, z float64, order EulerOrder) *Euler {
	return &Euler{
		X:     x,
		Y:     y,
		Z:     z,
		Order: order,
	}
}

func (e *Euler) Set(x, y, z float64, order EulerOrder) *Euler {
	e.X = x
	e.Y = y
	e.Z = z
	e.Order = order
	return e
}

func (e *Euler) Clone() *Euler {
	return NewEuler(e.X, e.Y, e.Z, e.Order)
}

func (e *Euler) Copy(euler *Euler) *Euler {
	e.X = euler.X
	e.Y = euler.Y
	e.Z = euler.Z
	e.Order = euler.Order
	return e
}

func (e *Euler) SetFromRotationMatrix(m Matrix4, order EulerOrder, update bool) *Euler {
	// assumes the upper 3x3 of m is a pure rotation matrix (i.e, unscaled)

	m11, m12, m13 := m[0], m[4], m[8]
	m21, m22, m23 := m[1], m[5], m[9]
	m31, m32, m33 := m[2], m[6], m[10]

	switch order {
	case EULER_ORDER_XYZ:
		e.Y = math.Asin(Clamp(m13, -1, 1))

		if math.Abs(m13) < 0.9999999 {
			e.X = math.Atan2(-m23, m33)
			e.Z = math.Atan2(-m12, m11)
		} else {
			e.X = math.Atan2(m32, m22)
			e.Z = 0
		}
	case EULER_ORDER_YXZ:
		e.X = math.Asin(-Clamp(m23, -1, 1))

		if math.Abs(m23) < 0.9999999 {
			e.Y = math.Atan2(m13, m33)
			e.Z = math.Atan2(m21, m22)
		} else {
			e.Y = math.Atan2(-m31, m11)
			e.Z = 0
		}
	case EULER_ORDER_ZXY:
		e.X = math.Asin(Clamp(m32, -1, 1))

		if math.Abs(m32) < 0.9999999 {
			e.Y = math.Atan2(-m31, m33)
			e.Z = math.Atan2(-m12, m22)
		} else {
			e.Y = 0
			e.Z = math.Atan2(m21, m11)
		}
	case EULER_ORDER_ZYX:
		e.Y = math.Asin(-Clamp(m31, -1, 1))

		if math.Abs(m31) < 0.9999999 {
			e.X = math.Atan2(m32, m33)
			e.Z = math.Atan2(m21, m11)
		} else {
			e.X = 0
			e.Z = math.Atan2(-m12, m22)
		}
	case EULER_ORDER_YZX:
		e.Z = math.Asin(Clamp(m21, -1, 1))

		if math.Abs(m21) < 0.9999999 {
			e.X = math.Atan2(-m23, m22)
			e.Y = math.Atan2(-m31, m11)
		} else {
			e.X = 0
			e.Y = math.Atan2(m13, m33)
		}
	case EULER_ORDER_XZY:
		e.Z = math.Asin(-Clamp(m12, -1, 1))

		if math.Abs(m12) < 0.9999999 {
			e.X = math.Atan2(m32, m22)
			e.Y = math.Atan2(m13, m11)
		} else {
			e.X = math.Atan2(-m23, m33)
			e.Y = 0
		}
	default:
		panic("Invalid order")
	}

	e.Order = order

	return e
}

func (e *Euler) SetFromQuaternion(q Quaternion, order EulerOrder, update bool) *Euler {
	m := MakeRotationFromQuaternion(q)
	return e.SetFromRotationMatrix(*m, order, update)
}

func (e *Euler) SetFromVector3(v Vector3, order EulerOrder) *Euler {
	return e.Set(v.X, v.Y, v.Z, order)
}

func (e *Euler) Reorder(newOrder EulerOrder) *Euler {
	// WARNING: this discards revolution information -bhouston
	q := NewIdentityQuaternion().SetFromEuler(*e)
	return e.SetFromQuaternion(*q, newOrder, false)
}

func (e *Euler) Equals(euler *Euler) bool {
	return euler.X == e.X && euler.Y == e.Y && euler.Z == e.Z && euler.Order == e.Order
}

func (e *Euler) FromArray(array []float64) *Euler {
	e.X = array[0]
	e.Y = array[1]
	e.Z = array[2]
	if len(array) > 3 {
		e.Order = EulerOrder(array[3])
	}

	return e
}

func (e *Euler) ToArray(array []float64, offset int) []float64 {
	array[offset] = e.X
	array[offset+1] = e.Y
	array[offset+2] = e.Z
	array[offset+3] = float64(e.Order)
	return array
}
