package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Matrix4[T constraints.Float] [16]T

func NewMatrix4[T constraints.Float](n11, n12, n13, n14, n21, n22, n23, n24, n31, n32, n33, n34, n41, n42, n43, n44 T) *Matrix4[T] {
	m := &Matrix4[T]{}
	m.Set(n11, n12, n13, n14, n21, n22, n23, n24, n31, n32, n33, n34, n41, n42, n43, n44)
	return m
}

func NewMatrix4Identity[T constraints.Float]() *Matrix4[T] {
	m := &Matrix4[T]{}
	m.Identity()
	return m
}

func (m *Matrix4[T]) Set(n11, n12, n13, n14, n21, n22, n23, n24, n31, n32, n33, n34, n41, n42, n43, n44 T) *Matrix4[T] {
	m[0], m[4], m[8], m[12] = n11, n12, n13, n14
	m[1], m[5], m[9], m[13] = n21, n22, n23, n24
	m[2], m[6], m[10], m[14] = n31, n32, n33, n34
	m[3], m[7], m[11], m[15] = n41, n42, n43, n44
	return m
}

func (m *Matrix4[T]) Identity() *Matrix4[T] {
	return m.Set(
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
}

func (m *Matrix4[T]) Clone() *Matrix4[T] {
	m2 := &Matrix4[T]{}
	m2.Copy(*m)
	return m2
}

func (m *Matrix4[T]) Copy(matrix Matrix4[T]) *Matrix4[T] {
	*m = matrix
	return m
}

func (m *Matrix4[T]) CopyPosition(matrix Matrix4[T]) *Matrix4[T] {
	m[12], m[13], m[14] = matrix[12], matrix[13], matrix[14]
	return m
}

func (m *Matrix4[T]) SetFromMatrix3(matrix Matrix3[T]) *Matrix4[T] {
	m[0], m[1], m[2] = matrix[0], matrix[1], matrix[2]
	m[4], m[5], m[6] = matrix[3], matrix[4], matrix[5]
	m[8], m[9], m[10] = matrix[6], matrix[7], matrix[8]
	return m
}

func (m *Matrix4[T]) ExtractBasis() (xAxis, yAxis, zAxis Vector3[T]) {
	xAxis.SetFromMatrixColumn(*m, 0)
	yAxis.SetFromMatrixColumn(*m, 1)
	zAxis.SetFromMatrixColumn(*m, 2)
	return xAxis, yAxis, zAxis
}

func (m *Matrix4[T]) MakeBasis(xAxis, yAxis, zAxis Vector3[T]) *Matrix4[T] {
	m[0], m[4], m[8] = xAxis.X, yAxis.X, zAxis.X
	m[1], m[5], m[9] = xAxis.Y, yAxis.Y, zAxis.Y
	m[2], m[6], m[10] = xAxis.Z, yAxis.Z, zAxis.Z
	return m
}

func (m *Matrix4[T]) ExtractRotation(matrix Matrix4[T]) *Matrix4[T] {
	v1 := Vector3[T]{}
	scaleX := 1 / v1.SetFromMatrixColumn(matrix, 0).Length()
	scaleY := 1 / v1.SetFromMatrixColumn(matrix, 1).Length()
	scaleZ := 1 / v1.SetFromMatrixColumn(matrix, 2).Length()
	m[0] = matrix[0] * scaleX
	m[1] = matrix[1] * scaleX
	m[2] = matrix[2] * scaleX
	m[4] = matrix[4] * scaleY
	m[5] = matrix[5] * scaleY
	m[6] = matrix[6] * scaleY
	m[8] = matrix[8] * scaleZ
	m[9] = matrix[9] * scaleZ
	m[10] = matrix[10] * scaleZ
	m[3] = 0
	m[7] = 0
	m[11] = 0
	m[12] = 0
	m[13] = 0
	m[14] = 0
	m[15] = 1
	return m
}

func (m *Matrix4[T]) MakeRotationFromEuler(euler Euler[T]) *Matrix4[T] {
	x, y, z := float64(euler.X), float64(euler.Y), float64(euler.Z)
	a, b := math.Cos(x), math.Sin(x)
	c, d := math.Cos(y), math.Sin(y)
	e, f := math.Cos(z), math.Sin(z)

	if euler.Order == EULER_ORDER_XYZ {
		ae, af, be, bf := a*e, a*f, b*e, b*f
		m[0] = T(c * e)
		m[4] = T(-c * f)
		m[8] = T(d)
		m[1] = T(af + be*d)
		m[5] = T(ae - bf*d)
		m[9] = T(-b * c)
		m[2] = T(bf - ae*d)
		m[6] = T(be + af*d)
		m[10] = T(a * c)
	} else if euler.Order == EULER_ORDER_YXZ {
		ce, cf, de, df := c*e, c*f, d*e, d*f
		m[0] = T(ce + df*b)
		m[4] = T(de*b - cf)
		m[8] = T(a * d)
		m[1] = T(a * f)
		m[5] = T(a * e)
		m[9] = T(-b)
		m[2] = T(cf*b - de)
		m[6] = T(df + ce*b)
		m[10] = T(a * c)
	} else if euler.Order == EULER_ORDER_ZXY {
		ce, cf, de, df := c*e, c*f, d*e, d*f
		m[0] = T(ce - df*b)
		m[4] = T(-a * f)
		m[8] = T(de + cf*b)
		m[1] = T(cf + de*b)
		m[5] = T(a * e)
		m[9] = T(df - ce*b)
		m[2] = T(-a * d)
		m[6] = T(b)
		m[10] = T(a * c)
	} else if euler.Order == EULER_ORDER_ZYX {
		ae, af, be, bf := a*e, a*f, b*e, b*f
		m[0] = T(c * e)
		m[4] = T(be*d - af)
		m[8] = T(ae*d + bf)
		m[1] = T(c * f)
		m[5] = T(bf*d + ae)
		m[9] = T(af*d - be)
		m[2] = T(-d)
		m[6] = T(b * c)
		m[10] = T(a * c)
	} else if euler.Order == EULER_ORDER_YZX {
		ac, ad, bc, bd := a*c, a*d, b*c, b*d
		m[0] = T(c * e)
		m[4] = T(bd - ac*f)
		m[8] = T(bc*f + ad)
		m[1] = T(f)
		m[5] = T(a * e)
		m[9] = T(-b * e)
		m[2] = T(-d * e)
		m[6] = T(ad*f + bc)
		m[10] = T(ac - bd*f)
	} else if euler.Order == EULER_ORDER_XZY {
		ac, ad, bc, bd := a*c, a*d, b*c, b*d
		m[0] = T(c * e)
		m[4] = T(-f)
		m[8] = T(d * e)
		m[1] = T(ac*f + bd)
		m[5] = T(a * e)
		m[9] = T(ad*f - bc)
		m[2] = T(bc*f - ad)
		m[6] = T(b * e)
		m[10] = T(bd*f + ac)
	}

	// bottom row
	m[3] = 0
	m[7] = 0
	m[11] = 0

	// last column
	m[12] = 0
	m[13] = 0
	m[14] = 0
	m[15] = 1

	return m
}

func MakeRotationFromQuaternion[T constraints.Float](q Quaternion[T]) *Matrix4[T] {
	m := Matrix4[T]{}
	return m.Compose(Vector3[T]{}, q, *NewOneVector3[T]())
}

func (m *Matrix4[T]) LookAt(eye, target, up Vector3[T]) *Matrix4[T] {
	z := SubVector3s(eye, target)
	if z.LengthSq() == 0 {
		// eye and target are in the same position
		z.Z = 1
	}
	z.Normalize()

	x := CrossVector3s(up, *z)
	if x.LengthSq() == 0 {
		// up and z are parallel
		if math.Abs(float64(up.Z)) == 1 {
			z.X += EPSILON
		} else {
			z.Z += EPSILON
		}
		z.Normalize()
		x = CrossVector3s(up, *z)
	}
	x.Normalize()

	y := CrossVector3s(*z, *x)

	m[0], m[4], m[8] = x.X, y.X, z.X
	m[1], m[5], m[9] = x.Y, y.Y, z.Y
	m[2], m[6], m[10] = x.Z, y.Z, z.Z

	return m
}

func (m *Matrix4[T]) Multiply(m2 Matrix4[T]) *Matrix4[T] {
	a11, a12, a13, a14 := m[0], m[4], m[8], m[12]
	a21, a22, a23, a24 := m[1], m[5], m[9], m[13]
	a31, a32, a33, a34 := m[2], m[6], m[10], m[14]
	a41, a42, a43, a44 := m[3], m[7], m[11], m[15]

	b11, b12, b13, b14 := m2[0], m2[4], m2[8], m2[12]
	b21, b22, b23, b24 := m2[1], m2[5], m2[9], m2[13]
	b31, b32, b33, b34 := m2[2], m2[6], m2[10], m2[14]
	b41, b42, b43, b44 := m2[3], m2[7], m2[11], m2[15]

	m[0] = a11*b11 + a12*b21 + a13*b31 + a14*b41
	m[4] = a11*b12 + a12*b22 + a13*b32 + a14*b42
	m[8] = a11*b13 + a12*b23 + a13*b33 + a14*b43
	m[12] = a11*b14 + a12*b24 + a13*b34 + a14*b44

	m[1] = a21*b11 + a22*b21 + a23*b31 + a24*b41
	m[5] = a21*b12 + a22*b22 + a23*b32 + a24*b42
	m[9] = a21*b13 + a22*b23 + a23*b33 + a24*b43
	m[13] = a21*b14 + a22*b24 + a23*b34 + a24*b44

	m[2] = a31*b11 + a32*b21 + a33*b31 + a34*b41
	m[6] = a31*b12 + a32*b22 + a33*b32 + a34*b42
	m[10] = a31*b13 + a32*b23 + a33*b33 + a34*b43
	m[14] = a31*b14 + a32*b24 + a33*b34 + a34*b44

	m[3] = a41*b11 + a42*b21 + a43*b31 + a44*b41
	m[7] = a41*b12 + a42*b22 + a43*b32 + a44*b42
	m[11] = a41*b13 + a42*b23 + a43*b33 + a44*b43
	m[15] = a41*b14 + a42*b24 + a43*b34 + a44*b44

	return m
}

func (m *Matrix4[T]) Premultiply(m2 Matrix4[T]) *Matrix4[T] {
	tmp := m2.Clone().Multiply(*m)
	m.Copy(*tmp)
	return m
}

func MultiplyMatrice4s[T constraints.Float](a, b Matrix4[T]) *Matrix4[T] {
	return a.Clone().Multiply(b)
}

func (m *Matrix4[T]) MultiplyScalar(s T) *Matrix4[T] {
	for i := 0; i < 16; i++ {
		m[i] *= s
	}
	return m
}

func (m *Matrix4[T]) Demrminant() T {
	n11, n12, n13, n14 := m[0], m[4], m[8], m[12]
	n21, n22, n23, n24 := m[1], m[5], m[9], m[13]
	n31, n32, n33, n34 := m[2], m[6], m[10], m[14]
	n41, n42, n43, n44 := m[3], m[7], m[11], m[15]

	//TODO: make this more efficient
	//( based on http://www.euclideanspace.com/maths/algebra/matrix/functions/inverse/fourD/index.htm )

	return (n41*(n14*n23*n32-
		n13*n24*n32-
		n14*n22*n33+
		n12*n24*n33+
		n13*n22*n34-
		n12*n23*n34) +
		n42*(n11*n23*n34-
			n11*n24*n33+
			n14*n21*n33-
			n13*n21*n34+
			n13*n24*n31-
			n14*n23*n31) +
		n43*(n11*n24*n32-
			n11*n22*n34-
			n14*n21*n32+
			n12*n21*n34+
			n14*n22*n31-
			n12*n24*n31) +
		n44*(-n13*n22*n31-
			n11*n23*n32+
			n11*n22*n33+
			n13*n21*n32-
			n12*n21*n33+
			n12*n23*n31))

}

func (m *Matrix4[T]) Transpose() *Matrix4[T] {
	var tmp T
	tmp = m[1]
	m[1] = m[4]
	m[4] = tmp
	tmp = m[2]
	m[2] = m[8]
	m[8] = tmp
	tmp = m[6]
	m[6] = m[9]
	m[9] = tmp
	tmp = m[3]
	m[3] = m[12]
	m[12] = tmp
	tmp = m[7]
	m[7] = m[13]
	m[13] = tmp
	tmp = m[11]
	m[11] = m[14]
	m[14] = tmp
	return m
}

func (m *Matrix4[T]) SetPosition(v Vector3[T]) *Matrix4[T] {
	m[12], m[13], m[14] = v.X, v.Y, v.Z
	return m
}

func (m *Matrix4[T]) Invert() *Matrix4[T] {
	// based on http://www.euclideanspace.com/maths/algebra/matrix/functions/inverse/fourD/index.htm
	n11, n21, n31, n41 := m[0], m[1], m[2], m[3]
	n12, n22, n32, n42 := m[4], m[5], m[6], m[7]
	n13, n23, n33, n43 := m[8], m[9], m[10], m[11]
	n14, n24, n34, n44 := m[12], m[13], m[14], m[15]

	t11 := n23*n34*n42 - n24*n33*n42 + n24*n32*n43 - n22*n34*n43 - n23*n32*n44 + n22*n33*n44
	t12 := n14*n33*n42 - n13*n34*n42 - n14*n32*n43 + n12*n34*n43 + n13*n32*n44 - n12*n33*n44
	t13 := n13*n24*n42 - n14*n23*n42 + n14*n22*n43 - n12*n24*n43 - n13*n22*n44 + n12*n23*n44
	t14 := n14*n23*n32 - n13*n24*n32 - n14*n22*n33 + n12*n24*n33 + n13*n22*n34 - n12*n23*n34

	det := n11*t11 + n21*t12 + n31*t13 + n41*t14
	if det == 0 {
		return m.Set(
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
			0, 0, 0, 0,
		)
	}

	detInv := 1 / det

	m[0] = t11 * detInv
	m[1] = (n24*n33*n41 - n23*n34*n41 - n24*n31*n43 + n21*n34*n43 + n23*n31*n44 - n21*n33*n44) * detInv
	m[2] = (n22*n34*n41 - n24*n32*n41 + n24*n31*n42 - n21*n34*n42 - n22*n31*n44 + n21*n32*n44) * detInv
	m[3] = (n23*n32*n41 - n22*n33*n41 - n23*n31*n42 + n21*n33*n42 + n22*n31*n43 - n21*n32*n43) * detInv

	m[4] = t12 * detInv
	m[5] = (n13*n34*n41 - n14*n33*n41 + n14*n31*n43 - n11*n34*n43 - n13*n31*n44 + n11*n33*n44) * detInv
	m[6] = (n14*n32*n41 - n12*n34*n41 - n14*n31*n42 + n11*n34*n42 + n12*n31*n44 - n11*n32*n44) * detInv
	m[7] = (n12*n33*n41 - n13*n32*n41 + n13*n31*n42 - n11*n33*n42 - n12*n31*n43 + n11*n32*n43) * detInv

	m[8] = t13 * detInv
	m[9] = (n14*n23*n41 - n13*n24*n41 - n14*n21*n43 + n11*n24*n43 + n13*n21*n44 - n11*n23*n44) * detInv
	m[10] = (n12*n24*n41 - n14*n22*n41 + n14*n21*n42 - n11*n24*n42 - n12*n21*n44 + n11*n22*n44) * detInv
	m[11] = (n13*n22*n41 - n12*n23*n41 - n13*n21*n42 + n11*n23*n42 + n12*n21*n43 - n11*n22*n43) * detInv

	m[12] = t14 * detInv
	m[13] = (n13*n22*n31 - n12*n23*n31 - n13*n21*n32 + n11*n23*n32 + n12*n21*n33 - n11*n22*n33) * detInv
	m[14] = (n12*n23*n34 - n13*n22*n34 + n13*n21*n32 - n11*n23*n32 - n12*n21*n33 + n11*n22*n33) * detInv
	m[15] = (n13*n22*n31 - n12*n23*n31 - n13*n21*n32 + n11*n23*n32 + n12*n21*n33 - n11*n22*n33) * detInv

	return m
}

func (m *Matrix4[T]) Scale(v Vector3[T]) *Matrix4[T] {
	x, y, z := v.X, v.Y, v.Z
	m[0], m[4], m[8] = m[0]*x, m[4]*y, m[8]*z
	m[1], m[5], m[9] = m[1]*x, m[5]*y, m[9]*z
	m[2], m[6], m[10] = m[2]*x, m[6]*y, m[10]*z
	m[3], m[7], m[11] = m[3]*x, m[7]*y, m[11]*z
	return m
}

func (m *Matrix4[T]) MaxScaleOnAxis() T {
	scaleXSq := m[0]*m[0] + m[1]*m[1] + m[2]*m[2]
	scaleYSq := m[4]*m[4] + m[5]*m[5] + m[6]*m[6]
	scaleZSq := m[8]*m[8] + m[9]*m[9] + m[10]*m[10]
	return T(math.Sqrt(float64(max(scaleXSq, max(scaleYSq, scaleZSq)))))
}

func (m *Matrix4[T]) MakeTranslation(v Vector3[T]) *Matrix4[T] {
	m.Set(
		1, 0, 0, v.X,
		0, 1, 0, v.Y,
		0, 0, 1, v.Z,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) MakeRotationX(theta T) *Matrix4[T] {
	c := T(math.Cos(float64(theta)))
	s := T(math.Sin(float64(theta)))
	m.Set(
		1, 0, 0, 0,
		0, c, -s, 0,
		0, s, c, 0,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) MakeRotationY(theta T) *Matrix4[T] {
	c := T(math.Cos(float64(theta)))
	s := T(math.Sin(float64(theta)))
	m.Set(
		c, 0, s, 0,
		0, 1, 0, 0,
		-s, 0, c, 0,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) MakeRotationZ(theta T) *Matrix4[T] {
	c := T(math.Cos(float64(theta)))
	s := T(math.Sin(float64(theta)))
	m.Set(
		c, -s, 0, 0,
		s, c, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) MakeRotationAxis(axis Vector3[T], angle T) *Matrix4[T] {
	// Based on http://www.gamedev.net/reference/articles/article1199.asp
	c := T(math.Cos(float64(angle)))
	s := T(math.Sin(float64(angle)))
	t := 1 - c
	x, y, z := axis.X, axis.Y, axis.Z
	tx, ty := t*x, t*y
	m.Set(
		tx*x+c, tx*y-s*z, tx*z+s*y, 0,
		tx*y+s*z, ty*y+c, ty*z-s*x, 0,
		tx*z-s*y, ty*z+s*x, t*z*z+c, 0,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) MakeScale(v Vector3[T]) *Matrix4[T] {
	m.Set(
		v.X, 0, 0, 0,
		0, v.Y, 0, 0,
		0, 0, v.Z, 0,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) MakeShear(xy, xz, yx, yz, zx, zy T) *Matrix4[T] {
	m.Set(
		1, yx, zx, 0,
		xy, 1, zy, 0,
		xz, yz, 1, 0,
		0, 0, 0, 1,
	)
	return m
}

func (m *Matrix4[T]) Compose(position Vector3[T], rotation Quaternion[T], scale Vector3[T]) *Matrix4[T] {
	x, y, z, w := rotation.X, rotation.Y, rotation.Z, rotation.W
	x2, y2, z2 := x+x, y+y, z+z
	xx, xy, xz := x*x2, x*y2, x*z2
	yy, yz, zz := y*y2, y*z2, z*z2
	wx, wy, wz := w*x2, w*y2, w*z2

	sx, sy, sz := scale.X, scale.Y, scale.Z

	m[0] = (1 - (yy + zz)) * sx
	m[1] = (xy + wz) * sx
	m[2] = (xz - wy) * sx
	m[3] = 0

	m[4] = (xy - wz) * sy
	m[5] = (1 - (xx + zz)) * sy
	m[6] = (yz + wx) * sy
	m[7] = 0

	m[8] = (xz + wy) * sz
	m[9] = (yz - wx) * sz
	m[10] = (1 - (xx + yy)) * sz
	m[11] = 0

	m[12] = position.X
	m[13] = position.Y
	m[14] = position.Z
	m[15] = 1

	return m
}

func (m *Matrix4[T]) Decompose() (position Vector3[T], rotation Quaternion[T], scale Vector3[T]) {
	// based on http://www.geometrictools.com/Documentation/ExtracmulerAngles.pdf
	v1 := NewZeroVector3[T]()
	sx := v1.Set(m[0], m[1], m[2]).Length()
	sy := v1.Set(m[4], m[5], m[6]).Length()
	sz := v1.Set(m[8], m[9], m[10]).Length()

	// if demrmine is negative, we need to invert one scale
	det := m.Demrminant()
	if det < 0 {
		sx = -sx
	}

	position.X = m[12]
	position.Y = m[13]
	position.Z = m[14]

	// scale the rotation part
	m1 := m.Clone()
	invSX := 1 / sx
	invSY := 1 / sy
	invSZ := 1 / sz

	m1[0] *= invSX
	m1[1] *= invSX
	m1[2] *= invSX

	m1[4] *= invSY
	m1[5] *= invSY
	m1[6] *= invSY

	m1[8] *= invSZ
	m1[9] *= invSZ
	m1[10] *= invSZ

	rotation.SetFromRotationMatrix(*m1)

	scale.X = sx
	scale.Y = sy
	scale.Z = sz

	return position, rotation, scale

}

func (m *Matrix4[T]) MakePerspective(left, right, top, bottom, near, far T) *Matrix4[T] {
	x := 2 * near / (right - left)
	y := 2 * near / (top - bottom)

	a := (right + left) / (right - left)
	b := (top + bottom) / (top - bottom)

	var c, d T
	switch coordinateSystem {
	case CoordinateSystemWebGL:
		c = -(far + near) / (far - near)
		d = -(2 * far * near) / (far - near)
	case CoordinateSystemWebGPU:
		c = -far / (far - near)
		d = -far * near / (far - near)
	default:
		panic("invalid coordinate system")

	}

	m.Set(
		x, 0, a, 0,
		0, y, b, 0,
		0, 0, c, d,
		0, 0, -1, 0,
	)

	return m
}

func (m *Matrix4[T]) MakeOrthographic(left, right, top, bottom, near, far T) *Matrix4[T] {
	x := 1 / (right - left)
	y := 1 / (top - bottom)
	z := 1 / (far - near)

	a := (right + left) * x
	b := (top + bottom) * y
	c := (far + near) * z

	m.Set(
		2*x, 0, 0, -a,
		0, 2*y, 0, -b,
		0, 0, -2*z, -c,
		0, 0, 0, 1,
	)

	return m
}

func (m *Matrix4[T]) Equals(matrix Matrix4[T]) bool {
	for i := 0; i < 16; i++ {
		if m[i] != matrix[i] {
			return false
		}
	}
	return true
}

func (m *Matrix4[T]) FromArray(array []T, offset int) *Matrix4[T] {
	for i := 0; i < 16; i++ {
		m[i] = array[i+offset]
	}
	return m
}

func (m *Matrix4[T]) ToArray(array []T, offset int) []T {
	for i := 0; i < 16; i++ {
		array[i+offset] = m[i]
	}
	return array
}
