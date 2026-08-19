package heatx

import "math"

type ThermalCheck struct {
	HotOutAboveColdIn bool
	ColdOutBelowHotIn bool
	PinchAtHotEnd     bool
	PinchAtColdEnd    bool
}

func CheckThermal(s Spec, hotOut, coldOut float64) ThermalCheck {
	return ThermalCheck{
		HotOutAboveColdIn: hotOut > s.Cold.Inlet,
		ColdOutBelowHotIn: coldOut < s.Hot.Inlet,
		PinchAtHotEnd:     math.Abs(hotOut-s.Cold.Inlet) < 1e-9,
		PinchAtColdEnd:    math.Abs(coldOut-s.Hot.Inlet) < 1e-9,
	}
}

func ApproachHot(s Spec, hotOut float64) float64 {
	return hotOut - s.Cold.Inlet
}

func ApproachCold(s Spec, coldOut float64) float64 {
	return s.Hot.Inlet - coldOut
}

func MinimumApproach(s Spec, hotOut, coldOut float64) float64 {
	return math.Min(ApproachHot(s, hotOut), ApproachCold(s, coldOut))
}

func ApproachViolated(s Spec, hotOut, coldOut, minApproach float64) bool {
	return MinimumApproach(s, hotOut, coldOut) < minApproach-1e-12
}

func TemperatureCrossFlag(s Spec, hotOut, coldOut float64) bool {
	if s.Flow == Parallel {
		return ParallelOutletsCrossed(hotOut, coldOut)
	}
	return coldOut > hotOut
}
