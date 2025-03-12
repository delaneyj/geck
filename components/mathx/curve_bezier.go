package mathx

import "golang.org/x/exp/constraints"

type CubicBezierCurve[T constraints.Float] struct {
	baseCurve[T]
	V0, V1, V2, V3 Vector3[T]
}

func NewCubicBezierCurve3[T constraints.Float](v0, v1, v2, v3 Vector3[T]) *CubicBezierCurve[T] {
	return &CubicBezierCurve[T]{
		baseCurve: *newBaseCurve[T](),
		V0:        v0,
		V1:        v1,
		V2:        v2,
		V3:        v3,
	}
}

func (c *CubicBezierCurve[T]) Point(t T) *Vector3[T] {
	v0 := c.V0
	v1 := c.V1
	v2 := c.V2
	v3 := c.V3
	return NewVector3[T](
		CubicBezier[T](t, v0.X, v1.X, v2.X, v3.X),
		CubicBezier[T](t, v0.Y, v1.Y, v2.Y, v3.Y),
		CubicBezier[T](t, v0.Z, v1.Z, v2.Z, v3.Z),
	)
}

func (c *CubicBezierCurve[T]) PointAt(u T) *Vector3[T] {
	return c.Point(u)
}

func (c *CubicBezierCurve[T]) Points(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(T(d) / T(divisions))
	}
	return points
}

func (c *CubicBezierCurve[T]) SpacedPoints(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(T(d) / T(divisions))
	}

	return points
}

func (c *CubicBezierCurve[T]) LengthsDefault() []T {
	return c.Lengths(200)
}

func (c *CubicBezierCurve[T]) Length() T {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *CubicBezierCurve[T]) Lengths(divisions int) []T {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions {
		return c.cacheArcLengths
	}

	lengths := make([]T, 0)
	lengths = append(lengths, 0)

	var sum T
	var current, last *Vector3[T]
	last = c.Point(0)
	for p := 1; p <= divisions; p++ {
		current = c.Point(T(p) / T(divisions))
		sum += current.DistanceTo(*last)
		lengths = append(lengths, sum)
		last = current
	}

	c.cacheArcLengths = lengths
	return lengths
}

func (c *CubicBezierCurve[T]) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(200)
}
