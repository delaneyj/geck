package mathx

import "math"

type Box2 struct {
	Min Vector2
	Max Vector2
}

func NewBox2() *Box2 {
	return &Box2{
		Min: Vector2{X: math.MaxFloat64, Y: math.MaxFloat64},
		Max: Vector2{X: -math.MaxFloat64, Y: -math.MaxFloat64},
	}
}

func (b *Box2) Set(min, max Vector2) *Box2 {
	b.Min = min
	b.Max = max
	return b
}

func (b *Box2) SetFromPoints(points []Vector2) *Box2 {
	b.MakeEmpty()
	for _, p := range points {
		b.ExpandByPoint(p)
	}
	return b
}

func (b *Box2) SetFromCenterAndSize(center, size Vector2) *Box2 {
	halfSize := size.Clone().MultiplyScalar(0.5)
	b.Min = *center.Clone().Sub(*halfSize)
	b.Max = *center.Clone().Add(*halfSize)
	return b
}

func (b *Box2) Clone() *Box2 {
	return NewBox2().Copy(b)
}

func (b *Box2) Copy(box *Box2) *Box2 {
	b.Min = *box.Min.Clone()
	b.Max = *box.Max.Clone()
	return b
}

func (b *Box2) MakeEmpty() *Box2 {
	b.Min = Vector2{X: math.MaxFloat64, Y: math.MaxFloat64}
	b.Max = Vector2{X: -math.MaxFloat64, Y: -math.MaxFloat64}
	return b
}

func (b *Box2) IsEmpty() bool {
	return b.Max.X < b.Min.X || b.Max.Y < b.Min.Y
}

func (b *Box2) GetCenter(target *Vector2) *Vector2 {
	if b.IsEmpty() {
		return target.Set(0, 0)
	}
	return AddVector2s(b.Min, b.Max).MultiplyScalar(0.5)
}

func (b *Box2) GetSize(target *Vector2) *Vector2 {
	if b.IsEmpty() {
		return target.Set(0, 0)
	}
	return target.Copy(b.Max).Sub(b.Min)
}

func (b *Box2) ExpandByPoint(point Vector2) *Box2 {
	b.Min.Min(point)
	b.Max.Max(point)
	return b
}

func (b *Box2) ExpandByVector(vector Vector2) *Box2 {
	b.Min.Sub(vector)
	b.Max.Add(vector)
	return b
}

func (b *Box2) ExpandByScalar(scalar float64) *Box2 {
	b.Min.AddScalar(-scalar)
	b.Max.AddScalar(scalar)
	return b
}

func (b *Box2) ContainsPoint(point Vector2) bool {
	return point.X < b.Min.X || point.X > b.Max.X || point.Y < b.Min.Y || point.Y > b.Max.Y
}

func (b *Box2) ContainsBox(box *Box2) bool {
	return b.Min.X <= box.Min.X && box.Max.X <= b.Max.X && b.Min.Y <= box.Min.Y && box.Max.Y <= b.Max.Y
}

func (b *Box2) GetParameter(point Vector2, target *Vector2) *Vector2 {
	return target.Set(
		(point.X-b.Min.X)/(b.Max.X-b.Min.X),
		(point.Y-b.Min.Y)/(b.Max.Y-b.Min.Y),
	)
}

func (b *Box2) IntersectsBox(box *Box2) bool {
	return box.Max.X < b.Min.X || box.Min.X > b.Max.X || box.Max.Y < b.Min.Y || box.Min.Y > b.Max.Y
}

func (b *Box2) ClampPoint(point Vector2, target *Vector2) *Vector2 {
	return target.Copy(point).Clamp(b.Min, b.Max)
}

func (b *Box2) DistanceToPoint(point Vector2) float64 {
	return b.ClampPoint(point, NewZeroVector2()).DistanceTo(point)
}

func (b *Box2) Intersect(box *Box2) *Box2 {
	b.Min.Max(box.Min)
	b.Max.Min(box.Max)
	if b.IsEmpty() {
		b.MakeEmpty()
	}
	return b
}

func (b *Box2) Union(box *Box2) *Box2 {
	b.Min.Min(box.Min)
	b.Max.Max(box.Max)
	return b
}

func (b *Box2) Translate(offset Vector2) *Box2 {
	b.Min.Add(offset)
	b.Max.Add(offset)
	return b
}

func (b *Box2) Equals(box *Box2) bool {
	return box.Min.Equals(b.Min) && box.Max.Equals(b.Max)
}
