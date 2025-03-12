package mathx

import "golang.org/x/exp/constraints"

type LineCurve[T constraints.Float] struct {
	baseCurve[T]
	V1, V2 Vector3[T]
}

func NewLineCurve3[T constraints.Float](v1, v2 Vector3[T]) *LineCurve[T] {
	return &LineCurve[T]{
		V1: v1,
		V2: v2,
	}
}

func (c *LineCurve[T]) Point(t T) *Vector3[T] {
	point := c.V2.Clone().Sub(c.V1)
	point.MultiplyScalar(t).Add(c.V1)
	return point
}

func (c *LineCurve[T]) PointAt(u T) *Vector3[T] {
	return c.Point(u)
}

func (c *LineCurve[T]) Points(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(T(d) / T(divisions))
	}
	return points
}

func (c *LineCurve[T]) SpacedPoints(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(T(d) / T(divisions))
	}

	return points
}

func (c *LineCurve[T]) LengthsDefault() []T {
	return c.Lengths(200)
}

func (c *LineCurve[T]) Length() T {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *LineCurve[T]) Lengths(divisions int) []T {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions {
		return c.cacheArcLengths
	}

	lengths := make([]T, 0)
	lengths = append(lengths, 0)

	var (
		sum           T
		current, last *Vector3[T]
	)
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

func (c *LineCurve[T]) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(200)
}
