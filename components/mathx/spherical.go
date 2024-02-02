package mathx

import "math"

/**
 * Ref: https://en.wikipedia.org/wiki/Spherical_coordinate_system
 *
 * The polar angle (phi) is measured from the positive y-axis. The positive y-axis is up.
 * The azimuthal angle (theta) is measured from the positive z-axis.
 */

type Spherical struct {
	Radius, Phi, Theta float64
}

func NewSpherical(radius, phi, theta float64) *Spherical {
	return &Spherical{Radius: radius, Phi: phi, Theta: theta}
}

func (s *Spherical) Set(radius, phi, theta float64) *Spherical {
	s.Radius = radius
	s.Phi = phi
	s.Theta = theta
	return s
}

func (s *Spherical) Copy(other *Spherical) *Spherical {
	s.Radius = other.Radius
	s.Phi = other.Phi
	s.Theta = other.Theta
	return s
}

func (s *Spherical) MakeSafe() *Spherical {
	s.Phi = math.Max(EPSILON64, math.Min(math.Pi-EPSILON64, s.Phi))
	return s
}

func (s *Spherical) SetFromVector3(v *Vector3) *Spherical {
	return s.SetFromCartesianCoords(v.X, v.Y, v.Z)
}

func (s *Spherical) SetFromCartesianCoords(x, y, z float64) *Spherical {
	s.Radius = math.Sqrt(x*x + y*y + z*z)
	if s.Radius == 0 {
		s.Theta = 0
		s.Phi = 0
	} else {
		s.Theta = math.Atan2(x, z)
		s.Phi = math.Acos(Clamp(y/s.Radius, -1, 1))
	}
	return s
}

func (s *Spherical) Clone() *Spherical {
	return NewSpherical(s.Radius, s.Phi, s.Theta)
}
