package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Cylindrical[T constraints.Float] struct {
	Radius T // distance from the origin to a point in the x-z plane
	Theta  T // counterclockwise angle in the x-z plane measured in radians from the positive z-axis
	Y      T // height above the x-z plane
}

func NewCylindrical[T constraints.Float](radius, theta, y T) *Cylindrical[T] {
	return &Cylindrical[T]{Radius: radius, Theta: theta, Y: y}
}

func (c *Cylindrical[T]) Set(radius, theta, y T) *Cylindrical[T] {
	c.Radius = radius
	c.Theta = theta
	c.Y = y
	return c
}

func (c *Cylindrical[T]) Copy(other *Cylindrical[T]) *Cylindrical[T] {
	c.Radius = other.Radius
	c.Theta = other.Theta
	c.Y = other.Y
	return c
}

func (c *Cylindrical[T]) SetFromVector3(v *Vector3[T]) *Cylindrical[T] {
	return c.SetFromCartesianCoords(v.X, v.Y, v.Z)
}

func (c *Cylindrical[T]) SetFromCartesianCoords(x, y, z T) *Cylindrical[T] {
	xf, zf := float64(x), float64(z)
	c.Radius = T(math.Sqrt(xf*xf + zf*zf))
	c.Theta = T(math.Atan2(xf, zf))
	c.Y = y
	return c
}

func (c *Cylindrical[T]) Clone() *Cylindrical[T] {
	return NewCylindrical(c.Radius, c.Theta, c.Y)
}
