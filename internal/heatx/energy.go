package heatx

import "math"

type EnergyCheck struct {
	Released float64
	Gained   float64
	Residual float64
	Balanced bool
}

func EnergyBalance(s Spec, q, hotOut, coldOut float64) EnergyCheck {
	released := HeatReleased(s.Hot, hotOut)
	gained := HeatGained(s.Cold, coldOut)
	residual := released - gained
	absQ := math.Max(math.Abs(q), 1e-15)
	rel := math.Abs(residual) / absQ
	return EnergyCheck{
		Released: released,
		Gained:   gained,
		Residual: residual,
		Balanced: rel < 1e-9,
	}
}

func ToleranceError(a, b, tol float64) float64 {
	return math.Abs(a-b) / math.Max(math.Abs(b), 1e-15)
}

func CrossCheckNTULMTD(qNtu, qLmtd float64) float64 {
	return RelativeError(qNtu, qLmtd)
}

func WithinTolerance(a, b, tol float64) bool {
	return ToleranceError(a, b, tol) <= tol
}

func MaxPossibleQ(s Spec, p CapacityPair) float64 {
	d := InletDifference(s)
	if d <= 0 {
		return 0
	}
	return p.Min * d
}

func FeasibleTarget(s Spec, p CapacityPair, targetQ float64) bool {
	return targetQ <= MaxPossibleQ(s, p)+1e-12
}
