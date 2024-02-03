package mathx

import "math"

type EllipseCurve struct {
	baseCurve
	aX, aY, xRadius, yRadius, aStartAngle, aEndAngle, aRotation float64
	aClockwise                                                  bool
}

func NewEllipseCurve(aX, aY, xRadius, yRadius, aStartAngle, aEndAngle, aRotation float64, aClockwise bool) *EllipseCurve {
	return &EllipseCurve{
		aX:          aX,
		aY:          aY,
		xRadius:     xRadius,
		yRadius:     yRadius,
		aStartAngle: aStartAngle,
		aEndAngle:   aEndAngle,
		aRotation:   aRotation,
		aClockwise:  aClockwise,
	}
}

func (c *EllipseCurve) Point(t float64) *Vector3 {
	point := &Vector3{}
	twoPi := 2 * math.Pi
	deltaAngle := c.aEndAngle - c.aStartAngle
	samePoints := math.Abs(deltaAngle) < EPSILON64

	for deltaAngle < 0 {
		deltaAngle += twoPi
	}
	for deltaAngle > twoPi {
		deltaAngle -= twoPi
	}

	if deltaAngle < EPSILON64 {
		if samePoints {
			deltaAngle = 0
		} else {
			deltaAngle = twoPi
		}
	}

	if c.aClockwise && !samePoints {
		if deltaAngle == twoPi {
			deltaAngle = -twoPi
		} else {
			deltaAngle = deltaAngle - twoPi
		}
	}

	angle := c.aStartAngle + t*deltaAngle
	x := c.aX + c.xRadius*math.Cos(angle)
	y := c.aY + c.yRadius*math.Sin(angle)

	if c.aRotation != 0 {
		cos := math.Cos(c.aRotation)
		sin := math.Sin(c.aRotation)
		tx := x - c.aX
		ty := y - c.aY
		x = tx*cos - ty*sin + c.aX
		y = tx*sin + ty*cos + c.aY
	}

	point.Set(x, y, 0)
	return point
}

func (c *EllipseCurve) PointAt(u float64) *Vector3 {
	return c.Point(u)
}

func (c *EllipseCurve) Points(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(float64(d) / float64(divisions))
	}
	return points
}

func (c *EllipseCurve) SpacedPoints(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(float64(d) / float64(divisions))
	}
	return points
}

func (c *EllipseCurve) LengthsDefault() []float64 {
	return c.Lengths(c.ArcLengthDivisions)
}

func (c *EllipseCurve) Length() float64 {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *EllipseCurve) Lengths(divisions int) []float64 {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions {
		return c.cacheArcLengths
	}

	lengths := make([]float64, 0)
	sum := 0.0
	var current, last *Vector3
	last = c.Point(0)
	lengths = append(lengths, 0)

	for p := 1; p <= divisions; p++ {
		current = c.Point(float64(p) / float64(divisions))
		sum += current.DistanceTo(*last)

		lengths = append(lengths, sum)
		last = current
	}

	c.cacheArcLengths = lengths
	return lengths
}

func (c *EllipseCurve) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(c.ArcLengthDivisions)
}

func (c *EllipseCurve) Copy(source *EllipseCurve) *EllipseCurve {
	c.baseCurve.Copy(source.baseCurve)
	c.aX = source.aX
	c.aY = source.aY
	c.xRadius = source.xRadius
	c.yRadius = source.yRadius
	c.aStartAngle = source.aStartAngle
	c.aEndAngle = source.aEndAngle
	c.aClockwise = source.aClockwise
	c.aRotation = source.aRotation
	return c
}
