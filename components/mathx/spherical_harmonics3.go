package mathx

import "golang.org/x/exp/constraints"

/**
 * Primary reference:
 *   https://graphics.stanford.edu/papers/envmap/envmap.pdf
 *
 * Secondary reference:
 *   https://www.ppsloan.org/publications/StupidSH36.pdf
 */

// 3-band SH defined by 9 coefficients

type SphericalHarmonics3[T constraints.Float] [9]Vector3[T]

func NewSphericalHarmonics3[T constraints.Float]() *SphericalHarmonics3[T] {
	return &SphericalHarmonics3[T]{}
}

func (sh *SphericalHarmonics3[T]) Set(coefficients ...Vector3[T]) *SphericalHarmonics3[T] {
	for i := 0; i < 9; i++ {
		sh[i].Copy(coefficients[i])
	}
	return sh
}

func (sh *SphericalHarmonics3[T]) Zero() *SphericalHarmonics3[T] {
	for i := 0; i < 9; i++ {
		sh[i].Set(0, 0, 0)
	}
	return sh
}

// GetAt returns the radiance in the direction of the normal
func (sh *SphericalHarmonics3[T]) GetAt(normal *Vector3[T]) *Vector3[T] {
	// normal is assumed to be unit length

	x, y, z := normal.X, normal.Y, normal.Z

	// band 0
	target := sh[0].Clone().MultiplyScalar(0.282095)

	// band 1
	target.AddScaledVector(sh[1], 0.488603*y)
	target.AddScaledVector(sh[2], 0.488603*z)
	target.AddScaledVector(sh[3], 0.488603*x)

	// band 2
	target.AddScaledVector(sh[4], 1.092548*(x*y))
	target.AddScaledVector(sh[5], 1.092548*(y*z))
	target.AddScaledVector(sh[6], 0.315392*(3.0*z*z-1.0))
	target.AddScaledVector(sh[7], 1.092548*(x*z))
	target.AddScaledVector(sh[8], 0.546274*(x*x-y*y))

	return target
}

// GetIrradianceAt returns the irradiance (radiance convolved with cosine lobe) in the direction of the normal
func (sh *SphericalHarmonics3[T]) GetIrradianceAt(normal Vector3[T]) *Vector3[T] {
	// normal is assumed to be unit length

	x, y, z := normal.X, normal.Y, normal.Z

	target := sh[0].Clone()

	// band 0
	target.MultiplyScalar(0.886227) // π * 0.282095

	// band 1
	target.AddScaledVector(sh[1], 2.0*0.511664*y) // ( 2 * π / 3 ) * 0.488603
	target.AddScaledVector(sh[2], 2.0*0.511664*z)
	target.AddScaledVector(sh[3], 2.0*0.511664*x)

	// band 2
	target.AddScaledVector(sh[4], 2.0*0.429043*x*y) // ( π / 4 ) * 1.092548
	target.AddScaledVector(sh[5], 2.0*0.429043*y*z)
	target.AddScaledVector(sh[6], 0.743125*z*z-0.247708) // ( π / 4 ) * 0.315392 * 3
	target.AddScaledVector(sh[7], 2.0*0.429043*x*z)
	target.AddScaledVector(sh[8], 0.429043*(x*x-y*y)) // ( π / 4 ) * 0.546274

	return target
}

func (sh *SphericalHarmonics3[T]) Add(sh2 SphericalHarmonics3[T]) *SphericalHarmonics3[T] {
	for i := 0; i < 9; i++ {
		sh[i].Add(sh2[i])
	}
	return sh
}

func (sh *SphericalHarmonics3[T]) AddScaledSH(sh2 SphericalHarmonics3[T], s T) *SphericalHarmonics3[T] {
	for i := 0; i < 9; i++ {
		sh[i].AddScaledVector(sh2[i], s)
	}
	return sh
}

func (sh *SphericalHarmonics3[T]) Scale(s T) *SphericalHarmonics3[T] {
	for i := 0; i < 9; i++ {
		sh[i].MultiplyScalar(s)
	}
	return sh
}

func (sh *SphericalHarmonics3[T]) Lerp(sh2 SphericalHarmonics3[T], alpha T) *SphericalHarmonics3[T] {
	for i := 0; i < 9; i++ {
		sh[i].Lerp(sh2[i], alpha)
	}
	return sh
}

func (sh *SphericalHarmonics3[T]) Equals(sh2 SphericalHarmonics3[T]) bool {
	for i := 0; i < 9; i++ {
		if !sh[i].Equals(sh2[i]) {
			return false
		}
	}
	return true
}

func (sh *SphericalHarmonics3[T]) Copy(sh2 SphericalHarmonics3[T]) *SphericalHarmonics3[T] {
	copy(sh[:], sh2[:])
	return sh
}

func (sh *SphericalHarmonics3[T]) Clone() *SphericalHarmonics3[T] {
	return NewSphericalHarmonics3[T]().Copy(*sh)
}

func (sh *SphericalHarmonics3[T]) FromArray(array []T, offset int) *SphericalHarmonics3[T] {
	coefficients := sh

	for i := 0; i < 9; i++ {
		coefficients[i].FromArray(array, offset+(i*3))
	}

	return sh
}

func (sh *SphericalHarmonics3[T]) ToArray(array []T, offset int) []T {
	coefficients := sh

	for i := 0; i < 9; i++ {
		coefficients[i].ToArray(array, offset+(i*3))
	}

	return array
}

// GetBasisAt evaluates the basis functions
func (sh *SphericalHarmonics3[T]) BasisAt(normal Vector3[T], shBasis []T) []T {
	// normal is assumed to be unit length
	x, y, z := normal.X, normal.Y, normal.Z

	// band 0
	shBasis[0] = 0.282095

	// band 1
	shBasis[1] = 0.488603 * y
	shBasis[2] = 0.488603 * z
	shBasis[3] = 0.488603 * x

	// band 2
	shBasis[4] = 1.092548 * x * y
	shBasis[5] = 1.092548 * y * z
	shBasis[6] = 0.315392 * (3*z*z - 1)
	shBasis[7] = 1.092548 * x * z
	shBasis[8] = 0.546274 * (x*x - y*y)

	return shBasis
}
