package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Quaternion[T constraints.Float] struct {
	X, Y, Z, W T
}

func NewQuaternion[T constraints.Float](x, y, z, w T) *Quaternion[T] {
	return &Quaternion[T]{X: x, Y: y, Z: z, W: w}
}

func NewIdentityQuaternion[T constraints.Float]() *Quaternion[T] {
	return &Quaternion[T]{X: 0, Y: 0, Z: 0, W: 1}
}

func (q *Quaternion[T]) SlerpFlat(dst []T, dstOffset int, src0 []T, srcOffset0 int, src1 []T, srcOffset1 int, t T) {
	// fuzz-free, array-based Quaternion SLERP operation
	x0, y0, z0, w0 := src0[srcOffset0], src0[srcOffset0+1], src0[srcOffset0+2], src0[srcOffset0+3]
	x1, y1, z1, w1 := src1[srcOffset1], src1[srcOffset1+1], src1[srcOffset1+2], src1[srcOffset1+3]
	if t == 0 {
		dst[dstOffset] = x0
		dst[dstOffset+1] = y0
		dst[dstOffset+2] = z0
		dst[dstOffset+3] = w0
		return
	}
	if t == 1 {
		dst[dstOffset] = x1
		dst[dstOffset+1] = y1
		dst[dstOffset+2] = z1
		dst[dstOffset+3] = w1
		return
	}
	if w0 != w1 || x0 != x1 || y0 != y1 || z0 != z1 {
		s := 1 - t
		cos := x0*x1 + y0*y1 + z0*z1 + w0*w1
		dir := 1
		if cos < 0 {
			dir = -1
		}
		sqrSin := 1 - cos*cos
		if sqrSin > T(EPSILON) {
			sin := math.Sqrt(float64(sqrSin))
			len := math.Atan2(sin, float64(cos*T(dir)))
			s = T(math.Sin(float64(s)*len) / sin)
			t = T(math.Sin(float64(t)*len) / sin)
		}
		tDir := t * T(dir)
		x0 = x0*s + x1*tDir
		y0 = y0*s + y1*tDir
		z0 = z0*s + z1*tDir
		w0 = w0*s + w1*tDir
		if s == 1-t {
			f := T(1 / math.Sqrt(float64(x0*x0+y0*y0+z0*z0+w0*w0)))
			x0 *= f
			y0 *= f
			z0 *= f
			w0 *= f
		}
	}
	dst[dstOffset] = x0
	dst[dstOffset+1] = y0
	dst[dstOffset+2] = z0
	dst[dstOffset+3] = w0
}

func (q *Quaternion[T]) MultiplyQuaternionsFlat(dst []T, dstOffset int, src0 []T, srcOffset0 int, src1 []T, srcOffset1 int) []T {

	x0, y0, z0, w0 := src0[srcOffset0], src0[srcOffset0+1], src0[srcOffset0+2], src0[srcOffset0+3]
	x1, y1, z1, w1 := src1[srcOffset1], src1[srcOffset1+1], src1[srcOffset1+2], src1[srcOffset1+3]

	dst[dstOffset] = x0*w1 + w0*x1 + y0*z1 - z0*y1
	dst[dstOffset+1] = y0*w1 + w0*y1 + z0*x1 - x0*z1
	dst[dstOffset+2] = z0*w1 + w0*z1 + x0*y1 - y0*x1
	dst[dstOffset+3] = w0*w1 - x0*x1 - y0*y1 - z0*z1

	return dst
}

func (q *Quaternion[T]) Set(x, y, z, w T) *Quaternion[T] {
	q.X = x
	q.Y = y
	q.Z = z
	q.W = w
	return q
}

func (q *Quaternion[T]) Clone() *Quaternion[T] {
	return NewQuaternion(q.X, q.Y, q.Z, q.W)
}

func (q *Quaternion[T]) Copy(quaternion Quaternion[T]) *Quaternion[T] {
	q.X = quaternion.X
	q.Y = quaternion.Y
	q.Z = quaternion.Z
	q.W = quaternion.W
	return q
}

