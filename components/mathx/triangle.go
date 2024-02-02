package mathx

import "math"

type Triangle struct {
	A, B, C Vector3
}

func NewTriangle(a, b, c Vector3) Triangle {
	return Triangle{A: a, B: b, C: c}
}

func (t *Triangle) Normal() *Vector3 {
	target := SubVector3s(t.C, t.B)
	v0 := SubVector3s(t.A, t.B)
	target.Cross(*v0)

	targetLengthSq := target.LengthSq()
	if targetLengthSq > 0 {
		return target.MultiplyScalar(1 / math.Sqrt(targetLengthSq))
	}

	return target.Set(0, 0, 0)
}

func (t *Triangle) Barycoord(point Vector3) *Vector3 {
	v0 := SubVector3s(t.C, t.A)
	v1 := SubVector3s(t.B, t.A)
	v2 := SubVector3s(point, t.A)

	dot00 := v0.Dot(*v0)
	dot01 := v0.Dot(*v1)
	dot02 := v0.Dot(*v2)
	dot11 := v1.Dot(*v1)
	dot12 := v1.Dot(*v2)

	denom := (dot00*dot11 - dot01*dot01)

	if denom == 0 {
		return &Vector3{X: 0, Y: 0, Z: 0}
	}

	invDenom := 1 / denom
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return &Vector3{X: 1 - u - v, Y: v, Z: u}
}

func (t *Triangle) Interpolation(point, p1, p2, p3, v1, v2, v3 Vector3) *Vector3 {
	target := &Vector3{}
	if t.Barycoord(point) == nil {
		target.X = 0
		target.Y = 0
		target.Z = 0
		return target
	}

	target.SetScalar(0)
	target.AddScaledVector(v1, t.Barycoord(point).X)
	target.AddScaledVector(v2, t.Barycoord(point).Y)
	target.AddScaledVector(v3, t.Barycoord(point).Z)

	return target
}

func (t *Triangle) ContainsPoint(point Vector3) bool {
	if t.Barycoord(point) == nil {
		return false
	}

	return t.Barycoord(point).X >= 0 && t.Barycoord(point).Y >= 0 && (t.Barycoord(point).X+t.Barycoord(point).Y) <= 1
}

func (t *Triangle) IsFrontFacing(direction Vector3) bool {
	v0 := SubVector3s(t.C, t.B)
	v1 := SubVector3s(t.A, t.B)

	return v0.Cross(*v1).Dot(direction) < 0
}

func (t *Triangle) Set(a, b, c Vector3) *Triangle {
	t.A = a
	t.B = b
	t.C = c
	return t
}

func (t *Triangle) SetFromPointsAndIndices(points []Vector3, i0, i1, i2 int) *Triangle {
	t.A = points[i0]
	t.B = points[i1]
	t.C = points[i2]
	return t
}

func (t *Triangle) Clone() Triangle {
	return NewTriangle(t.A, t.B, t.C)
}

func (t *Triangle) Copy(triangle Triangle) *Triangle {
	t.A = triangle.A
	t.B = triangle.B
	t.C = triangle.C
	return t
}

func (t *Triangle) Area() float64 {
	v0 := SubVector3s(t.C, t.B)
	v1 := SubVector3s(t.A, t.B)
	return v0.Cross(*v1).Length() * 0.5
}

func (t *Triangle) Midpoint() *Vector3 {
	return AddVector3s(t.A, t.B).Add(t.C).MultiplyScalar(1.0 / 3)
}

func (t *Triangle) Plane() *Plane {
	p := &Plane{}
	return p.SetFromCoplanarPoints(t.A, t.B, t.C)
}

func (t *Triangle) ClosestPointToPoint(p Vector3) *Vector3 {
	a, b, c := t.A, t.B, t.C
	vab := SubVector3s(b, a)
	vac := SubVector3s(c, a)
	vap := SubVector3s(p, a)
	d1 := vab.Dot(*vap)
	d2 := vac.Dot(*vap)
	if d1 <= 0 && d2 <= 0 {
		return a.Clone()
	}

	vbp := SubVector3s(p, b)
	d3 := vab.Dot(*vbp)
	d4 := vac.Dot(*vbp)
	if d3 >= 0 && d4 <= d3 {
		return b.Clone()
	}

	vcp := SubVector3s(p, c)
	d5 := vab.Dot(*vcp)
	d6 := vac.Dot(*vcp)
	if d6 >= 0 && d5 <= d6 {
		return c.Clone()
	}

	vbc := SubVector3s(c, b)
	va := d3*d6 - d5*d4
	if va <= 0 && (d4-d3) >= 0 && (d5-d6) >= 0 {
		w := (d4 - d3) / ((d4 - d3) + (d5 - d6))
		return b.Clone().AddScaledVector(*vbc, w)
	}

	vb := d5*d2 - d1*d6
	if vb <= 0 && d2 >= 0 && d6 <= 0 {
		w := d2 / (d2 - d6)
		return a.Clone().AddScaledVector(*vac, w)
	}

	vc := d1*d4 - d3*d2
	if vc <= 0 && (d1-d3) >= 0 && (d3-d4) >= 0 {
		v := (d1 - d3) / ((d1 - d3) + (d3 - d4))
		return a.Clone().AddScaledVector(*vab, v)
	}

	denom := 1 / (va + vb + vc)
	v := vb * denom
	w := vc * denom
	return a.Clone().AddScaledVector(*vab, v).AddScaledVector(*vac, w)
}

func (t *Triangle) Equals(triangle Triangle) bool {
	return triangle.A.Equals(t.A) && triangle.B.Equals(t.B) && triangle.C.Equals(t.C)
}
