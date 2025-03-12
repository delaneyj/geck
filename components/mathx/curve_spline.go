package mathx

import "golang.org/x/exp/constraints"

type SplineCurve[T constraints.Float] struct {
	baseCurve[T]
	points []Vector3[T]
}

func NewSplineCurve[T constraints.Float](points ...Vector3[T]) *SplineCurve[T] {
	return &SplineCurve[T]{
		points: points,
	}
}

func (c *SplineCurve[T]) Point(t T) *Vector3[T] {
	point := &Vector3[T]{}
	p := T(len(c.points)-1) * t
	pInt := int(p)

	weight := p - T(pInt)

	var p0, p1, p2, p3 Vector3[T]
	if pInt == 0 {
		p0 = c.points[pInt]
	} else {
		p0 = c.points[pInt-1]
	}
	p1 = c.points[pInt]
	if pInt > len(c.points)-2 {
		p2 = c.points[len(c.points)-1]
	} else {
		p2 = c.points[pInt+1]
	}
	if pInt > len(c.points)-3 {
		p3 = c.points[len(c.points)-1]
	} else {
		p3 = c.points[pInt+2]
	}

	point.Set(
		CatmullRom[T](weight, p0.X, p1.X, p2.X, p3.X),
		CatmullRom[T](weight, p0.Y, p1.Y, p2.Y, p3.Y),
		CatmullRom[T](weight, p0.Z, p1.Z, p2.Z, p3.Z),
	)

	return point
}

func (c *SplineCurve[T]) PointAt(u T) *Vector3[T] {
	return c.Point(u)
}

func (c *SplineCurve[T]) Points(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(T(d) / T(divisions))
	}
	return points
}

func (c *SplineCurve[T]) SpacedPoints(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(T(d) / T(divisions))
	}

	return points
}

func (c *SplineCurve[T]) LengthsDefault() []T {
	return c.Lengths(200)
}

func (c *SplineCurve[T]) Length() T {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *SplineCurve[T]) Lengths(divisions int) []T {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions {
		return c.cacheArcLengths
	}

	lengths := make([]T, 0)
	lengths = append(lengths, 0)
	divF := T(divisions)

	var sum T
	var current, last *Vector3[T]
	last = c.Point(0)
	for p := 1; p <= divisions; p++ {
		current = c.Point(T(p) / divF)
		sum += current.DistanceTo(*last)
		lengths = append(lengths, sum)
		last = current
	}

	c.cacheArcLengths = lengths
	return lengths
}

func (c *SplineCurve[T]) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(200)
}
