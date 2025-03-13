package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Box2[T constraints.Float] struct {
	Min Vector2[T]
	Max Vector2[T]
}

func NewBox2[T constraints.Float]() *Box2[T] {
	var (
		maxF = float64(math.MaxFloat64)
		minF = float64(-math.MaxFloat64)
	)

	return &Box2[T]{
		Min: Vector2[T]{X: T(maxF), Y: T(maxF)},
		Max: Vector2[T]{X: T(minF), Y: T(minF)},
	}
}

func (b *Box2[T]) Set(min, max Vector2[T]) *Box2[T] {
	b.Min = min
	b.Max = max
	return b
}

func (b *Box2[T]) SetFromPoints(points []Vector2[T]) *Box2[T] {
	b.MakeEmpty()
	for _, p := range points {
		b.ExpandByPoint(p)
	}
	return b
}

func (b *Box2[T]) SetFromCenterAndSize(center, size Vector2[T]) *Box2[T] {
	halfSize := size.Clone().MultiplyScalar(0.5)
	b.Min = *center.Clone().Sub(*halfSize)
	b.Max = *center.Clone().Add(*halfSize)
	return b
}

func (b *Box2[T]) Clone() *Box2[T] {
	return NewBox2[T]().Copy(b)
}

func (b *Box2[T]) Copy(box *Box2[T]) *Box2[T] {
	b.Min = *box.Min.Clone()
	b.Max = *box.Max.Clone()
	return b
}

func (b *Box2[T]) MakeEmpty() *Box2[T] {
	var (
		maxF = float64(math.MaxFloat64)
		minF = float64(-math.MaxFloat64)
	)
	b.Min.Set(T(maxF), T(maxF))
	b.Max.Set(T(minF), T(minF))
	return b
}

func (b *Box2[T]) IsEmpty() bool {
	return b.Max.X < b.Min.X || b.Max.Y < b.Min.Y
}

func (b *Box2[T]) GetCenter(target *Vector2[T]) *Vector2[T] {
	if b.IsEmpty() {
		return target.Set(0, 0)
	}
	return AddVector2s(b.Min, b.Max).MultiplyScalar(0.5)
}

func (b *Box2[T]) GetSize(target *Vector2[T]) *Vector2[T] {
	if b.IsEmpty() {
		return target.Set(0, 0)
	}
	return target.Copy(b.Max).Sub(b.Min)
}

func (b *Box2[T]) ExpandByPoint(point Vector2[T]) *Box2[T] {
	b.Min.Min(point)
	b.Max.Max(point)
	return b
}

func (b *Box2[T]) ExpandByVector(vector Vector2[T]) *Box2[T] {
	b.Min.Sub(vector)
	b.Max.Add(vector)
	return b
}

func (b *Box2[T]) ExpandByScalar(scalar T) *Box2[T] {
	b.Min.AddScalar(-scalar)
	b.Max.AddScalar(scalar)
	return b
}

func (b *Box2[T]) ContainsPoint(point Vector2[T]) bool {
	return point.X >= b.Min.X && point.X <= b.Max.X &&
		point.Y >= b.Min.Y && point.Y <= b.Max.Y
}

func (b *Box2[T]) ContainsBox(box *Box2[T]) bool {
	return b.Min.X <= box.Min.X && box.Max.X <= b.Max.X && b.Min.Y <= box.Min.Y && box.Max.Y <= b.Max.Y
}

func (b *Box2[T]) GetParameter(point Vector2[T], target *Vector2[T]) *Vector2[T] {
	return target.Set(
		(point.X-b.Min.X)/(b.Max.X-b.Min.X),
		(point.Y-b.Min.Y)/(b.Max.Y-b.Min.Y),
	)
}

func (b *Box2[T]) IntersectsBox(box *Box2[T]) bool {
	return box.Max.X < b.Min.X || box.Min.X > b.Max.X || box.Max.Y < b.Min.Y || box.Min.Y > b.Max.Y
}

func (b *Box2[T]) ClampPoint(point Vector2[T], target *Vector2[T]) *Vector2[T] {
	return target.Copy(point).Clamp(b.Min, b.Max)
}

func (b *Box2[T]) DistanceToPoint(point Vector2[T]) T {
	return b.ClampPoint(point, NewZeroVector2[T]()).DistanceTo(point)
}

func (b *Box2[T]) Intersect(box *Box2[T]) *Box2[T] {
	b.Min.Max(box.Min)
	b.Max.Min(box.Max)
	if b.IsEmpty() {
		b.MakeEmpty()
	}
	return b
}

func (b *Box2[T]) Union(box *Box2[T]) *Box2[T] {
	b.Min.Min(box.Min)
	b.Max.Max(box.Max)
	return b
}

func (b *Box2[T]) Translate(offset Vector2[T]) *Box2[T] {
	b.Min.Add(offset)
	b.Max.Add(offset)
	return b
}

func (b *Box2[T]) Equals(box *Box2[T]) bool {
	return box.Min.Equals(b.Min) && box.Max.Equals(b.Max)
}
