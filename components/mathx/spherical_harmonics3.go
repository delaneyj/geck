package mathx

/**
 * Primary reference:
 *   https://graphics.stanford.edu/papers/envmap/envmap.pdf
 *
 * Secondary reference:
 *   https://www.ppsloan.org/publications/StupidSH36.pdf
 */

// 3-band SH defined by 9 coefficients

type SphericalHarmonics3 [9]Vector3

func NewSphericalHarmonics3() *SphericalHarmonics3 {
	return &SphericalHarmonics3{}
}

func (sh *SphericalHarmonics3) Set(coefficients ...Vector3) *SphericalHarmonics3 {
	for i := 0; i < 9; i++ {
		sh[i].Copy(coefficients[i])
	}
	return sh
}

func (sh *SphericalHarmonics3) Zero() *SphericalHarmonics3 {
	for i := 0; i < 9; i++ {
		sh[i].Set(0, 0, 0)
	}
	return sh
}

// GetAt returns the radiance in the direction of the normal
func (sh *SphericalHarmonics3) GetAt(normal *Vector3) *Vector3 {
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
func (sh *SphericalHarmonics3) GetIrradianceAt(normal Vector3) *Vector3 {
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

func (sh *SphericalHarmonics3) Add(sh2 SphericalHarmonics3) *SphericalHarmonics3 {
	for i := 0; i < 9; i++ {
		sh[i].Add(sh2[i])
	}
	return sh
}

func (sh *SphericalHarmonics3) AddScaledSH(sh2 SphericalHarmonics3, s float64) *SphericalHarmonics3 {
	for i := 0; i < 9; i++ {
		sh[i].AddScaledVector(sh2[i], s)
	}
	return sh
}

func (sh *SphericalHarmonics3) Scale(s float64) *SphericalHarmonics3 {
	for i := 0; i < 9; i++ {
		sh[i].MultiplyScalar(s)
	}
	return sh
}

func (sh *SphericalHarmonics3) Lerp(sh2 SphericalHarmonics3, alpha float64) *SphericalHarmonics3 {
	for i := 0; i < 9; i++ {
		sh[i].Lerp(sh2[i], alpha)
	}
	return sh
}

func (sh *SphericalHarmonics3) Equals(sh2 SphericalHarmonics3) bool {
	for i := 0; i < 9; i++ {
		if !sh[i].Equals(sh2[i]) {
			return false
		}
	}
	return true
}

func (sh *SphericalHarmonics3) Copy(sh2 SphericalHarmonics3) *SphericalHarmonics3 {
	copy(sh[:], sh2[:])
	return sh
}

func (sh *SphericalHarmonics3) Clone() *SphericalHarmonics3 {
	return NewSphericalHarmonics3().Copy(*sh)
}

func (sh *SphericalHarmonics3) FromArray(array []float64, offset int) *SphericalHarmonics3 {
	coefficients := sh

	for i := 0; i < 9; i++ {
		coefficients[i].FromArray(array, offset+(i*3))
	}

	return sh
}

func (sh *SphericalHarmonics3) ToArray(array []float64, offset int) []float64 {
	coefficients := sh

	for i := 0; i < 9; i++ {
		coefficients[i].ToArray(array, offset+(i*3))
	}

	return array
}

// GetBasisAt evaluates the basis functions
func (sh *SphericalHarmonics3) BasisAt(normal Vector3, shBasis []float64) []float64 {
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
