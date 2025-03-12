package mathx

import "golang.org/x/exp/constraints"

/**
 * Bezier Curves formulas obtained from
 * https://en.wikipedia.org/wiki/B%C3%A9zier_curve
 */

func CatmullRom[T constraints.Float](t, p0, p1, p2, p3 T) T {
	v0, v1 := (p2-p0)*0.5, (p3-p1)*0.5
	t2, t3 := t*t, t*t*t
	return (2*p1-2*p2+v0+v1)*t3 + (-3*p1+3*p2-2*v0-v1)*t2 + v0*t + p1
}

func QuadraticBezierP0[T constraints.Float](t, p T) T {
	k := 1 - t
	return k * k * p
}

func QuadraticBezierP1[T constraints.Float](t, p T) T {
	return 2 * (1 - t) * t * p
}

func QuadraticBezierP2[T constraints.Float](t, p T) T {
	return t * t * p
}

func QuadraticBezier[T constraints.Float](t, p0, p1, p2 T) T {
	return QuadraticBezierP0(t, p0) +
		QuadraticBezierP1(t, p1) +
		QuadraticBezierP2(t, p2)
}

func CubicBezierP0[T constraints.Float](t, p T) T {
	k := 1 - t
	return k * k * k * p
}

func CubicBezierP1[T constraints.Float](t, p T) T {
	k := 1 - t
	return 3 * k * k * t * p
}

func CubicBezierP2[T constraints.Float](t, p T) T {
	return 3 * (1 - t) * t * t * p
}

func CubicBezierP3[T constraints.Float](t, p T) T {
	return t * t * t * p
}

func CubicBezier[T constraints.Float](t, p0, p1, p2, p3 T) T {
	return CubicBezierP0(t, p0) +
		CubicBezierP1(t, p1) +
		CubicBezierP2(t, p2) +
		CubicBezierP3(t, p3)
}
