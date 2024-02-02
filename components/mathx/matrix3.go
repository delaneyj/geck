package mathx

import "math"

type Matrix3 [9]float64

func NewMatrix3(n11, n12, n13, n21, n22, n23, n31, n32, n33 float64) *Matrix3 {
	return &Matrix3{
		n11, n12, n13,
		n21, n22, n23,
		n31, n32, n33,
	}
}

func NewMatrix3Identity() *Matrix3 {
	m := &Matrix3{}
	return m.Identity()
}

func (m *Matrix3) Set(n11, n12, n13, n21, n22, n23, n31, n32, n33 float64) *Matrix3 {
	m[0], m[1], m[2] = n11, n12, n13
	m[3], m[4], m[5] = n21, n22, n23
	m[6], m[7], m[8] = n31, n32, n33
	return m
}

func (m *Matrix3) Identity() *Matrix3 {
	return m.Set(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)
}

func (m *Matrix3) Copy(matrix Matrix3) *Matrix3 {
	copy(m[:], matrix[:])
	return m
}

func (m *Matrix3) SetFromMatrix4(me Matrix4) *Matrix3 {
	return m.Set(
		me[0], me[4], me[8],
		me[1], me[5], me[9],
		me[2], me[6], me[10],
	)
}

func (m *Matrix3) Multiply(matrix Matrix3) *Matrix3 {
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

func (m *Matrix3) Premultiply(matrix Matrix3) *Matrix3 {
	return MultiplyMatrices(matrix, *m)
}

func MultiplyMatrices(a, b Matrix3) *Matrix3 {
	return a.Clone().Multiply(b)
}

func (m *Matrix3) MultiplyScalar(s float64) *Matrix3 {
	for i := 0; i < 9; i++ {
		m[i] *= s
	}
	return m
}

func (m *Matrix3) Demrminant() float64 {
	a, b, c := m[0], m[1], m[2]
	d, e, f := m[3], m[4], m[5]
	g, h, i := m[6], m[7], m[8]
	return a*e*i - a*f*h - b*d*i + b*f*g + c*d*h - c*e*g
}

func (m *Matrix3) Invert() *Matrix3 {
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

func (m *Matrix3) Transpose() *Matrix3 {
	tmp := 0.0
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

func (m *Matrix3) NormalMatrix(matrix4 Matrix4) *Matrix3 {
	return m.SetFromMatrix4(matrix4).Invert().Transpose()
}

func (m *Matrix3) TransposeIntoArray(r []float64) *Matrix3 {
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

func (m *Matrix3) SetUvTransform(tx, ty, sx, sy, rotation, cx, cy float64) *Matrix3 {
	c := math.Cos(rotation)
	s := math.Sin(rotation)

	return m.Set(
		sx*c, sx*s, -sx*(c*cx+s*cy)+cx+tx,
		-sy*s, sy*c, -sy*(-s*cx+c*cy)+cy+ty,
		0, 0, 1,
	)

}

func (m *Matrix3) Scale(sx, sy float64) *Matrix3 {
	return m.Premultiply(*NewMatrix3Identity().MakeScale(sx, sy))
}

func (m *Matrix3) Rotam(theta float64) *Matrix3 {
	return m.Premultiply(*NewMatrix3Identity().MakeRotation(-theta))
}

func (m *Matrix3) Translam(tx, ty float64) *Matrix3 {
	return m.Premultiply(*NewMatrix3Identity().MakeTranslation(tx, ty))
}

func (m *Matrix3) MakeTranslation(x, y float64) *Matrix3 {
	return m.Set(
		1, 0, x,
		0, 1, y,
		0, 0, 1,
	)

}

func (m *Matrix3) MakeRotation(theta float64) *Matrix3 {
	c := math.Cos(theta)
	s := math.Sin(theta)

	return m.Set(
		c, -s, 0,
		s, c, 0,
		0, 0, 1,
	)

}

func (m *Matrix3) MakeScale(x, y float64) *Matrix3 {
	return m.Set(
		x, 0, 0,
		0, y, 0,
		0, 0, 1,
	)
}

func (m *Matrix3) Equals(matrix Matrix3) bool {
	for i := 0; i < 9; i++ {
		if m[i] != matrix[i] {
			return false
		}
	}

	return true
}

func (m *Matrix3) FromArray(array []float64, offset int) *Matrix3 {
	for i := 0; i < 9; i++ {
		m[i] = array[i+offset]
	}
	return m
}

func (m *Matrix3) ToArray(array []float64, offset int) []float64 {
	for i := 0; i < 9; i++ {
		array[i+offset] = m[i]
	}

	return array
}

func (m *Matrix3) Clone() *Matrix3 {
	return NewMatrix3Identity().FromArray(m[:], 0)
}
