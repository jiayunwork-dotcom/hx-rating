package heatx

import "math"

type CapacityPair struct {
	Hot  float64
	Cold float64
	Min  float64
	Max  float64
	Cr   float64
}

func CapacityOf(s Stream) float64 {
	return s.Capacity()
}

func PairCapacities(hot, cold Stream) CapacityPair {
	ch := CapacityOf(hot)
	cc := CapacityOf(cold)
	min := math.Min(ch, cc)
	max := math.Max(ch, cc)
	cr := 0.0
	if max > 0 {
		cr = min / max
	}
	return CapacityPair{Hot: ch, Cold: cc, Min: min, Max: max, Cr: cr}
}

func MaxHeatTransfer(s Spec, p CapacityPair) float64 {
	return p.Min * (s.Hot.Inlet - s.Cold.Inlet)
}

func HeatReleased(s Stream, outlet float64) float64 {
	return signedHeat(s, outlet, -1)
}

func HeatGained(s Stream, outlet float64) float64 {
	return signedHeat(s, outlet, 1)
}

func InletDifference(s Spec) float64 {
	return s.Hot.Inlet - s.Cold.Inlet
}

func OutletDifferenceParallel(s Spec, hotOut, coldOut float64) float64 {
	return hotOut - coldOut
}

func OutletDifferenceCounter(s Spec, hotOut, coldOut float64) float64 {
	d1 := s.Hot.Inlet - coldOut
	d2 := hotOut - s.Cold.Inlet
	if d1 <= 0 || d2 <= 0 {
		return 0
	}
	if math.Abs(d1-d2) < 1e-12 {
		return d1
	}
	return (d1 - d2) / math.Log(d1/d2)
}
