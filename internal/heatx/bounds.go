package heatx

import "math"

type BoundsCheck struct {
	MaxQ       float64
	InletDiff  float64
	Feasible   bool
	Reason     string
	HeatAbove  bool
	CrossedOut bool
}

func InletSpan(s Spec) float64 {
	return InletDifference(s)
}

func CheckBounds(s Spec, p CapacityPair, q float64) BoundsCheck {
	maxQ := MaxPossibleQ(s, p)
	feasible := q <= maxQ+1e-12
	reason := ""
	if !feasible {
		reason = "unreachable: requested heat exceeds Cmin*(Th,in-Tc,in)"
	}
	return BoundsCheck{
		MaxQ:      maxQ,
		InletDiff: InletSpan(s),
		Feasible:  feasible,
		Reason:    reason,
		HeatAbove: !feasible,
	}
}

func TargetOutletsWithin(s Spec, p CapacityPair, hotOut, coldOut float64) bool {
	q := HeatReleased(s.Hot, hotOut)
	check := CheckBounds(s, p, q)
	return check.Feasible
}

func ParallelOutletsCrossed(hotOut, coldOut float64) bool {
	return coldOut >= hotOut
}

func AnyOutletCrossed(s Spec, hotOut, coldOut float64) bool {
	if s.Flow == Parallel {
		return ParallelOutletsCrossed(hotOut, coldOut)
	}
	d1 := CounterDeltaT1(s, hotOut, coldOut)
	d2 := CounterDeltaT2(s, hotOut, coldOut)
	return d1 <= 0 || d2 <= 0
}

func ThermalCrossPossible(hotOut, coldOut float64) bool {
	return hotOut <= coldOut
}

func EnsurePositive(d float64) float64 {
	return math.Max(d, 0)
}
