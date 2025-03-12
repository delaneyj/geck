package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type EllipseCurve[T constraints.Float] struct {
	baseCurve[T]
	aX, aY, xRadius, yRadius, aStartAngle, aEndAngle, aRotation T
	aClockwise                                                  bool
}

func NewEllipseCurve[T constraints.Float](aX, aY, xRadius, yRadius, aStartAngle, aEndAngle, aRotation T, aClockwise bool) *EllipseCurve[T] {
	return &EllipseCurve[T]{
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

func (c *EllipseCurve[T]) Point(t T) *Vector3[T] {
	point := &Vector3[T]{}
	twoPi := T(2 * math.Pi)
	deltaAngle := c.aEndAngle - c.aStartAngle
	samePoints := math.Abs(float64(deltaAngle)) < EPSILON

	for deltaAngle < 0 {
		deltaAngle += twoPi
	}
	for deltaAngle > twoPi {
		deltaAngle -= twoPi
	}

	if deltaAngle < EPSILON {
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
	af := float64(angle)
	x := c.aX + c.xRadius*T(math.Cos(af))
	y := c.aY + c.yRadius*T(math.Sin(af))

	if c.aRotation != 0 {
		rf := float64(c.aRotation)
		cos := T(math.Cos(rf))
		sin := T(math.Sin(rf))
		tx := x - c.aX
		ty := y - c.aY
		x = tx*cos - ty*sin + c.aX
		y = tx*sin + ty*cos + c.aY
	}

	point.Set(x, y, 0)
	return point
}

func (c *EllipseCurve[T]) PointAt(u T) *Vector3[T] {
	return c.Point(u)
}

func (c *EllipseCurve[T]) Points(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.Point(T(d) / T(divisions))
	}
	return points
}

func (c *EllipseCurve[T]) SpacedPoints(divisions int) []Vector3[T] {
	points := make([]Vector3[T], divisions)
	for d := 0; d <= divisions; d++ {
		points[d] = *c.PointAt(T(d) / T(divisions))
	}
	return points
}

func (c *EllipseCurve[T]) LengthsDefault() []T {
	return c.Lengths(c.ArcLengthDivisions)
}

func (c *EllipseCurve[T]) Length() T {
	lens := c.LengthsDefault()
	return lens[len(lens)-1]
}

func (c *EllipseCurve[T]) Lengths(divisions int) []T {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions {
		return c.cacheArcLengths
	}

	lengths := make([]T, 0)
	var sum T
	var current, last *Vector3[T]
	last = c.Point(0)
	lengths = append(lengths, 0)

	for p := 1; p <= divisions; p++ {
		current = c.Point(T(p) / T(divisions))
		sum += current.DistanceTo(*last)

		lengths = append(lengths, sum)
		last = current
	}

	c.cacheArcLengths = lengths
	return lengths
}

func (c *EllipseCurve[T]) UpdateArcLengths() {
	c.cacheArcLengths = nil
	c.Lengths(c.ArcLengthDivisions)
}

func (c *EllipseCurve[T]) Copy(source *EllipseCurve[T]) *EllipseCurve[T] {
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
