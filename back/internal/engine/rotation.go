package engine

type rotationCoefficients struct {
	cos, sin int
}

type precalculatedRotationCoefficients []*rotationCoefficients

var rotationCoefficientsList = precalculatedRotationCoefficients{
	&rotationCoefficients{cos: 1, sin: 0},
	&rotationCoefficients{cos: 0, sin: -1},
	&rotationCoefficients{cos: -1, sin: 0},
	&rotationCoefficients{cos: 0, sin: 1},
}

func coefficients(n int) *rotationCoefficients {
	return rotationCoefficientsList[n%4]
}
