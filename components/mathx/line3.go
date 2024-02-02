package mathx

type Line3 struct {
	Start, End Vector3
}

func NewLine3(start, end Vector3) *Line3 {
	return &Line3{
		Start: start,
		End:   end,
	}
}

func (l *Line3) Set(start, end Vector3) *Line3 {
	l.Start = start
	l.End = end
	return l
}

func (l *Line3) Copy(line Line3) *Line3 {
	l.Start = line.Start
	l.End = line.End
	return l
}

func (l *Line3) Center() *Vector3 {
	return l.Start.Clone().Add(l.End).MultiplyScalar(0.5)
}

func (l *Line3) Delta() *Vector3 {
	return l.End.Clone().Sub(l.Start)
}

func (l *Line3) DistanceSq() float64 {
	return l.Start.DistanceToSquared(l.End)
}

func (l *Line3) Distance() float64 {
	return l.Start.DistanceTo(l.End)
}

func (l *Line3) At(t float64) *Vector3 {
	return l.Delta().MultiplyScalar(t).Add(l.Start)
}

func (l *Line3) ClosestPointToPointParameter(point Vector3, clampToLine bool) float64 {
	startP := point.Clone().Sub(l.Start)
	startEnd := l.End.Clone().Sub(l.Start)

	startEnd2 := startEnd.Dot(*startEnd)
	startEnd_startP := startEnd.Dot(*startP)

	t := startEnd_startP / startEnd2

	if clampToLine {
		t = Clamp(t, 0, 1)
	}

	return t
}

func (l *Line3) ClosestPointToPoint(point Vector3, clampToLine bool) *Vector3 {
	t := l.ClosestPointToPointParameter(point, clampToLine)
	return l.Delta().MultiplyScalar(t).Add(l.Start)
}

func (l *Line3) ApplyMatrix4(matrix Matrix4) *Line3 {
	l.Start.ApplyMatrix4(matrix)
	l.End.ApplyMatrix4(matrix)
	return l
}

func (l *Line3) Equals(line Line3) bool {
	return line.Start.Equals(l.Start) && line.End.Equals(l.End)
}

func (l *Line3) Clone() *Line3 {
	return NewLine3(l.Start, l.End)
}
