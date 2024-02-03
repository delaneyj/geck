package mathx

type CubicBezierCurve struct {
	baseCurve
	V0, V1, V2, V3 Vector3
}

func NewCubicBezierCurve3(v0, v1, v2, v3 Vector3) *CubicBezierCurve {
	return &CubicBezierCurve{
		baseCurve: *newBaseCurve(),
		V0:        v0,
		V1:        v1,
		V2:        v2,
		V3:        v3,
	}
}

func (c *CubicBezierCurve) Point(t float64) *Vector3 {
	v0 := c.V0
	v1 := c.V1
	v2 := c.V2
	v3 := c.V3
	return NewVector3(
		CubicBezier(t, v0.X, v1.X, v2.X, v3.X),
		CubicBezier(t, v0.Y, v1.Y, v2.Y, v3.Y),
		CubicBezier(t, v0.Z, v1.Z, v2.Z, v3.Z),
	)
}

func (c *CubicBezierCurve) PointAt(u float64) *Vector3 {
	return c.Point(u)
}

func (c *CubicBezierCurve) Points(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(float64(d) / float64(divisions))
	}
	return points
}

func (c *CubicBezierCurve) SpacedPoints(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(float64(d) / float64(divisions))
	}

	return points
}

func (c *CubicBezierCurve) LengthsDefault() []float64 {
	return c.Lengths(200)
}

func (c *CubicBezierCurve) Length() float64 {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *CubicBezierCurve) Lengths(divisions int) []float64 {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions {
		return c.cacheArcLengths
	}

	lengths := make([]float64, 0)
	lengths = append(lengths, 0)

	sum := 0.0
	var current, last *Vector3
	last = c.Point(0)
	for p := 1; p <= divisions; p++ {
		current = c.Point(float64(p) / float64(divisions))
		sum += current.DistanceTo(*last)
		lengths = append(lengths, sum)
		last = current
	}

	c.cacheArcLengths = lengths
	return lengths
}

func (c *CubicBezierCurve) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(200)
}
