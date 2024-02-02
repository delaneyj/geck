package mathx

import "math"

type Sphere struct {
	Center Vector3
	Radius float64
}

func NewSphere(center Vector3, radius float64) *Sphere {
	return &Sphere{
		Center: center,
		Radius: radius,
	}
}

func (s *Sphere) Set(center Vector3, radius float64) *Sphere {
	s.Center = center
	s.Radius = radius
	return s
}

func (s *Sphere) SetFromPoints(points []Vector3, optionalCenter *Vector3) *Sphere {
	center := s.Center
	if optionalCenter != nil {
		center = *optionalCenter
	} else {
		box := (&Box3{}).SetFromPoints(points...)
		center = box.Center()
	}
	maxRadiusSq := 0.0
	for _, p := range points {
		maxRadiusSq = math.Max(maxRadiusSq, center.DistanceToSquared(p))
	}

	s.Radius = math.Sqrt(maxRadiusSq)
	return s
}

func (s *Sphere) Copy(sphere Sphere) *Sphere {
	s.Center = sphere.Center
	s.Radius = sphere.Radius
	return s
}

func (s *Sphere) IsEmpty() bool {
	return s.Radius < 0
}

func (s *Sphere) MakeEmpty() *Sphere {
	s.Center.Set(0, 0, 0)
	s.Radius = -1
	return s
}

func (s *Sphere) ContainsPoint(point Vector3) bool {
	return point.DistanceToSquared(s.Center) <= (s.Radius * s.Radius)
}

func (s *Sphere) DistanceToPoint(point Vector3) float64 {
	return point.DistanceTo(s.Center) - s.Radius
}

func (s *Sphere) IntersectsSphere(sphere Sphere) bool {
	radiusSum := s.Radius + sphere.Radius
	return sphere.Center.DistanceToSquared(s.Center) <= (radiusSum * radiusSum)
}

func (s *Sphere) IntersectsBox(box Box3) bool {
	return box.IntersectsSphere(*s)
}

func (s *Sphere) IntersectsPlane(plane Plane) bool {
	return math.Abs(plane.DistanceToPoint(s.Center)) <= s.Radius
}

func (s *Sphere) ClampPoint(point Vector3) *Vector3 {
	deltaLengthSq := s.Center.DistanceToSquared(point)
	target := point.Clone()
	if deltaLengthSq > (s.Radius * s.Radius) {
		target.Sub(s.Center).Normalize().MultiplyScalar(s.Radius).Add(s.Center)
	}

	return target
}

func (s *Sphere) BoundingBox() *Box3 {
	target := &Box3{}
	if s.IsEmpty() {
		target.MakeEmpty()
		return target
	}
	target.Set(s.Center, s.Center)
	target.ExpandByScalar(s.Radius)
	return target
}

func (s *Sphere) ApplyMatrix4(matrix Matrix4) *Sphere {
	s.Center.ApplyMatrix4(matrix)
	s.Radius = s.Radius * matrix.MaxScaleOnAxis()
	return s
}

func (s *Sphere) Translate(offset Vector3) *Sphere {
	s.Center.Add(offset)
	return s
}

func (s *Sphere) ExpandByPoint(point Vector3) *Sphere {
	if s.IsEmpty() {
		s.Center = point
		s.Radius = 0
		return s
	}
	v1 := point.Clone().Sub(s.Center)
	lengthSq := v1.LengthSq()
	if lengthSq > (s.Radius * s.Radius) {
		length := math.Sqrt(lengthSq)
		delta := (length - s.Radius) * 0.5
		v1.AddScaledVector(s.Center, delta/length)
		s.Radius += delta
	}
	return s
}

func (s *Sphere) Union(sphere Sphere) *Sphere {
	if sphere.IsEmpty() {
		return s
	}
	if s.IsEmpty() {
		s.Copy(sphere)
		return s
	}
	if s.Center.Equals(sphere.Center) {
		s.Radius = math.Max(s.Radius, sphere.Radius)
	} else {
		v2 := sphere.Center.Clone().Sub(s.Center).SetLength(sphere.Radius)
		s.ExpandByPoint(*sphere.Center.Clone().Add(*v2))
		s.ExpandByPoint(*sphere.Center.Clone().Sub(*v2))
	}
	return s
}

func (s *Sphere) Equals(sphere Sphere) bool {
	return s.Center.Equals(sphere.Center) && (s.Radius == sphere.Radius)
}

func (s *Sphere) Clone() *Sphere {
	return NewSphere(s.Center, s.Radius)
}
