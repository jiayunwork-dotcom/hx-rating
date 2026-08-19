package heatx

import "math"

type ProfilePoint struct {
	Position float64
	Hot      float64
	Cold     float64
}

func TemperatureProfile(s Spec, p CapacityPair, n int) []ProfilePoint {
	if n < 2 {
		n = 2
	}
	o, err := Rate(s)
	if err != nil {
		return nil
	}
	points := make([]ProfilePoint, n)
	for i := 0; i < n; i++ {
		x := float64(i) / float64(n-1)
		th, tc := interpolate(s, o, p, x)
		points[i] = ProfilePoint{Position: x, Hot: th, Cold: tc}
	}
	return points
}

func interpolate(s Spec, o RateOutcome, p CapacityPair, x float64) (float64, float64) {
	if x <= 0 {
		return s.Hot.Inlet, s.Cold.Inlet
	}
	if x >= 1 {
		return o.HotOut, o.ColdOut
	}
	if s.Flow == Parallel {
		fraction := x * o.QNtu / (p.Min * InletDifference(s) * o.Eff)
		if o.Eff <= 0 {
			fraction = x
		}
		qx := o.QNtu * clamp01(fraction)
		return s.Hot.Inlet - qx/p.Hot, s.Cold.Inlet + qx/p.Cold
	}
	qx := o.QNtu * x
	return s.Hot.Inlet - qx/p.Hot, s.Cold.Inlet + qx/p.Cold
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func PinchPoint(s Spec, p CapacityPair) (float64, float64) {
	if _, err := Rate(s); err != nil {
		return 0, 0
	}
	profile := TemperatureProfile(s, p, 101)
	min := math.Inf(1)
	var at float64
	for _, pt := range profile {
		d := pt.Hot - pt.Cold
		if d < min {
			min = d
			at = pt.Position
		}
	}
	return min, at
}
