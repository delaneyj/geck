package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Sphere[T constraints.Float] struct {
	Center Vector3[T]
	Radius T
}

func NewSphere[T constraints.Float](center Vector3[T], radius T) *Sphere[T] {
	return &Sphere[T]{
		Center: center,
		Radius: radius,
	}
}

func (s *Sphere[T]) Set(center Vector3[T], radius T) *Sphere[T] {
	s.Center = center
	s.Radius = radius
	return s
}

func (s *Sphere[T]) SetFromPoints(points []Vector3[T], optionalCenter *Vector3[T]) *Sphere[T] {
	center := s.Center
	if optionalCenter != nil {
		center = *optionalCenter
	} else {
		box := (&Box3[T]{}).SetFromPoints(points...)
		center = box.Center()
	}
	var maxRadiusSq T
	for _, p := range points {
		maxRadiusSq = max(maxRadiusSq, center.DistanceToSquared(p))
	}

	s.Radius = T(math.Sqrt(float64(maxRadiusSq)))
	return s
}

func (s *Sphere[T]) Copy(sphere Sphere[T]) *Sphere[T] {
	s.Center = sphere.Center
	s.Radius = sphere.Radius
	return s
}

func (s *Sphere[T]) IsEmpty() bool {
	return s.Radius < 0
}

func (s *Sphere[T]) MakeEmpty() *Sphere[T] {
	s.Center.Set(0, 0, 0)
	s.Radius = -1
	return s
}

func (s *Sphere[T]) ContainsPoint(point Vector3[T]) bool {
	return point.DistanceToSquared(s.Center) <= (s.Radius * s.Radius)
}

func (s *Sphere[T]) DistanceToPoint(point Vector3[T]) T {
	return point.DistanceTo(s.Center) - s.Radius
}

func (s *Sphere[T]) IntersectsSphere(sphere Sphere[T]) bool {
	radiusSum := s.Radius + sphere.Radius
	return sphere.Center.DistanceToSquared(s.Center) <= (radiusSum * radiusSum)
}

func (s *Sphere[T]) IntersectsBox(box Box3[T]) bool {
	return box.IntersectsSphere(*s)
}

func (s *Sphere[T]) IntersectsPlane(plane Plane[T]) bool {
	return T(math.Abs(float64(plane.DistanceToPoint(s.Center)))) <= s.Radius
}

func (s *Sphere[T]) ClampPoint(point Vector3[T]) *Vector3[T] {
	deltaLengthSq := s.Center.DistanceToSquared(point)
	target := point.Clone()
	if deltaLengthSq > (s.Radius * s.Radius) {
		target.Sub(s.Center).Normalize().MultiplyScalar(s.Radius).Add(s.Center)
	}

	return target
}

func (s *Sphere[T]) BoundingBox() *Box3[T] {
	target := &Box3[T]{}
	if s.IsEmpty() {
		target.MakeEmpty()
		return target
	}
	target.Set(s.Center, s.Center)
	target.ExpandByScalar(s.Radius)
	return target
}

func (s *Sphere[T]) ApplyMatrix4(matrix Matrix4[T]) *Sphere[T] {
	s.Center.ApplyMatrix4(matrix)
	s.Radius = s.Radius * matrix.MaxScaleOnAxis()
	return s
}

func (s *Sphere[T]) Translate(offset Vector3[T]) *Sphere[T] {
	s.Center.Add(offset)
	return s
}

func (s *Sphere[T]) ExpandByPoint(point Vector3[T]) *Sphere[T] {
	if s.IsEmpty() {
		s.Center = point
		s.Radius = 0
		return s
	}
	v1 := point.Clone().Sub(s.Center)
	lengthSq := v1.LengthSq()
	if lengthSq > (s.Radius * s.Radius) {
		length := T(math.Sqrt(float64(lengthSq)))
		delta := (length - s.Radius) * 0.5
		v1.AddScaledVector(s.Center, delta/length)
		s.Radius += delta
	}
	return s
}

func (s *Sphere[T]) Union(sphere Sphere[T]) *Sphere[T] {
	if sphere.IsEmpty() {
		return s
	}
	if s.IsEmpty() {
		s.Copy(sphere)
		return s
	}
	if s.Center.Equals(sphere.Center) {
		s.Radius = max(s.Radius, sphere.Radius)
	} else {
		v2 := sphere.Center.Clone().Sub(s.Center).SetLength(sphere.Radius)
		s.ExpandByPoint(*sphere.Center.Clone().Add(*v2))
		s.ExpandByPoint(*sphere.Center.Clone().Sub(*v2))
	}
	return s
}

func (s *Sphere[T]) Equals(sphere Sphere[T]) bool {
	return s.Center.Equals(sphere.Center) && (s.Radius == sphere.Radius)
}

func (s *Sphere[T]) Clone() *Sphere[T] {
	return NewSphere(s.Center, s.Radius)
}
