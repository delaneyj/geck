package mathx

import "math"

type Quaternion struct {
	X, Y, Z, W float64
}

func NewQuaternion(x, y, z, w float64) *Quaternion {
	return &Quaternion{X: x, Y: y, Z: z, W: w}
}

func NewIdentityQuaternion() *Quaternion {
	return &Quaternion{X: 0, Y: 0, Z: 0, W: 1}
}

func (q *Quaternion) SlerpFlat(dst []float64, dstOffset int, src0 []float64, srcOffset0 int, src1 []float64, srcOffset1 int, t float64) {
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
		if sqrSin > EPSILON64 {
			sin := math.Sqrt(sqrSin)
			len := math.Atan2(sin, cos*float64(dir))
			s = math.Sin(s*len) / sin
			t = math.Sin(t*len) / sin
		}
		tDir := t * float64(dir)
		x0 = x0*s + x1*tDir
		y0 = y0*s + y1*tDir
		z0 = z0*s + z1*tDir
		w0 = w0*s + w1*tDir
		if s == 1-t {
			f := 1 / math.Sqrt(x0*x0+y0*y0+z0*z0+w0*w0)
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

func (q *Quaternion) MultiplyQuaternionsFlat(dst []float64, dstOffset int, src0 []float64, srcOffset0 int, src1 []float64, srcOffset1 int) []float64 {

	x0, y0, z0, w0 := src0[srcOffset0], src0[srcOffset0+1], src0[srcOffset0+2], src0[srcOffset0+3]
	x1, y1, z1, w1 := src1[srcOffset1], src1[srcOffset1+1], src1[srcOffset1+2], src1[srcOffset1+3]

	dst[dstOffset] = x0*w1 + w0*x1 + y0*z1 - z0*y1
	dst[dstOffset+1] = y0*w1 + w0*y1 + z0*x1 - x0*z1
	dst[dstOffset+2] = z0*w1 + w0*z1 + x0*y1 - y0*x1
	dst[dstOffset+3] = w0*w1 - x0*x1 - y0*y1 - z0*z1

	return dst
}

func (q *Quaternion) Set(x, y, z, w float64) *Quaternion {
	q.X = x
	q.Y = y
	q.Z = z
	q.W = w
	return q
}

func (q *Quaternion) Clone() *Quaternion {
	return NewQuaternion(q.X, q.Y, q.Z, q.W)
}

func (q *Quaternion) Copy(quaternion Quaternion) *Quaternion {
	q.X = quaternion.X
	q.Y = quaternion.Y
	q.Z = quaternion.Z
	q.W = quaternion.W
	return q
}

func (q *Quaternion) SetFromEuler(euler Euler) *Quaternion {
	x, y, z, order := euler.X, euler.Y, euler.Z, euler.Order
	c1 := math.Cos(x / 2)
	c2 := math.Cos(y / 2)
	c3 := math.Cos(z / 2)
	s1 := math.Sin(x / 2)
	s2 := math.Sin(y / 2)
	s3 := math.Sin(z / 2)
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

func (q *Quaternion) SetFromAxisAngle(axis Vector3, angle float64) *Quaternion {
	halfAngle := angle / 2
	s := math.Sin(halfAngle)
	q.X = axis.X * s
	q.Y = axis.Y * s
	q.Z = axis.Z * s
	q.W = math.Cos(halfAngle)
	return q
}

func (q *Quaternion) SetFromRotationMatrix(m Matrix4) *Quaternion {
	m11, m12, m13 := m[0], m[4], m[8]
	m21, m22, m23 := m[1], m[5], m[9]
	m31, m32, m33 := m[2], m[6], m[10]
	trace := m11 + m22 + m33
	var s float64
	if trace > 0 {
		s = 0.5 / math.Sqrt(trace+1)
		q.W = 0.25 / s
		q.X = (m32 - m23) * s
		q.Y = (m13 - m31) * s
		q.Z = (m21 - m12) * s
	}
	if m11 > m22 && m11 > m33 {
		s = 2 * math.Sqrt(1+m11-m22-m33)
		q.W = (m32 - m23) / s
		q.X = 0.25 * s
		q.Y = (m12 + m21) / s
		q.Z = (m13 + m31) / s
	}

	if m22 > m33 {
		s = 2 * math.Sqrt(1+m22-m11-m33)
		q.W = (m13 - m31) / s
		q.X = (m12 + m21) / s
		q.Y = 0.25 * s
		q.Z = (m23 + m32) / s
	}

	s = 2 * math.Sqrt(1+m33-m11-m22)
	q.W = (m21 - m12) / s
	q.X = (m13 + m31) / s
	q.Y = (m23 + m32) / s
	q.Z = 0.25 * s
	return q
}

func (q *Quaternion) SetFromUnitVectors(vFrom, vTo Vector3) *Quaternion {
	r := vFrom.Dot(vTo) + 1
	if r < EPSILON64 {
		r = 0
		if math.Abs(vFrom.X) > math.Abs(vFrom.Z) {
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

func (q *Quaternion) AngleTo(qb Quaternion) float64 {
	return 2 * math.Acos(math.Abs(Clamp(q.Dot(qb), -1, 1)))
}

func (q *Quaternion) RotateTowards(qb Quaternion, step float64) *Quaternion {
	angle := q.AngleTo(qb)
	if angle == 0 {
		return q
	}

	t := math.Min(1, step/angle)
	return q.Slerp(qb, t)
}

func (q *Quaternion) Identity() *Quaternion {
	return q.Set(0, 0, 0, 1)
}

func (q *Quaternion) Invert() *Quaternion {
	return q.Conjugate()
}

func (q *Quaternion) Conjugate() *Quaternion {
	q.X *= -1
	q.Y *= -1
	q.Z *= -1
	return q
}

func (q *Quaternion) Dot(v Quaternion) float64 {
	return q.X*v.X + q.Y*v.Y + q.Z*v.Z + q.W*v.W
}

func (q *Quaternion) LengthSq() float64 {
	return q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W
}

func (q *Quaternion) Length() float64 {
	return math.Sqrt(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)
}

func (q *Quaternion) Normalize() *Quaternion {
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

func (q *Quaternion) Multiply(qb Quaternion) *Quaternion {
	qax, qay, qaz, qaw := q.X, q.Y, q.Z, q.W
	qbx, qby, qbz, qbw := qb.X, qb.Y, qb.Z, qb.W
	return NewQuaternion(
		qax*qbw+qaw*qbx+qay*qbz-qaz*qby,
		qay*qbw+qaw*qby+qaz*qbx-qax*qbz,
		qaz*qbw+qaw*qbz+qax*qby-qay*qbx,
		qaw*qbw-qax*qbx-qay*qby-qaz*qbz,
	)

}

func (q *Quaternion) Premultiply(qb Quaternion) *Quaternion {
	q2 := MultiplyQuaternions(qb, *q)
	q.Copy(*q2)
	return q
}

func MultiplyQuaternions(a, b Quaternion) *Quaternion {
	q := a.Clone()
	return q.Multiply(b)
}

func (q *Quaternion) Slerp(qb Quaternion, t float64) *Quaternion {
	if t == 0 {
		return q
	}
	if t == 1 {
		return q.Copy(qb)
	}
	x, y, z, w := q.X, q.Y, q.Z, q.W
	cosHalfTheta := w*qb.W + x*qb.X + y*qb.Y + z*qb.Z
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
	sqrSinHalfTheta := 1.0 - cosHalfTheta*cosHalfTheta
	if sqrSinHalfTheta <= EPSILON64 {
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
	ratioA := math.Sin((1-t)*halfTheta) / sinHalfTheta
	ratioB := math.Sin(t*halfTheta) / sinHalfTheta
	q.W = w*ratioA + q.W*ratioB
	q.X = x*ratioA + q.X*ratioB
	q.Y = y*ratioA + q.Y*ratioB
	q.Z = z*ratioA + q.Z*ratioB
	return q.Normalize()
}

func (q *Quaternion) Equals(quaternion Quaternion) bool {
	return quaternion.X == q.X && quaternion.Y == q.Y && quaternion.Z == q.Z && quaternion.W == q.W
}

func (q *Quaternion) FromArray(array []float64, offset int) *Quaternion {
	q.X = array[offset]
	q.Y = array[offset+1]
	q.Z = array[offset+2]
	q.W = array[offset+3]
	return q
}

func (q *Quaternion) ToArray(array []float64, offset int) []float64 {
	array[offset] = q.X
	array[offset+1] = q.Y
	array[offset+2] = q.Z
	array[offset+3] = q.W
	return array
}
