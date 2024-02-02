package mathx

type Plane struct {
	Normal   Vector3
	Constant float64
}

func NewPlane(normal Vector3, constant float64) *Plane {
	return &Plane{
		Normal:   normal,
		Constant: constant,
	}

}

func (p *Plane) Set(normal Vector3, constant float64) *Plane {
	p.Normal = normal
	p.Constant = constant
	return p
}

func (p *Plane) SetComponents(x, y, z, w float64) *Plane {
	p.Normal.Set(x, y, z)
	p.Constant = w
	return p
}

func (p *Plane) SetFromNormalAndCoplanarPoint(normal, point Vector3) *Plane {
	p.Normal = normal
	p.Constant = -point.Dot(p.Normal)
	return p
}

func (p *Plane) SetFromCoplanarPoints(a, b, c Vector3) *Plane {
	normal := SubVector3s(c, b).Cross(*SubVector3s(a, b)).Normalize()
	p.SetFromNormalAndCoplanarPoint(*normal, a)
	return p
}

func (p *Plane) Copy(plane Plane) *Plane {
	p.Normal = plane.Normal
	p.Constant = plane.Constant
	return p
}

func (p *Plane) Normalize() *Plane {
	inverseNormalLength := 1.0 / p.Normal.Length()
	p.Normal.MultiplyScalar(inverseNormalLength)
	p.Constant *= inverseNormalLength
	return p
}

func (p *Plane) Negate() *Plane {
	p.Constant *= -1
	p.Normal.Negate()
	return p
}

func (p *Plane) DistanceToPoint(point Vector3) float64 {
	return p.Normal.Dot(point) + p.Constant
}

func (p *Plane) DistanceToSphere(sphere Sphere) float64 {
	return p.DistanceToPoint(sphere.Center) - sphere.Radius
}

func (p *Plane) ProjectPoint(point Vector3) *Vector3 {
	return point.Clone().AddScaledVector(p.Normal, -p.DistanceToPoint(point))
}

func (p *Plane) IntersectLine(line Line3) *Vector3 {
	direction := line.Delta()

	denominator := p.Normal.Dot(*direction)

	if denominator == 0 {
		if p.DistanceToPoint(line.Start) == 0 {
			return line.Start.Clone()
		}
		return nil
	}

	t := -(line.Start.Dot(p.Normal) + p.Constant) / denominator

	if t < 0 || t > 1 {
		return nil
	}

	return line.Start.Clone().AddScaledVector(*direction, t)
}

func (p *Plane) IntersectsLine(line Line3) bool {
	startSign := p.DistanceToPoint(line.Start)
	endSign := p.DistanceToPoint(line.End)

	return (startSign < 0 && endSign > 0) || (endSign < 0 && startSign > 0)
}

func (p *Plane) IntersectsBox(box Box3) bool {
	return box.IntersectsPlane(*p)
}

func (p *Plane) IntersectsSphere(sphere Sphere) bool {
	return sphere.IntersectsPlane(*p)
}

func (p *Plane) CoplanarPoint() *Vector3 {
	return p.Normal.Clone().MultiplyScalar(-p.Constant)
}

func (p *Plane) ApplyMatrix4(matrix Matrix4, optionalNormalMatrix *Matrix3) *Plane {
	normalMatrix := optionalNormalMatrix
	if normalMatrix == nil {
		normalMatrix = normalMatrix.NormalMatrix(matrix)
	}

	referencePoint := p.CoplanarPoint().ApplyMatrix4(matrix)

	normal := p.Normal.ApplyMatrix3(*normalMatrix).Normalize()

	p.Constant = -referencePoint.Dot(*normal)

	return p
}

func (p *Plane) Translate(offset Vector3) *Plane {
	p.Constant -= offset.Dot(p.Normal)
	return p
}

func (p *Plane) Equals(plane Plane) bool {
	return plane.Normal.Equals(p.Normal) && plane.Constant == p.Constant
}

func (p *Plane) Clone() *Plane {
	return NewPlane(p.Normal, p.Constant)
}
