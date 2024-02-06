package mathx

import "math"

/**
 * Extensible curve object.
 *
 * Some common of curve methods:
 * .getPoint( t, optionalTarget ), .getTangent( t, optionalTarget )
 * .getPointAt( u, optionalTarget ), .getTangentAt( u, optionalTarget )
 * .getPoints(), .getSpacedPoints()
 * .getLength()
 * .updateArcLengths()
 *
 * This following curves inherit from THREE.Curve:
 *
 * -- 2D curves --
 * THREE.ArcCurve
 * THREE.CubicBezierCurve
 * THREE.EllipseCurve
 * THREE.LineCurve
 * THREE.QuadraticBezierCurve
 * THREE.SplineCurve
 *
 * -- 3D curves --
 * THREE.CatmullRomCurve3
 * THREE.CubicBezierCurve3
 * THREE.LineCurve3
 * THREE.QuadraticBezierCurve3
 *
 * A series of curves can be represented as a THREE.CurvePath.
 *
 **/

type Curve interface {
	Point(t float64) *Vector3
	PointAt(u float64) *Vector3
	Points(divisions int) []Vector3
	SpacedPoints(divisions int) []Vector3
	LengthsDefault() []float64
	Length() float64
	Lengths(divisions int) []float64
	UpdateArcLengths()
}

type baseCurve struct {
	ArcLengthDivisions int
	cacheArcLengths    []float64
	needsUpdate        bool
}

func newBaseCurve() *baseCurve {
	return &baseCurve{
		ArcLengthDivisions: 200,
	}
}

// Virtual base class method to overwrite and implement in subclasses
//   - t [0 .. 1]
func (c *baseCurve) Point(t float64) *Vector3 {
	panic("GetPoint() not implemented")
}

// Get point at relative position in curve according to arc length
// - u [0 .. 1]
func (c *baseCurve) PointAt(u float64) *Vector3 {
	t := c.uToTmapping(u)
	return c.Point(t)
}

// Get sequence of points using getPoint( t )
func (c *baseCurve) Points(divisions int) []Vector3 {
	points := make([]Vector3, divisions)

	d := 0
	for d <= divisions {
		points[d] = *c.Point(float64(d) / float64(divisions))
		d++
	}
	return points
}

// Get sequence of points using getPointAt( u )
func (c *baseCurve) SpacedPoints(divisions int) []Vector3 {
	points := make([]Vector3, divisions)
	d := 0
	for d <= divisions {
		points[d] = *c.PointAt(float64(d) / float64(divisions))
		d++
	}
	return points
}

func (c *baseCurve) LengthsDefault() []float64 {
	return c.Lengths(c.ArcLengthDivisions)
}

// Get total curve arc length
func (c *baseCurve) Length() float64 {
	lengths := c.LengthsDefault()
	return lengths[len(lengths)-1]
}

// Get list of cumulative segment lengths
func (c *baseCurve) Lengths(divisions int) []float64 {
	if c.cacheArcLengths != nil && len(c.cacheArcLengths) == divisions+1 && !c.needsUpdate {
		return c.cacheArcLengths
	}
	c.needsUpdate = false
	last := c.Point(0)
	sum := 0.0
	cache := []float64{0}
	current := last
	for p := 1; p <= divisions; p++ {
		current = c.Point(float64(p) / float64(divisions))
		sum += current.DistanceTo(*last)
		cache = append(cache, sum)
		last = current
	}
	c.cacheArcLengths = cache
	return cache
}

func (c *baseCurve) UpdateArcLengths() {
	c.needsUpdate = true
	c.Lengths(c.ArcLengthDivisions)
}

