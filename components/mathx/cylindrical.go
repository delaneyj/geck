package mathx

import "math"

type Cylindrical struct {
	Radius float64 // distance from the origin to a point in the x-z plane
	Theta  float64 // counterclockwise angle in the x-z plane measured in radians from the positive z-axis
	Y      float64 // height above the x-z plane
}

func NewCylindrical(radius, theta, y float64) *Cylindrical {
	return &Cylindrical{Radius: radius, Theta: theta, Y: y}
}

func (c *Cylindrical) Set(radius, theta, y float64) *Cylindrical {
	c.Radius = radius
	c.Theta = theta
	c.Y = y
	return c
}

func (c *Cylindrical) Copy(other *Cylindrical) *Cylindrical {
	c.Radius = other.Radius
	c.Theta = other.Theta
	c.Y = other.Y
	return c
}

func (c *Cylindrical) SetFromVector3(v *Vector3) *Cylindrical {
	return c.SetFromCartesianCoords(v.X, v.Y, v.Z)
}

func (c *Cylindrical) SetFromCartesianCoords(x, y, z float64) *Cylindrical {
	c.Radius = math.Sqrt(x*x + z*z)
	c.Theta = math.Atan2(x, z)
	c.Y = y
	return c
}

func (c *Cylindrical) Clone() *Cylindrical {
	return NewCylindrical(c.Radius, c.Theta, c.Y)
}
