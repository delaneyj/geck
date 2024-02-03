package mathx

type ArcCurve struct {
	EllipseCurve
}

func NewArcCurve(aX, aY, aRadius, aStartAngle, aEndAngle float64, aClockwise bool) *ArcCurve {
	return &ArcCurve{
		EllipseCurve: *NewEllipseCurve(
			aX, aY,
			aRadius, aRadius,
			aStartAngle, aEndAngle,
			0, aClockwise,
		),
	}
}
