package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Box3[T constraints.Float] struct {
	Min Vector3[T]
	Max Vector3[T]
}

func NewBox3[T constraints.Float](minVal, maxVal Vector3[T]) *Box3[T] {
	return &Box3[T]{
		Min: minVal,
		Max: maxVal,
	}
}

func (b *Box3[T]) Set(min, max Vector3[T]) *Box3[T] {
	b.Min = min
	b.Max = max
	return b
}

func (b *Box3[T]) SetFromPoints(points ...Vector3[T]) *Box3[T] {
	b.MakeEmpty()
	for _, p := range points {
		b.ExpandByPoint(p)
	}
	return b
}

func (b *Box3[T]) SetFromCenterAndSize(center, size Vector3[T]) *Box3[T] {
	halfSize := size.Clone().MultiplyScalar(0.5)
	b.Min = *center.Clone().Sub(*halfSize)
	b.Max = *center.Clone().Add(*halfSize)
	return b
}

func (b *Box3[T]) Clone() *Box3[T] {
	return NewBox3(b.Min, b.Max)
}

func (b *Box3[T]) Copy(box *Box3[T]) *Box3[T] {
	b.Min = *box.Min.Clone()
	b.Max = *box.Max.Clone()
	return b
}

func (b *Box3[T]) MakeEmpty() *Box3[T] {
	minF, maxF := math.MaxFloat64, -math.MaxFloat64
	minTF, maxTF := T(minF), T(maxF)
	b.Min.Set(minTF, minTF, minTF)
	b.Max.Set(maxTF, maxTF, maxTF)
	return b
}

func (b *Box3[T]) IsEmpty() bool {
	return b.Max.X < b.Min.X || b.Max.Y < b.Min.Y || b.Max.Z < b.Min.Z
}

func (b *Box3[T]) Center() Vector3[T] {
	if b.IsEmpty() {
		return *NewZeroVector3[T]()
	}

	return *AddVector3s(b.Min, b.Max).MultiplyScalar(0.5)
}

func (b *Box3[T]) Size() *Vector3[T] {
	if b.IsEmpty() {
		return NewZeroVector3[T]()
	}

	return SubVector3s(b.Max, b.Min)
}

func (b *Box3[T]) ExpandByPoint(point Vector3[T]) *Box3[T] {
	b.Min = *b.Min.Min(point)
	b.Max = *b.Max.Max(point)
	return b
}

func (b *Box3[T]) ExpandByVector(vector Vector3[T]) *Box3[T] {
	b.Min = *b.Min.Sub(vector)
	b.Max = *b.Max.Add(vector)
	return b
}

func (b *Box3[T]) ExpandByScalar(scalar T) *Box3[T] {
	b.Min = *b.Min.AddScalar(-scalar)
	b.Max = *b.Max.AddScalar(scalar)
	return b
}

func (b *Box3[T]) ContainsPoint(point Vector3[T]) bool {
	return point.X >= b.Min.X && point.X <= b.Max.X &&
		point.Y >= b.Min.Y && point.Y <= b.Max.Y &&
		point.Z >= b.Min.Z && point.Z <= b.Max.Z
}

func (b *Box3[T]) ContainsBox(box Box3[T]) bool {
	return b.Min.X <= box.Min.X && box.Max.X <= b.Max.X && b.Min.Y <= box.Min.Y && box.Max.Y <= b.Max.Y && b.Min.Z <= box.Min.Z && box.Max.Z <= b.Max.Z
}

func (b *Box3[T]) GetParameter(point Vector3[T]) (target Vector3[T]) {
	return *target.Set(
		(point.X-b.Min.X)/(b.Max.X-b.Min.X),
		(point.Y-b.Min.Y)/(b.Max.Y-b.Min.Y),
		(point.Z-b.Min.Z)/(b.Max.Z-b.Min.Z),
	)
}

func (b *Box3[T]) IntersectsBox(box Box3[T]) bool {
	return box.Max.X < b.Min.X || box.Min.X > b.Max.X || box.Max.Y < b.Min.Y || box.Min.Y > b.Max.Y || box.Max.Z < b.Min.Z || box.Min.Z > b.Max.Z
}

func (b *Box3[T]) IntersectsSphere(sphere Sphere[T]) bool {
	_vector := b.ClampPoint(sphere.Center)
	return _vector.DistanceToSquared(sphere.Center) <= (sphere.Radius * sphere.Radius)
}

func (b *Box3[T]) IntersectsPlane(plane Plane[T]) bool {
	var min, max T
	if plane.Normal.X > 0 {
		min = plane.Normal.X * b.Min.X
		max = plane.Normal.X * b.Max.X
	} else {
		min = plane.Normal.X * b.Max.X
		max = plane.Normal.X * b.Min.X
	}

	if plane.Normal.Y > 0 {
		min += plane.Normal.Y * b.Min.Y
		max += plane.Normal.Y * b.Max.Y
	} else {
		min += plane.Normal.Y * b.Max.Y
		max += plane.Normal.Y * b.Min.Y
	}

	if plane.Normal.Z > 0 {
		min += plane.Normal.Z * b.Min.Z
		max += plane.Normal.Z * b.Max.Z
	} else {
		min += plane.Normal.Z * b.Max.Z
		max += plane.Normal.Z * b.Min.Z
	}

	return min <= -plane.Constant && max >= -plane.Constant
}

