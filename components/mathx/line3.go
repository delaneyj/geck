package mathx

import "golang.org/x/exp/constraints"

type Line3[T constraints.Float] struct {
	Start, End Vector3[T]
}

func NewLine3[T constraints.Float](start, end Vector3[T]) *Line3[T] {
	return &Line3[T]{
		Start: start,
		End:   end,
	}
}

func (l *Line3[T]) Set(start, end Vector3[T]) *Line3[T] {
	l.Start = start
	l.End = end
	return l
}

func (l *Line3[T]) Copy(line Line3[T]) *Line3[T] {
	l.Start = line.Start
	l.End = line.End
	return l
}

func (l *Line3[T]) Center() *Vector3[T] {
	return l.Start.Clone().Add(l.End).MultiplyScalar(0.5)
}

func (l *Line3[T]) Delta() *Vector3[T] {
	return l.End.Clone().Sub(l.Start)
}

func (l *Line3[T]) DistanceSq() T {
	return l.Start.DistanceToSquared(l.End)
}

func (l *Line3[T]) Distance() T {
	return l.Start.DistanceTo(l.End)
}

func (l *Line3[T]) At(t T) *Vector3[T] {
	return l.Delta().MultiplyScalar(t).Add(l.Start)
}

func (l *Line3[T]) ClosestPointToPointParameter(point Vector3[T], clampToLine bool) T {
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

func (l *Line3[T]) ClosestPointToPoint(point Vector3[T], clampToLine bool) *Vector3[T] {
	t := l.ClosestPointToPointParameter(point, clampToLine)
	return l.Delta().MultiplyScalar(t).Add(l.Start)
}

func (l *Line3[T]) ApplyMatrix4(matrix Matrix4[T]) *Line3[T] {
	l.Start.ApplyMatrix4(matrix)
	l.End.ApplyMatrix4(matrix)
	return l
}

func (l *Line3[T]) Equals(line Line3[T]) bool {
	return line.Start.Equals(l.Start) && line.End.Equals(l.End)
}

func (l *Line3[T]) Clone() *Line3[T] {
	return NewLine3(l.Start, l.End)
}
