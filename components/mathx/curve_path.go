package mathx

import "golang.org/x/exp/constraints"

/**************************************************************
 *	Curved Path - a curve path is simply a array of connected
 *  curves, but retains the api of a curve
 **************************************************************/

type CurvePath[T constraints.Float] struct {
	autoClose    bool
	curves       []Curve[T]
	needsUpdate  bool
	cacheLengths []T
}

func NewCurvePath[T constraints.Float]() *CurvePath[T] {
	return &CurvePath[T]{}
}

func (c *CurvePath[T]) Add(curve Curve[T]) {
	c.curves = append(c.curves, curve)
}

func (c *CurvePath[T]) ClosePath() *CurvePath[T] {
	// Add a line curve if start and end of lines are not connected
	startPoint := c.curves[0].Point(0)
	endPoint := c.curves[len(c.curves)-1].Point(1)

	if !startPoint.Equals(*endPoint) {
		c.curves = append(c.curves, NewLineCurve3(*endPoint, *startPoint))
	}

	return c
}

func (c *CurvePath[T]) Point(t T) *Vector3[T] {
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

func (c *CurvePath[T]) Length() T {
	lens := c.Lengths()
	return lens[len(lens)-1]
}

func (c *CurvePath[T]) UpdateArcLengths() {
	c.needsUpdate = true
	c.cacheLengths = nil
	c.Lengths()
}

func (c *CurvePath[T]) Lengths() []T {
	if c.cacheLengths != nil && len(c.cacheLengths) == len(c.curves) {
		return c.cacheLengths
	}

	var (
		lengths = make([]T, 0)
		sums    T
	)

	for i := 0; i < len(c.curves); i++ {
		sums += c.curves[i].Length()
		lengths = append(lengths, sums)
	}

	c.cacheLengths = lengths
	return lengths
}

func (c *CurvePath[T]) SpacedPoints(divisions int) []*Vector3[T] {
	points := make([]*Vector3[T], 0)

	for i := 0; i <= divisions; i++ {
		points = append(points, c.Point(T(i)/T(divisions)))
	}

	if c.autoClose {
		points = append(points, points[0])
	}

	return points
}

func (c *CurvePath[T]) Points(divisions int) []Vector3[T] {
	points := make([]Vector3[T], 0)
	var last *Vector3[T]

	for i := 0; i < len(c.curves); i++ {
		curve := c.curves[i]
		resolution := divisions

		switch x := curve.(type) {
		case *EllipseCurve[T]:
			resolution = divisions * 2
		case *LineCurve[T]:
			resolution = 1
		case *SplineCurve[T]:
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

func (c *CurvePath[T]) Copy(source *CurvePath[T]) *CurvePath[T] {
	c.curves = make([]Curve[T], len(source.curves))
	copy(c.curves, source.curves)
	c.autoClose = source.autoClose
	return c
}