// Given u ( 0 .. 1 ), get a t to find p. This gives you points which are equidistant
func (c *baseCurve) uToTmapping(u float64) float64 {
	arcLengths := c.Lengths(c.ArcLengthDivisions)
	var i int
	il := len(arcLengths)
	var targetArcLength float64
	if targetArcLength == 0 {
		targetArcLength = u * arcLengths[il-1]
	}
	low := 0
	high := il - 1
	var comparison float64
	for low <= high {
		i = low + (high-low)/2
		comparison = arcLengths[i] - targetArcLength
		if comparison < 0 {
			low = i + 1
		} else if comparison > 0 {
			high = i - 1
		} else {
			high = i
			break
		}
	}
	i = high
	if arcLengths[i] == targetArcLength {
		return float64(i) / float64(il-1)
	}
	lengthBefore := arcLengths[i]
	lengthAfter := arcLengths[i+1]
	segmentLength := lengthAfter - lengthBefore
	segmentFraction := (targetArcLength - lengthBefore) / segmentLength
	t := (float64(i) + segmentFraction) / float64(il-1)
	return t
}

// Returns a unit vector tangent at t
// In case any sub curve does not implement its tangent derivation,
// 2 points a small delta apart will be used to find its gradient
// which seems to give a reasonable approximation
func (c *baseCurve) Tangent(t float64) *Vector3 {
	delta := 0.0001
	t1 := t - delta
	t2 := t + delta
	if t1 < 0 {
		t1 = 0
	}
	if t2 > 1 {
		t2 = 1
	}
	pt1 := c.Point(t1)
	pt2 := c.Point(t2)
	tangent := pt2.Clone().Sub(*pt1).Normalize()
	return tangent
}

func (c *baseCurve) TangentAt(u float64) *Vector3 {
	t := c.uToTmapping(u)
	return c.Tangent(t)
}

func (c *baseCurve) ComputeFrenetFrames(segments int, closed bool) map[string][]Vector3 {
	normal := NewZeroVector3()
	tangents := make([]Vector3, segments+1)
	normals := make([]Vector3, segments+1)
	binormals := make([]Vector3, segments+1)
	mat := NewMatrix4Identity()
	i := 0
	for i <= segments {
		u := float64(i) / float64(segments)
		tangents[i] = *c.TangentAt(u)
		i++
	}
	normals[0] = *NewZeroVector3()
	binormals[0] = *NewZeroVector3()
	min := 999999.0
	tx := math.Abs(tangents[0].X)
	ty := math.Abs(tangents[0].Y)
	tz := math.Abs(tangents[0].Z)
	if tx <= min {
		min = tx
		normal.Set(1, 0, 0)
	}
	if ty <= min {
		min = ty
		normal.Set(0, 1, 0)
	}
	if tz <= min {
		normal.Set(0, 0, 1)
	}
	vec := CrossVector3s(tangents[0], *normal).Normalize()
	normals[0] = *CrossVector3s(tangents[0], *vec)
	binormals[0] = *CrossVector3s(tangents[0], normals[0])
	for i := 1; i <= segments; i++ {
		normals[i] = normals[i-1]
		binormals[i] = binormals[i-1]
		vec = CrossVector3s(tangents[i-1], tangents[i])
		if vec.Length() > 0.0001 {
			vec.Normalize()
			theta := math.Acos(math.Max(-1, math.Min(1, tangents[i-1].Dot(tangents[i]))))
			normals[i].ApplyMatrix4(*mat.MakeRotationAxis(*vec, theta))
		}
		binormals[i] = *CrossVector3s(tangents[i], normals[i])
	}
	if closed {
		theta := math.Acos(math.Max(-1, math.Min(1, normals[0].Dot(normals[segments]))))
		theta /= float64(segments)
		if tangents[0].Dot(*CrossVector3s(normals[0], normals[segments])) > 0 {
			theta = -theta
		}
		for i := 1; i <= segments; i++ {
			normals[i].ApplyMatrix4(*mat.MakeRotationAxis(tangents[i], theta*float64(i)))
			binormals[i] = *CrossVector3s(tangents[i], normals[i])
		}

	}
	return map[string][]Vector3{
		"tangents":  tangents,
		"normals":   normals,
		"binormals": binormals,
	}
}

func (c *baseCurve) Clone() *baseCurve {
	return &baseCurve{
		ArcLengthDivisions: c.ArcLengthDivisions,
	}
}

func (c *baseCurve) Copy(source baseCurve) *baseCurve {
	c.ArcLengthDivisions = source.ArcLengthDivisions
	return c
}
