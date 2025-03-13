package mathx

import "golang.org/x/exp/constraints"

type ArcCurve[T constraints.Float] struct {
	EllipseCurve[T]
}

func NewArcCurve[T constraints.Float](aX, aY, aRadius, aStartAngle, aEndAngle T, aClockwise bool) *ArcCurve[T] {
	return &ArcCurve[T]{
		EllipseCurve: *NewEllipseCurve(
			aX, aY,
			aRadius, aRadius,
			aStartAngle, aEndAngle,
			0, aClockwise,
		),
	}
}