func (b *Box3[T]) IntersectsTriangle(triangle Triangle[T]) bool {
	if b.IsEmpty() {
		return false
	}

	center := b.Center()
	_extents := *SubVector3s(b.Max, center)

	_v0 := *SubVector3s(triangle.A, center)
	_v1 := *SubVector3s(triangle.B, center)
	_v2 := *SubVector3s(triangle.C, center)

	_f0 := *SubVector3s(_v1, _v0)
	_f1 := *SubVector3s(_v2, _v1)
	_f2 := *SubVector3s(_v0, _v2)

	axes := []T{
		0, -_f0.Z, _f0.Y, 0, -_f1.Z, _f1.Y, 0, -_f2.Z, _f2.Y,
		_f0.Z, 0, -_f0.X, _f1.Z, 0, -_f1.X, _f2.Z, 0, -_f2.X,
		-_f0.Y, _f0.X, 0, -_f1.Y, _f1.X, 0, -_f2.Y, _f2.X, 0,
	}

	if !satForAxes(axes, _v0, _v1, _v2, _extents) {
		return false
	}

	axes = []T{1, 0, 0, 0, 1, 0, 0, 0, 1}
	if !satForAxes(axes, _v0, _v1, _v2, _extents) {
		return false
	}

	_triangleNormal := CrossVector3s(_f0, _f1)
	axes = []T{_triangleNormal.X, _triangleNormal.Y, _triangleNormal.Z}

	return satForAxes(axes, _v0, _v1, _v2, _extents)
}

func (b *Box3[T]) ClampPoint(point Vector3[T]) (target *Vector3[T]) {
	return NewVector3(
		max(b.Min.X, min(b.Max.X, point.X)),
		max(b.Min.Y, min(b.Max.Y, point.Y)),
		max(b.Min.Z, min(b.Max.Z, point.Z)),
	)
}

func (b *Box3[T]) DistanceToPoint(point Vector3[T]) T {
	return b.ClampPoint(point).DistanceTo(point)
}

func (b *Box3[T]) GetBoundingSphere() *Sphere[T] {
	if b.IsEmpty() {
		zero := NewZeroVector3[T]()
		return NewSphere(*zero, 0)
	}

	return NewSphere(b.Center(), b.Size().Length()*0.5)
}

func (b *Box3[T]) Intersect(box Box3[T]) *Box3[T] {
	b.Min = *b.Min.Max(box.Min)
	b.Max = *b.Max.Min(box.Max)

	if b.IsEmpty() {
		b.MakeEmpty()
	}

	return b
}

func (b *Box3[T]) Union(box Box3[T]) *Box3[T] {
	b.Min = *b.Min.Min(box.Min)
	b.Max = *b.Max.Max(box.Max)
	return b
}

func (b *Box3[T]) ApplyMatrix4(matrix Matrix4[T]) *Box3[T] {
	if b.IsEmpty() {
		return b
	}

	b.SetFromPoints(
		*b.Min.Clone().ApplyMatrix4(matrix),
		*b.Min.Clone().Set(b.Min.X, b.Min.Y, b.Max.Z).ApplyMatrix4(matrix),
		*b.Min.Clone().Set(b.Min.X, b.Max.Y, b.Min.Z).ApplyMatrix4(matrix),
		*b.Min.Clone().Set(b.Min.X, b.Max.Y, b.Max.Z).ApplyMatrix4(matrix),
		*b.Max.Clone().Set(b.Max.X, b.Min.Y, b.Min.Z).ApplyMatrix4(matrix),
		*b.Max.Clone().Set(b.Max.X, b.Min.Y, b.Max.Z).ApplyMatrix4(matrix),
		*b.Max.Clone().Set(b.Max.X, b.Max.Y, b.Min.Z).ApplyMatrix4(matrix),
		*b.Max.Clone().ApplyMatrix4(matrix),
	)

	return b
}

func (b *Box3[T]) Translate(offset Vector3[T]) *Box3[T] {
	b.Min = *b.Min.Add(offset)
	b.Max = *b.Max.Add(offset)
	return b
}

func (b *Box3[T]) Equals(box Box3[T]) bool {
	return b.Min.Equals(box.Min) && b.Max.Equals(box.Max)
}

func satForAxes[T constraints.Float](axes []T, v0, v1, v2, extents Vector3[T]) bool {
	testAxis := Vector3[T]{}
	xf, yf, zf := float64(extents.X), float64(extents.Y), float64(extents.Z)
	for i := 0; i <= len(axes)-3; i += 3 {
		testAxis.Set(axes[i], axes[i+1], axes[i+2])
		tXf, tYf, tZf := float64(testAxis.X), float64(testAxis.Y), float64(testAxis.Z)
		r := tXf*math.Abs(tXf)*xf + tYf*math.Abs(tYf)*yf + tZf*math.Abs(tZf)*zf
		p0 := v0.Clone().Dot(testAxis)
		p1 := v1.Clone().Dot(testAxis)
		p2 := v2.Clone().Dot(testAxis)
		if max(-max(p0, p1, p2), min(p0, p1, p2)) > T(r) {
			return false
		}
	}

	return true
}