func (q *Quaternion[T]) SetFromEuler(euler Euler[T]) *Quaternion[T] {
	x, y, z, order := euler.X, euler.Y, euler.Z, euler.Order
	c1 := T(math.Cos(float64(x / 2)))
	c2 := T(math.Cos(float64(y / 2)))
	c3 := T(math.Cos(float64(z / 2)))
	s1 := T(math.Sin(float64(x / 2)))
	s2 := T(math.Sin(float64(y / 2)))
	s3 := T(math.Sin(float64(z / 2)))
	switch order {
	case EULER_ORDER_XYZ:
		q.Set(s1*c2*c3+c1*s2*s3, c1*s2*c3-s1*c2*s3, c1*c2*s3+s1*s2*c3, c1*c2*c3-s1*s2*s3)
	case EULER_ORDER_YXZ:
		q.Set(s1*c2*c3+c1*s2*s3, c1*s2*c3-s1*c2*s3, c1*c2*s3-s1*s2*c3, c1*c2*c3+s1*s2*s3)
	case EULER_ORDER_ZXY:
		q.Set(s1*c2*c3-c1*s2*s3, c1*s2*c3+s1*c2*s3, c1*c2*s3+s1*s2*c3, c1*c2*c3-s1*s2*s3)
	case EULER_ORDER_ZYX:
		q.Set(s1*c2*c3-c1*s2*s3, c1*s2*c3+s1*c2*s3, c1*c2*s3-s1*s2*c3, c1*c2*c3+s1*s2*s3)
	case EULER_ORDER_YZX:
		q.Set(s1*c2*c3+c1*s2*s3, c1*s2*c3-s1*c2*s3, c1*c2*s3-s1*s2*c3, c1*c2*c3+s1*s2*s3)
	case EULER_ORDER_XZY:
		q.Set(s1*c2*c3-c1*s2*s3, c1*s2*c3+s1*c2*s3, c1*c2*s3+s1*s2*c3, c1*c2*c3-s1*s2*s3)
	default:
		panic("Invalid order")
	}
	return q
}

func (q *Quaternion[T]) SetFromAxisAngle(axis Vector3[T], angle T) *Quaternion[T] {
	halfAngle := float64(angle / 2)
	s := T(math.Sin(halfAngle))
	q.X = T(axis.X * s)
	q.Y = T(axis.Y * s)
	q.Z = T(axis.Z * s)
	q.W = T(math.Cos(halfAngle))
	return q
}

func (q *Quaternion[T]) SetFromRotationMatrix(m Matrix4[T]) *Quaternion[T] {
	m11, m12, m13 := float64(m[0]), float64(m[4]), float64(m[8])
	m21, m22, m23 := float64(m[1]), float64(m[5]), float64(m[9])
	m31, m32, m33 := float64(m[2]), float64(m[6]), float64(m[10])
	trace := m11 + m22 + m33
	var s T
	if trace > 0 {
		s = T(0.5 / math.Sqrt(trace+1))
		q.W = 0.25 / s
		q.X = T((m32 - m23)) * s
		q.Y = T((m13 - m31)) * s
		q.Z = T((m21 - m12)) * s
	}
	if m11 > m22 && m11 > m33 {
		s = T(2 * math.Sqrt(1+m11-m22-m33))
		q.W = T(m32-m23) / s
		q.X = 0.25 * s
		q.Y = T(m12+m21) / s
		q.Z = T(m13+m31) / s
	}

	if m22 > m33 {
		s = T(2 * math.Sqrt(1+m22-m11-m33))
		q.W = T(m13-m31) / s
		q.X = T(m12+m21) / s
		q.Y = 0.25 * s
		q.Z = T(m23+m32) / s
	}

	s = T(2 * math.Sqrt(1+m33-m11-m22))
	q.W = T(m21-m12) / s
	q.X = T(m13+m31) / s
	q.Y = T(m23+m32) / s
	q.Z = 0.25 * s
	return q
}

func (q *Quaternion[T]) SetFromUnitVectors(vFrom, vTo Vector3[T]) *Quaternion[T] {
	r := vFrom.Dot(vTo) + 1
	if r < T(EPSILON) {
		r = 0
		if math.Abs(float64(vFrom.X)) > math.Abs(float64(vFrom.Z)) {
			q.X = -vFrom.Y
			q.Y = vFrom.X
			q.Z = 0
			q.W = r
		} else {
			q.X = 0
			q.Y = -vFrom.Z
			q.Z = vFrom.Y
			q.W = r
		}
	} else {
		q.X = vFrom.Y*vTo.Z - vFrom.Z*vTo.Y
		q.Y = vFrom.Z*vTo.X - vFrom.X*vTo.Z
		q.Z = vFrom.X*vTo.Y - vFrom.Y*vTo.X
		q.W = r
	}

	return q.Normalize()
}

func (q *Quaternion[T]) AngleTo(qb Quaternion[T]) T {
	return T(2 * math.Acos(math.Abs(float64(Clamp(q.Dot(qb), -1, 1)))))
}

func (q *Quaternion[T]) RotateTowards(qb Quaternion[T], step T) *Quaternion[T] {
	angle := q.AngleTo(qb)
	if angle == 0 {
		return q
	}

	t := min(1, step/angle)
	return q.Slerp(qb, t)
}

