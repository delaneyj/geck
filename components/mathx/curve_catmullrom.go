package mathx

import "math"

/**
 * Centripetal CatmullRom Curve - which is useful for avoiding
 * cusps and self-intersections in non-uniform catmull rom curves.
 * http://www.cemyuksel.com/research/catmullrom_param/catmullrom.pdf
 *
 * curve.type accepts centripetal(default), chordal and catmullrom
 * curve.tension is used for catmullrom which defaults to 0.5
 */

/*
Based on an optimized c++ solution in
 - http://stackoverflow.com/questions/9489736/catmull-rom-curve-with-no-cusps-and-no-self-intersections/
 - http://ideone.com/NoEbVM

This CubicPoly class could be used for reusing some variables and calculations,
but for three.js curve use, it could be possible inlined and flatten into a single function call
which can be placed in CurveUtils.
*/

type cubicPoly struct {
	C0, C1, C2, C3 float64
}

type CurveType int

const (
	CurveTypeUnknown CurveType = iota
	CurveTypeCentripetal
	CurveTypeChordal
	CurveTypeCatmullRom
)

func (c *cubicPoly) Init(x0, x1, t0, t1 float64) {
	c.C0 = x0
	c.C1 = t0
	c.C2 = -3*x0 + 3*x1 - 2*t0 - t1
	c.C3 = 2*x0 - 2*x1 + t0 + t1
}

func (c *cubicPoly) InitCatmullRom(x0, x1, x2, x3, tension float64) {
	c.Init(x1, x2, tension*(x2-x0), tension*(x3-x1))
}

func (c *cubicPoly) InitNonuniformCatmullRom(x0, x1, x2, x3, dt0, dt1, dt2 float64) {
	// compute tangents when parameterized in [t1,t2]
	t1 := (x1-x0)/dt0 - (x2-x0)/(dt0+dt1) + (x2-x1)/dt1
	t2 := (x2-x1)/dt1 - (x3-x1)/(dt1+dt2) + (x3-x2)/dt2

	// rescale tangents for parametrization in [0,1]
	t1 *= dt1
	t2 *= dt1

	c.Init(x1, x2, t1, t2)
}

func (c *cubicPoly) Calc(t float64) float64 {
	t2 := t * t
	t3 := t2 * t
	return c.C0 + c.C1*t + c.C2*t2 + c.C3*t3
}

type CatmullRomCurve3 struct {
	baseCurve
	Points    []*Vector3
	Closed    bool
	CurveType CurveType
	Tension   float64
}

func NewCatmullRomCurve3(points []*Vector3, closed bool, curveType CurveType, tension float64) *CatmullRomCurve3 {
	c := &CatmullRomCurve3{
		Points:    points,
		Closed:    closed,
		CurveType: curveType,
		Tension:   tension,
	}
	return c
}

func (c *CatmullRomCurve3) GetPoint(t float64, optionalTarget *Vector3) *Vector3 {
	point := optionalTarget
	points := c.Points
	l := len(points)
	b := 0
	if c.Closed {
		b = 1
	}
	p := float64(l-b) * t
	intPoint := int(p)
	weight := p - float64(intPoint)

	if c.Closed {
		if intPoint <= 0 {
			intPoint += (int(math.Abs(float64(intPoint))/float64(l)) + 1) * l
		}
	} else if weight == 0 && intPoint == l-1 {
		intPoint = l - 2
		weight = 1
	}

	var p0, p3 *Vector3 // 4 points (p1 & p2 defined below)

	if c.Closed || intPoint > 0 {
		p0 = points[(intPoint-1)%l]
	} else {
		// extrapolate first point
		p0 = points[0].Clone().Sub(*points[1]).Add(*points[0])
	}

	p1 := points[intPoint%l]
	p2 := points[(intPoint+1)%l]

	if c.Closed || intPoint+2 < l {
		p3 = points[(intPoint+2)%l]
	} else {
		// extrapolate last point
		p3 = points[l-1].Clone().Sub(*points[l-2]).Add(*points[l-1])
	}

	if c.CurveType == CurveTypeCentripetal || c.CurveType == CurveTypeChordal {
		// init Centripetal / Chordal Catmull-Rom
		pow := 0.25
		if c.CurveType == CurveTypeChordal {
			pow = 0.5
		}
		dt0 := math.Pow(p0.DistanceToSquared(*p1), pow)
		dt1 := math.Pow(p1.DistanceToSquared(*p2), pow)
		dt2 := math.Pow(p2.DistanceToSquared(*p3), pow)

		// safety check for repeated points
		if dt1 < 1e-4 {
			dt1 = 1.0
		}
		if dt0 < 1e-4 {
			dt0 = dt1
		}
		if dt2 < 1e-4 {
			dt2 = dt1
		}

		px := &cubicPoly{}
		py := &cubicPoly{}
		pz := &cubicPoly{}
		px.InitNonuniformCatmullRom(p0.X, p1.X, p2.X, p3.X, dt0, dt1, dt2)
		py.InitNonuniformCatmullRom(p0.Y, p1.Y, p2.Y, p3.Y, dt0, dt1, dt2)
		pz.InitNonuniformCatmullRom(p0.Z, p1.Z, p2.Z, p3.Z, dt0, dt1, dt2)
		point.Set(
			px.Calc(weight),
			py.Calc(weight),
			pz.Calc(weight),
		)
	} else if c.CurveType == CurveTypeCatmullRom {
		px := &cubicPoly{}
		py := &cubicPoly{}
		pz := &cubicPoly{}
		px.InitCatmullRom(p0.X, p1.X, p2.X, p3.X, c.Tension)
		py.InitCatmullRom(p0.Y, p1.Y, p2.Y, p3.Y, c.Tension)
		pz.InitCatmullRom(p0.Z, p1.Z, p2.Z, p3.Z, c.Tension)
		point.Set(
			px.Calc(weight),
			py.Calc(weight),
			pz.Calc(weight),
		)
	}
	return point
}

func (c *CatmullRomCurve3) Copy(source *CatmullRomCurve3) *CatmullRomCurve3 {
	c.Points = make([]*Vector3, len(source.Points))
	for i := 0; i < len(source.Points); i++ {
		c.Points[i] = source.Points[i].Clone()
	}
	c.Closed = source.Closed
	c.CurveType = source.CurveType
	c.Tension = source.Tension
	return c
}
