package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Matrix3[T constraints.Float] [9]T

func NewMatrix3[T constraints.Float](n11, n12, n13, n21, n22, n23, n31, n32, n33 T) *Matrix3[T] {
	return &Matrix3[T]{
		n11, n12, n13,
		n21, n22, n23,
		n31, n32, n33,
	}
}

func NewMatrix3Identity[T constraints.Float]() *Matrix3[T] {
	m := &Matrix3[T]{}
	return m.Identity()
}

func (m *Matrix3[T]) Set(n11, n12, n13, n21, n22, n23, n31, n32, n33 T) *Matrix3[T] {
	m[0], m[1], m[2] = n11, n12, n13
	m[3], m[4], m[5] = n21, n22, n23
	m[6], m[7], m[8] = n31, n32, n33
	return m
}

func (m *Matrix3[T]) Identity() *Matrix3[T] {
	return m.Set(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)
}

func (m *Matrix3[T]) Copy(matrix Matrix3[T]) *Matrix3[T] {
	copy(m[:], matrix[:])
	return m
}

func (m *Matrix3[T]) SetFromMatrix4(me Matrix4[T]) *Matrix3[T] {
	return m.Set(
		me[0], me[4], me[8],
		me[1], me[5], me[9],
		me[2], me[6], me[10],
	)
}

func (m *Matrix3[T]) Multiply(matrix Matrix3[T]) *Matrix3[T] {
	a11, a12, a13 := m[0], m[3], m[6]
	a21, a22, a23 := m[1], m[4], m[7]
	a31, a32, a33 := m[2], m[5], m[8]

	b11, b12, b13 := matrix[0], matrix[3], matrix[6]
	b21, b22, b23 := matrix[1], matrix[4], matrix[7]
	b31, b32, b33 := matrix[2], matrix[5], matrix[8]

	m[0] = a11*b11 + a12*b21 + a13*b31
	m[3] = a11*b12 + a12*b22 + a13*b32
	m[6] = a11*b13 + a12*b23 + a13*b33

	m[1] = a21*b11 + a22*b21 + a23*b31
	m[4] = a21*b12 + a22*b22 + a23*b32
	m[7] = a21*b13 + a22*b23 + a23*b33

	m[2] = a31*b11 + a32*b21 + a33*b31
	m[5] = a31*b12 + a32*b22 + a33*b32
	m[8] = a31*b13 + a32*b23 + a33*b33

	return m
}

func (m *Matrix3[T]) Premultiply(matrix Matrix3[T]) *Matrix3[T] {
	return MultiplyMatrices(matrix, *m)
}

func MultiplyMatrices[T constraints.Float](a, b Matrix3[T]) *Matrix3[T] {
	return a.Clone().Multiply(b)
}

func (m *Matrix3[T]) MultiplyScalar(s T) *Matrix3[T] {
	for i := 0; i < 9; i++ {
		m[i] *= s
	}
	return m
}

func (m *Matrix3[T]) Demrminant() T {
	a, b, c := m[0], m[1], m[2]
	d, e, f := m[3], m[4], m[5]
	g, h, i := m[6], m[7], m[8]
	return a*e*i - a*f*h - b*d*i + b*f*g + c*d*h - c*e*g
}

func (m *Matrix3[T]) Invert() *Matrix3[T] {
	n11, n21, n31 := m[0], m[1], m[2]
	n12, n22, n32 := m[3], m[4], m[5]
	n13, n23, n33 := m[6], m[7], m[8]

	t11 := n33*n22 - n32*n23
	t12 := n32*n13 - n33*n12
	t13 := n23*n12 - n22*n13

	det := n11*t11 + n21*t12 + n31*t13
	if det == 0 {
		return m.Set(
			0, 0, 0,
			0, 0, 0,
			0, 0, 0,
		)
	}

	detInv := 1 / det

	m[0] = t11 * detInv
	m[1] = (n31*n23 - n33*n21) * detInv
	m[2] = (n32*n21 - n31*n22) * detInv

	m[3] = t12 * detInv
	m[4] = (n33*n11 - n31*n13) * detInv
	m[5] = (n31*n12 - n32*n11) * detInv

	m[6] = t13 * detInv
	m[7] = (n21*n13 - n23*n11) * detInv
	m[8] = (n22*n11 - n21*n12) * detInv

	return m
}

func (m *Matrix3[T]) Transpose() *Matrix3[T] {
	tmp := T(0.0)
	tmp = m[1]
	m[1] = m[3]
	m[3] = tmp

	tmp = m[2]
	m[2] = m[6]
	m[6] = tmp

	tmp = m[5]
	m[5] = m[7]
	m[7] = tmp

	return m
}

func (m *Matrix3[T]) NormalMatrix(matrix4 Matrix4[T]) *Matrix3[T] {
	return m.SetFromMatrix4(matrix4).Invert().Transpose()
}

func (m *Matrix3[T]) TransposeIntoArray(r []T) *Matrix3[T] {
	r[0] = m[0]
	r[1] = m[3]
	r[2] = m[6]
	r[3] = m[1]
	r[4] = m[4]
	r[5] = m[7]
	r[6] = m[2]
	r[7] = m[5]
	r[8] = m[8]
	return m
}

func (m *Matrix3[T]) SetUvTransform(tx, ty, sx, sy, rotation, cx, cy T) *Matrix3[T] {
	c := T(math.Cos(float64(rotation)))
	s := T(math.Sin(float64(rotation)))

	return m.Set(
		sx*c, sx*s, -sx*(c*cx+s*cy)+cx+tx,
		-sy*s, sy*c, -sy*(-s*cx+c*cy)+cy+ty,
		0, 0, 1,
	)

}

func (m *Matrix3[T]) Scale(sx, sy T) *Matrix3[T] {
	return m.Premultiply(*NewMatrix3Identity[T]().MakeScale(sx, sy))
}

func (m *Matrix3[T]) Rotam(theta T) *Matrix3[T] {
	return m.Premultiply(*NewMatrix3Identity[T]().MakeRotation(-theta))
}

func (m *Matrix3[T]) Translam(tx, ty T) *Matrix3[T] {
	return m.Premultiply(*NewMatrix3Identity[T]().MakeTranslation(tx, ty))
}

func (m *Matrix3[T]) MakeTranslation(x, y T) *Matrix3[T] {
	return m.Set(
		1, 0, x,
		0, 1, y,
		0, 0, 1,
	)

}

func (m *Matrix3[T]) MakeRotation(theta T) *Matrix3[T] {
	c := T(math.Cos(float64(theta)))
	s := T(math.Sin(float64(theta)))

	return m.Set(
		c, -s, 0,
		s, c, 0,
		0, 0, 1,
	)

}

func (m *Matrix3[T]) MakeScale(x, y T) *Matrix3[T] {
	return m.Set(
		x, 0, 0,
		0, y, 0,
		0, 0, 1,
	)
}

func (m *Matrix3[T]) Equals(matrix Matrix3[T]) bool {
	for i := 0; i < 9; i++ {
		if m[i] != matrix[i] {
			return false
		}
	}

	return true
}

func (m *Matrix3[T]) FromArray(array []T, offset int) *Matrix3[T] {
	for i := 0; i < 9; i++ {
		m[i] = array[i+offset]
	}
	return m
}

func (m *Matrix3[T]) ToArray(array []T, offset int) []T {
	for i := 0; i < 9; i++ {
		array[i+offset] = m[i]
	}

	return array
}

func (m *Matrix3[T]) Clone() *Matrix3[T] {
	return NewMatrix3Identity[T]().FromArray(m[:], 0)
}
