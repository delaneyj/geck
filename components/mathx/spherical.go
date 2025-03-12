package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

/**
 * Ref: https://en.wikipedia.org/wiki/Spherical_coordinate_system
 *
 * The polar angle (phi) is measured from the positive y-axis. The positive y-axis is up.
 * The azimuthal angle (theta) is measured from the positive z-axis.
 */

type Spherical[T constraints.Float] struct {
	Radius, Phi, Theta T
}

func NewSpherical[T constraints.Float](radius, phi, theta T) *Spherical[T] {
	return &Spherical[T]{Radius: radius, Phi: phi, Theta: theta}
}

func (s *Spherical[T]) Set(radius, phi, theta T) *Spherical[T] {
	s.Radius = radius
	s.Phi = phi
	s.Theta = theta
	return s
}

func (s *Spherical[T]) Copy(other *Spherical[T]) *Spherical[T] {
	s.Radius = other.Radius
	s.Phi = other.Phi
	s.Theta = other.Theta
	return s
}

func (s *Spherical[T]) MakeSafe() *Spherical[T] {
	s.Phi = max(EPSILON, min(math.Pi-EPSILON, s.Phi))
	return s
}

func (s *Spherical[T]) SetFromVector3(v *Vector3[T]) *Spherical[T] {
	return s.SetFromCartesianCoords(v.X, v.Y, v.Z)
}

func (s *Spherical[T]) SetFromCartesianCoords(x, y, z T) *Spherical[T] {
	s.Radius = T(math.Sqrt(float64(x*x + y*y + z*z)))
	if s.Radius == 0 {
		s.Theta = 0
		s.Phi = 0
	} else {
		s.Theta = T(math.Atan2(float64(x), float64(z)))
		s.Phi = T(math.Acos(float64(Clamp(y/s.Radius, -1, 1))))
	}
	return s
}

func (s *Spherical[T]) Clone() *Spherical[T] {
	return NewSpherical(s.Radius, s.Phi, s.Theta)
}
