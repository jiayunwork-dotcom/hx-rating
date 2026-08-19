package heatx

import "math"

type StreamResult struct {
	CapacityRate float64
	Heat         float64
	Inlet        float64
	Outlet       float64
	TemperatureChange float64
}

func StreamRates(s Spec, o RateOutcome) (hot, cold StreamResult) {
	return StreamResult{
			CapacityRate: o.Pair.Hot,
			Heat:         o.QNtu,
			Inlet:        s.Hot.Inlet,
			Outlet:       o.HotOut,
			TemperatureChange: s.Hot.Inlet - o.HotOut,
		}, StreamResult{
			CapacityRate: o.Pair.Cold,
			Heat:         o.QNtu,
			Inlet:        s.Cold.Inlet,
			Outlet:       o.ColdOut,
			TemperatureChange: o.ColdOut - s.Cold.Inlet,
		}
}

func CapacityRates(s Spec) (hotC, coldC float64) {
	return CapacityOf(s.Hot), CapacityOf(s.Cold)
}

func TemperatureEffectiveness(s Spec, o RateOutcome) float64 {
	d := InletDifference(s)
	if d <= 0 {
		return 0
	}
	return (s.Hot.Inlet - o.HotOut) / d
}

func HeatLoadFromHot(s Stream, hotOut float64) float64 {
	return HeatReleased(s, hotOut)
}

func HeatLoadFromCold(s Stream, coldOut float64) float64 {
	return HeatGained(s, coldOut)
}

func ApproachDifference(s Spec, o RateOutcome) float64 {
	return math.Abs(o.HotOut - o.ColdOut)
}
