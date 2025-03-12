package mathx

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Ray[T constraints.Float] struct {
	Origin, Dir Vector3[T]
}

func NewRay[T constraints.Float](origin, dir Vector3[T]) *Ray[T] {
	return &Ray[T]{origin, dir}
}

func (r *Ray[T]) Set(origin, dir Vector3[T]) *Ray[T] {
	r.Origin = origin
	r.Dir = dir
	return r
}

func (r *Ray[T]) Copy(ray *Ray[T]) *Ray[T] {
	r.Origin = ray.Origin
	r.Dir = ray.Dir
	return r
}

func (r *Ray[T]) At(t T) *Vector3[T] {
	return r.Origin.Clone().AddScaledVector(r.Dir, t)
}

func (r *Ray[T]) LookAt(v *Vector3[T]) *Ray[T] {
	r.Dir = *v.Clone().Sub(r.Origin).Normalize()
	return r
}

func (r *Ray[T]) Recast(t T) *Ray[T] {
	r.Origin = *r.At(t)
	return r
}

func (r *Ray[T]) ClosestPointToPoint(point Vector3[T]) *Vector3[T] {
	target := SubVector3s(point, r.Origin)
	directionDistance := target.Dot(r.Dir)
	if directionDistance < 0 {
		return target.Copy(r.Origin)
	}
	return target.Copy(r.Origin).AddScaledVector(r.Dir, directionDistance)
}

func (r *Ray[T]) DistanceToPoint(point Vector3[T]) T {
	return T(math.Sqrt(float64(r.DistanceSqToPoint(point))))
}

func (r *Ray[T]) DistanceSqToPoint(point Vector3[T]) T {
	directionDistance := SubVector3s(point, r.Origin).Dot(r.Dir)
	if directionDistance < 0 {
		return r.Origin.DistanceToSquared(point)
	}
	_vector := r.Origin.Clone().AddScaledVector(r.Dir, directionDistance)
	return _vector.DistanceToSquared(point)
}

func (r *Ray[T]) DistanceSqToSegment(v0, v1 Vector3[T], optionalPointOnRay, optionalPointOnSegment *Vector3[T]) T {
	// from
	// It returns the min distance between the ray and the segment
	// defined by v0 and v1
	// It can also set two optional targets :
	// - The closest point on the ray
	// - The closest point on the segment
	_segCenter := v0.Clone().Add(v1).MultiplyScalar(0.5)
	_segDir := v1.Clone().Sub(v0).Normalize()
	_diff := r.Origin.Clone().Sub(*_segCenter)
	segExtent := float64(v0.DistanceTo(v1) * 0.5)
	a01 := float64(-r.Dir.Dot(*_segDir))
	b0 := float64(_diff.Dot(r.Dir))
	b1 := float64(-_diff.Dot(*_segDir))
	c := float64(_diff.LengthSq())
	det := math.Abs(1 - float64(a01*a01))
	var s0, s1, sqrDist, extDet float64
	if det > 0 {
		// The ray and segment are not parallel.
		s0 = a01*b1 - b0
		s1 = a01*b0 - b1
		extDet = float64(segExtent) * det
		if s0 >= 0 {
			if s1 >= -extDet {
				if s1 <= extDet {
					// region 0
					// Minimum at interior points of ray and segment.
					invDet := 1 / det
					s0 *= invDet
					s1 *= invDet
					sqrDist = s0*(s0+a01*s1+2*b0) + s1*(a01*s0+s1+2*b1) + c
				} else {
					// region 1
					s1 = segExtent
					s0 = max(0, -(a01*s1 + b0))
					sqrDist = -s0*s0 + s1*(s1+2*b1) + c
				}
			} else {
				// region 5
				s1 = -segExtent
				s0 = max(0, -(a01*s1 + b0))
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c
			}
		} else {
			if s1 <= -extDet {
				// region 4
				s0 = max(0, -(-a01*segExtent + b0))
				s1 = -segExtent
				if s0 <= 0 {
					s1 = min(max(-segExtent, -b1), segExtent)
				}
				sqrDist = -s0*s0 + s1*(s1+2*b1) + c
			} else if s1 <= extDet {
				// region 3
				s0 = 0
				s1 = min(max(-segExtent, -b1), segExtent)
				sqrDist = s1*(s1+2*b1) + c
			} else {
				// region 2
				s0 = max(0, -(a01*segExtent + b0))
				s1 = segExtent
				if s0 <= 0 {
					s1 = min(max(-segExtent, -b1), segExtent)
				}
				sqrDist = -s0 * s0

				if s1 > 0 {
					sqrDist += s1*(s1+2*b1) + c
				} else {
					sqrDist += c
				}

			}
		}
	} else {
		// Ray and segment are parallel.
		s1 = segExtent
		if a01 > 0 {
			s1 = -segExtent
		}
		s0 = max(0, -(a01*s1 + b0))
		sqrDist = -s0*s0 + s1*(s1+2*b1) + c
	}
	if optionalPointOnRay != nil {
		optionalPointOnRay.Copy(r.Origin).AddScaledVector(r.Dir, T(s0))
	}
	if optionalPointOnSegment != nil {
		optionalPointOnSegment.Copy(*_segCenter).AddScaledVector(*_segDir, T(s1))
	}
	return T(sqrDist)
}

