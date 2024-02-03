package mathx

type SplineCurve struct {
	baseCurve
	points []Vector3
}

func NewSplineCurve(points ...Vector3) *SplineCurve {
	return &SplineCurve{
		points: points,
	}
}

func (c *SplineCurve) Point(t float64) *Vector3 {
	point := &Vector3{}
	p := float64(len(c.points)-1) * t
	pInt := int(p)

	weight := p - float64(pInt)

	var p0, p1, p2, p3 Vector3
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
		CatmullRom(weight, p0.X, p1.X, p2.X, p3.X),
		CatmullRom(weight, p0.Y, p1.Y, p2.Y, p3.Y),
		CatmullRom(weight, p0.Z, p1.Z, p2.Z, p3.Z),
	)

	return point
}

func (c *SplineCurve) PointAt(u float64) *Vector3 {
	return c.Point(u)
}

func (c *SplineCurve) Points(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(float64(d) / float64(divisions))
	}
	return points
}

func (c *SplineCurve) SpacedPoints(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(float64(d) / float64(divisions))
	}

	return points
}

func (c *SplineCurve) LengthsDefault() []float64 {
	return c.Lengths(200)
}

func (c *SplineCurve) Length() float64 {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *SplineCurve) Lengths(divisions int) []float64 {
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

func (c *SplineCurve) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(200)
}
