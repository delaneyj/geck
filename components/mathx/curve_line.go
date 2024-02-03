package mathx

type LineCurve struct {
	baseCurve
	V1, V2 Vector3
}

func NewLineCurve3(v1, v2 Vector3) *LineCurve {
	return &LineCurve{
		V1: v1,
		V2: v2,
	}
}

func (c *LineCurve) Point(t float64) *Vector3 {
	point := c.V2.Clone().Sub(c.V1)
	point.MultiplyScalar(t).Add(c.V1)
	return point
}

func (c *LineCurve) PointAt(u float64) *Vector3 {
	return c.Point(u)
}

func (c *LineCurve) Points(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(float64(d) / float64(divisions))
	}
	return points
}

func (c *LineCurve) SpacedPoints(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(float64(d) / float64(divisions))
	}

	return points
}

func (c *LineCurve) LengthsDefault() []float64 {
	return c.Lengths(200)
}

func (c *LineCurve) Length() float64 {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *LineCurve) Lengths(divisions int) []float64 {
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

func (c *LineCurve) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(200)
}
