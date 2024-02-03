package mathx

/**
 * Bezier Curves formulas obtained from
 * https://en.wikipedia.org/wiki/B%C3%A9zier_curve
 */

func CatmullRom(t, p0, p1, p2, p3 float64) float64 {
	v0, v1 := (p2-p0)*0.5, (p3-p1)*0.5
	t2, t3 := t*t, t*t*t
	return (2*p1-2*p2+v0+v1)*t3 + (-3*p1+3*p2-2*v0-v1)*t2 + v0*t + p1
}

func QuadraticBezierP0(t, p float64) float64 {
	k := 1 - t
	return k * k * p
}

func QuadraticBezierP1(t, p float64) float64 {
	return 2 * (1 - t) * t * p
}

func QuadraticBezierP2(t, p float64) float64 {
	return t * t * p
}

func QuadraticBezier(t, p0, p1, p2 float64) float64 {
	return QuadraticBezierP0(t, p0) +
		QuadraticBezierP1(t, p1) +
		QuadraticBezierP2(t, p2)
}

func CubicBezierP0(t, p float64) float64 {
	k := 1 - t
	return k * k * k * p
}

func CubicBezierP1(t, p float64) float64 {
	k := 1 - t
	return 3 * k * k * t * p
}

func CubicBezierP2(t, p float64) float64 {
	return 3 * (1 - t) * t * t * p
}

func CubicBezierP3(t, p float64) float64 {
	return t * t * t * p
}

func CubicBezier(t, p0, p1, p2, p3 float64) float64 {
	return CubicBezierP0(t, p0) +
		CubicBezierP1(t, p1) +
		CubicBezierP2(t, p2) +
		CubicBezierP3(t, p3)
}
