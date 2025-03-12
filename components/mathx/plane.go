package mathx

import "golang.org/x/exp/constraints"

type Plane[T constraints.Float] struct {
	Normal   Vector3[T]
	Constant T
}

func NewPlane[T constraints.Float](normal Vector3[T], constant T) *Plane[T] {
	return &Plane[T]{
		Normal:   normal,
		Constant: constant,
	}

}

func (p *Plane[T]) Set(normal Vector3[T], constant T) *Plane[T] {
	p.Normal = normal
	p.Constant = constant
	return p
}

func (p *Plane[T]) SetComponents(x, y, z, w T) *Plane[T] {
	p.Normal.Set(x, y, z)
	p.Constant = w
	return p
}

func (p *Plane[T]) SetFromNormalAndCoplanarPoint(normal, point Vector3[T]) *Plane[T] {
	p.Normal = normal
	p.Constant = -point.Dot(p.Normal)
	return p
}

func (p *Plane[T]) SetFromCoplanarPoints(a, b, c Vector3[T]) *Plane[T] {
	normal := SubVector3s(c, b).Cross(*SubVector3s(a, b)).Normalize()
	p.SetFromNormalAndCoplanarPoint(*normal, a)
	return p
}

func (p *Plane[T]) Copy(plane Plane[T]) *Plane[T] {
	p.Normal = plane.Normal
	p.Constant = plane.Constant
	return p
}

func (p *Plane[T]) Normalize() *Plane[T] {
	inverseNormalLength := 1.0 / p.Normal.Length()
	p.Normal.MultiplyScalar(inverseNormalLength)
	p.Constant *= inverseNormalLength
	return p
}

func (p *Plane[T]) Negate() *Plane[T] {
	p.Constant *= -1
	p.Normal.Negate()
	return p
}

func (p *Plane[T]) DistanceToPoint(point Vector3[T]) T {
	return p.Normal.Dot(point) + p.Constant
}

func (p *Plane[T]) DistanceToSphere(sphere Sphere[T]) T {
	return p.DistanceToPoint(sphere.Center) - sphere.Radius
}

func (p *Plane[T]) ProjectPoint(point Vector3[T]) *Vector3[T] {
	return point.Clone().AddScaledVector(p.Normal, -p.DistanceToPoint(point))
}

func (p *Plane[T]) IntersectLine(line Line3[T]) *Vector3[T] {
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

func (p *Plane[T]) IntersectsLine(line Line3[T]) bool {
	startSign := p.DistanceToPoint(line.Start)
	endSign := p.DistanceToPoint(line.End)

	return (startSign < 0 && endSign > 0) || (endSign < 0 && startSign > 0)
}

func (p *Plane[T]) IntersectsBox(box Box3[T]) bool {
	return box.IntersectsPlane(*p)
}

func (p *Plane[T]) IntersectsSphere(sphere Sphere[T]) bool {
	return sphere.IntersectsPlane(*p)
}

func (p *Plane[T]) CoplanarPoint() *Vector3[T] {
	return p.Normal.Clone().MultiplyScalar(-p.Constant)
}

func (p *Plane[T]) ApplyMatrix4(matrix Matrix4[T], optionalNormalMatrix *Matrix3[T]) *Plane[T] {
	normalMatrix := optionalNormalMatrix
	if normalMatrix == nil {
		normalMatrix = normalMatrix.NormalMatrix(matrix)
	}

	referencePoint := p.CoplanarPoint().ApplyMatrix4(matrix)

	normal := p.Normal.ApplyMatrix3(*normalMatrix).Normalize()

	p.Constant = -referencePoint.Dot(*normal)

	return p
}

func (p *Plane[T]) Translate(offset Vector3[T]) *Plane[T] {
	p.Constant -= offset.Dot(p.Normal)
	return p
}

func (p *Plane[T]) Equals(plane Plane[T]) bool {
	return plane.Normal.Equals(p.Normal) && plane.Constant == p.Constant
}

func (p *Plane[T]) Clone() *Plane[T] {
	return NewPlane(p.Normal, p.Constant)
}