func (r *Ray[T]) IntersectSphere(sphere *Sphere[T], target *Vector3[T]) *Vector3[T] {
	_vector := SubVector3s(sphere.Center, r.Origin)
	tca := _vector.Dot(r.Dir)
	d2 := _vector.Dot(*_vector) - tca*tca
	radius2 := sphere.Radius * sphere.Radius
	if d2 > radius2 {
		return nil
	}
	thc := T(math.Sqrt(float64(radius2 - d2)))
	t0 := tca - thc
	t1 := tca + thc
	if t1 < 0 {
		return nil
	}
	if t0 < 0 {
		return r.At(t1)
	}
	return r.At(t0)
}

func (r *Ray[T]) IntersectsSphere(sphere *Sphere[T]) bool {
	return r.DistanceSqToPoint(sphere.Center) <= sphere.Radius*sphere.Radius
}

func (r *Ray[T]) DistanceToPlane(plane *Plane[T]) T {
	denominator := plane.Normal.Dot(r.Dir)
	if denominator == 0 {
		if plane.DistanceToPoint(r.Origin) == 0 {
			return 0
		}
		return T(math.NaN())
	}
	t := -(r.Origin.Dot(plane.Normal) + plane.Constant) / denominator
	if t >= 0 {
		return t
	}
	return T(math.NaN())
}

func (r *Ray[T]) IntersectPlane(plane *Plane[T]) *Vector3[T] {
	t := r.DistanceToPlane(plane)
	if math.IsNaN(float64(t)) {
		return nil
	}
	return r.At(t)
}

func (r *Ray[T]) IntersectsPlane(plane *Plane[T]) bool {
	distToPoint := plane.DistanceToPoint(r.Origin)
	if distToPoint == 0 {
		return true
	}
	denominator := plane.Normal.Dot(r.Dir)
	return denominator*distToPoint < 0
}

