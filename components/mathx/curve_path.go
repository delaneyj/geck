package mathx

/**************************************************************
 *	Curved Path - a curve path is simply a array of connected
 *  curves, but retains the api of a curve
 **************************************************************/

type CurvePath struct {
	autoClose    bool
	curves       []Curve
	needsUpdate  bool
	cacheLengths []float64
}

func NewCurvePath() *CurvePath {
	return &CurvePath{}
}

func (c *CurvePath) Add(curve Curve) {
	c.curves = append(c.curves, curve)
}

func (c *CurvePath) ClosePath() *CurvePath {
	// Add a line curve if start and end of lines are not connected
	startPoint := c.curves[0].Point(0)
	endPoint := c.curves[len(c.curves)-1].Point(1)

	if !startPoint.Equals(*endPoint) {
		c.curves = append(c.curves, NewLineCurve3(*endPoint, *startPoint))
	}

	return c
}

func (c *CurvePath) Point(t float64) *Vector3 {
	d := t * c.Length()
	curveLengths := c.Lengths()
	i := 0

	for i < len(curveLengths) {
		if curveLengths[i] >= d {
			diff := curveLengths[i] - d
			curve := c.curves[i]
			segmentLength := curve.Length()
			u := 1 - diff/segmentLength
			return curve.PointAt(u)
		}
		i++
	}
	return nil
}

func (c *CurvePath) Length() float64 {
	lens := c.Lengths()
	return lens[len(lens)-1]
}

func (c *CurvePath) UpdateArcLengths() {
	c.needsUpdate = true
	c.cacheLengths = nil
	c.Lengths()
}

func (c *CurvePath) Lengths() []float64 {
	if c.cacheLengths != nil && len(c.cacheLengths) == len(c.curves) {
		return c.cacheLengths
	}

	lengths := make([]float64, 0)
	sums := 0.0

	for i := 0; i < len(c.curves); i++ {
		sums += c.curves[i].Length()
		lengths = append(lengths, sums)
	}

	c.cacheLengths = lengths
	return lengths
}

func (c *CurvePath) SpacedPoints(divisions int) []*Vector3 {
	points := make([]*Vector3, 0)

	for i := 0; i <= divisions; i++ {
		points = append(points, c.Point(float64(i)/float64(divisions)))
	}

	if c.autoClose {
		points = append(points, points[0])
	}

	return points
}

func (c *CurvePath) Points(divisions int) []Vector3 {
	points := make([]Vector3, 0)
	var last *Vector3

	for i := 0; i < len(c.curves); i++ {
		curve := c.curves[i]
		resolution := divisions

		switch x := curve.(type) {
		case *EllipseCurve:
			resolution = divisions * 2
		case *LineCurve:
			resolution = 1
		case *SplineCurve:
			resolution = divisions * len(x.points)
		default:
			panic("unknown curve type")
		}

		pts := curve.Points(resolution)

		for j := 0; j < len(pts); j++ {
			point := pts[j]

			if last != nil && last.Equals(point) {
				continue
			}

			points = append(points, point)
			last = &point
		}

	}

	if c.autoClose && len(points) > 1 && !points[len(points)-1].Equals(points[0]) {
		points = append(points, points[0])

	}

	return points
}

func (c *CurvePath) Copy(source *CurvePath) *CurvePath {
	c.curves = make([]Curve, len(source.curves))
	copy(c.curves, source.curves)
	c.autoClose = source.autoClose
	return c
}