func (q *Quaternion[T]) Identity() *Quaternion[T] {
	return q.Set(0, 0, 0, 1)
}

func (q *Quaternion[T]) Invert() *Quaternion[T] {
	return q.Conjugate()
}

func (q *Quaternion[T]) Conjugate() *Quaternion[T] {
	q.X *= -1
	q.Y *= -1
	q.Z *= -1
	return q
}

func (q *Quaternion[T]) Dot(v Quaternion[T]) T {
	return q.X*v.X + q.Y*v.Y + q.Z*v.Z + q.W*v.W
}

func (q *Quaternion[T]) LengthSq() T {
	return q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
}

func (q *Quaternion[T]) Length() T {
	return T(math.Sqrt(float64(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)))
}

func (q *Quaternion[T]) Normalize() *Quaternion[T] {
	l := q.Length()
	if l == 0 {
		q.X = 0
		q.Y = 0
		q.Z = 0
		q.W = 1
	} else {
		l = 1 / l
		q.X *= l
		q.Y *= l
		q.Z *= l
		q.W *= l
	}
	return q
}

func (q *Quaternion[T]) Multiply(qb Quaternion[T]) *Quaternion[T] {
	qax, qay, qaz, qaw := q.X, q.Y, q.Z, q.W
	qbx, qby, qbz, qbw := qb.X, qb.Y, qb.Z, qb.W
	return NewQuaternion(
		qax*qbw+qaw*qbx+qay*qbz-qaz*qby,
		qay*qbw+qaw*qby+qaz*qbx-qax*qbz,
		qaz*qbw+qaw*qbz+qax*qby-qay*qbx,
		qaw*qbw-qax*qbx-qay*qby-qaz*qbz,
	)

}

func (q *Quaternion[T]) Premultiply(qb Quaternion[T]) *Quaternion[T] {
	q2 := MultiplyQuaternions(qb, *q)
	q.Copy(*q2)
	return q
}

func MultiplyQuaternions[T constraints.Float](a, b Quaternion[T]) *Quaternion[T] {
	q := a.Clone()
	return q.Multiply(b)
}

func (q *Quaternion[T]) Slerp(qb Quaternion[T], t T) *Quaternion[T] {
	if t == 0 {
		return q
	}
	if t == 1 {
		return q.Copy(qb)
	}
	x, y, z, w := q.X, q.Y, q.Z, q.W
	cosHalfTheta := float64(w*qb.W + x*qb.X + y*qb.Y + z*qb.Z)
	if cosHalfTheta < 0 {
		q.W = -qb.W
		q.X = -qb.X
		q.Y = -qb.Y
		q.Z = -qb.Z
		cosHalfTheta = -cosHalfTheta
	} else {
		q.Copy(qb)
	}
	if cosHalfTheta >= 1.0 {
		q.W = w
		q.X = x
		q.Y = y
		q.Z = z
		return q
	}
	sqrSinHalfTheta := float64(1.0 - cosHalfTheta*cosHalfTheta)
	if sqrSinHalfTheta <= EPSILON {
		s := 1 - t
		q.W = s*w + t*q.W
		q.X = s*x + t*q.X
		q.Y = s*y + t*q.Y
		q.Z = s*z + t*q.Z
		q.Normalize()
		return q
	}
	sinHalfTheta := math.Sqrt(sqrSinHalfTheta)
	halfTheta := math.Atan2(sinHalfTheta, cosHalfTheta)
	ratioA := T(math.Sin(float64(1-t)*halfTheta) / sinHalfTheta)
	ratioB := T(math.Sin(float64(t)*halfTheta) / sinHalfTheta)
	q.W = w*ratioA + q.W*ratioB
	q.X = x*ratioA + q.X*ratioB
	q.Y = y*ratioA + q.Y*ratioB
	q.Z = z*ratioA + q.Z*ratioB
	return q.Normalize()
}

func (q *Quaternion[T]) Equals(quaternion Quaternion[T]) bool {
	return quaternion.X == q.X && quaternion.Y == q.Y && quaternion.Z == q.Z && quaternion.W == q.W
}

func (q *Quaternion[T]) FromArray(array []T, offset int) *Quaternion[T] {
	q.X = array[offset]
	q.Y = array[offset+1]
	q.Z = array[offset+2]
	q.W = array[offset+3]
	return q
}

func (q *Quaternion[T]) ToArray(array []T, offset int) []T {
	array[offset] = q.X
	array[offset+1] = q.Y
	array[offset+2] = q.Z
	array[offset+3] = q.W
	return array
}