func (r *Ray[T]) IntersectBox(box *Box3[T]) *Vector3[T] {
	var tmin, tmax, tymin, tymax, tzmin, tzmax T
	invdirx := 1 / r.Dir.X
	invdiry := 1 / r.Dir.Y
	invdirz := 1 / r.Dir.Z
	origin := r.Origin
	if invdirx >= 0 {
		tmin = (box.Min.X - origin.X) * invdirx
		tmax = (box.Max.X - origin.X) * invdirx
	} else {
		tmin = (box.Max.X - origin.X) * invdirx
		tmax = (box.Min.X - origin.X) * invdirx
	}
	if invdiry >= 0 {
		tymin = (box.Min.Y - origin.Y) * invdiry
		tymax = (box.Max.Y - origin.Y) * invdiry
	} else {
		tymin = (box.Max.Y - origin.Y) * invdiry
		tymax = (box.Min.Y - origin.Y) * invdiry
	}
	if tmin > tymax || tymin > tmax {
		return nil
	}
	if tymin > tmin || math.IsNaN(float64(tmin)) {
		tmin = tymin
	}
	if tymax < tmax || math.IsNaN(float64(tmax)) {
		tmax = tymax
	}
	if invdirz >= 0 {
		tzmin = (box.Min.Z - origin.Z) * invdirz
		tzmax = (box.Max.Z - origin.Z) * invdirz
	} else {
		tzmin = (box.Max.Z - origin.Z) * invdirz
		tzmax = (box.Min.Z - origin.Z) * invdirz
	}
	if tmin > tzmax || tzmin > tmax {
		return nil
	}
	if tzmin > tmin || math.IsNaN(float64(tmin)) {
		tmin = tzmin
	}
	if tzmax < tmax || math.IsNaN(float64(tmax)) {
		tmax = tzmax
	}
	if tmax < 0 {
		return nil
	}

	if tmin >= 0 {
		return r.At(tmin)
	}
	return r.At(tmax)
}

func (r *Ray[T]) IntersectsBox(box *Box3[T]) bool {
	return r.IntersectBox(box) != nil
}

func (r *Ray[T]) IntersectTriangle(a, b, c Vector3[T], backfaceCulling bool) *Vector3[T] {
	// Compute the offset origin, edges, and normal.
	// from
	// It returns the min distance between the ray and the segment
	// defined by v0 and v1
	// It can also set two optional targets :
	// - The closest point on the ray
	// - The closest point on the segment
	_edge1 := SubVector3s(b, a)
	_edge2 := SubVector3s(c, a)
	_normal := CrossVector3s(*_edge1, *_edge2)
	// Solve Q + t*D = b1*E1 + b2*E2 (Q = kDiff, D = ray direction,
	// E1 = kEdge1, E2 = kEdge2, N = Cross(E1,E2)) by
	//   |Dot(D,N)|*b1 = sign(Dot(D,N))*Dot(D,Cross(Q,E2))
	//   |Dot(D,N)|*b2 = sign(Dot(D,N))*Dot(D,Cross(E1,Q))
	//   |Dot(D,N)|*t = -sign(Dot(D,N))*Dot(Q,N)
	DdN := r.Dir.Dot(*_normal)
	var sign T
	if DdN > 0 {
		if backfaceCulling {
			return nil
		}
		sign = 1
	}
	if DdN < 0 {
		sign = -1
		DdN = -DdN
	}
	if DdN == 0 {
		return nil
	}
	_diff := SubVector3s(r.Origin, a)
	_edge2 = CrossVector3s(*_diff, *_edge2)
	DdQxE2 := sign * r.Dir.Dot(*_edge2)
	if DdQxE2 < 0 {
		return nil
	}
	DdE1xQ := sign * r.Dir.Dot(*_edge1.Cross(*_diff))
	if DdE1xQ < 0 {
		return nil
	}
	if DdQxE2+DdE1xQ > DdN {
		return nil
	}
	QdN := -sign * _diff.Dot(*_normal)
	if QdN < 0 {
		return nil
	}
	return r.At(QdN / DdN)
}

func (r *Ray[T]) ApplyMatrix4(matrix4 Matrix4[T]) *Ray[T] {
	r.Origin.ApplyMatrix4(matrix4)
	r.Dir.TransformDirection(matrix4)
	return r
}

func (r *Ray[T]) Equals(ray Ray[T]) bool {
	return ray.Origin.Equals(r.Origin) && ray.Dir.Equals(r.Dir)
}

func (r *Ray[T]) Clone() *Ray[T] {
	return NewRay(r.Origin, r.Dir)
}
