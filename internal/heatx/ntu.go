package heatx

import "math"

type NTUParams struct {
	NTU float64
	Cr  float64
}

func EffectivenessCounter(ntu, cr float64) float64 {
	if ntu <= 0 {
		return 0
	}
	if cr >= 1 {
		if ntu > 50 {
			return 1 - 1e-9
		}
		return ntu / (1 + ntu)
	}
	exp := math.Exp(-ntu * (1 - cr))
	return (1 - exp) / (1 - cr*exp)
}

func EffectivenessParallel(ntu, cr float64) float64 {
	if ntu <= 0 {
		return 0
	}
	return (1 - math.Exp(-ntu*(1+cr))) / (1 + cr)
}

func Effectiveness(s Spec, ntu float64, p CapacityPair) float64 {
	switch s.Flow {
	case Counter:
		return EffectivenessCounter(ntu, p.Cr)
	case Parallel:
		return EffectivenessParallel(ntu, p.Cr)
	default:
		return 0
	}
}

func NTUOf(s Spec, p CapacityPair) float64 {
	return s.UA / p.Min
}

func NTUFromEffectiveness(s Spec, eff float64, p CapacityPair) float64 {
	switch s.Flow {
	case Counter:
		return NTUCounterInverse(eff, p.Cr)
	case Parallel:
		return NTUParallelInverse(eff, p.Cr)
	default:
		return 0
	}
}

func NTUCounterInverse(eff, cr float64) float64 {
	if eff <= 0 {
		return 0
	}
	if eff >= 1 {
		return math.Inf(1)
	}
	if cr >= 1 {
		return eff / (1 - eff)
	}
	num := math.Log((1 - eff*cr) / (1 - eff))
	return num / (1 - cr)
}

func NTUParallelInverse(eff, cr float64) float64 {
	if eff <= 0 {
		return 0
	}
	if eff >= 1 {
		return math.Inf(1)
	}
	return -math.Log(1 - eff*(1+cr)) / (1 + cr)
}

func SolveNTU(s Spec, p CapacityPair) (q, hotOut, coldOut, eff, ntu float64) {
	ua := s.UA
	if ua <= 0 || p.Min <= 0 {
		return 0, s.Hot.Inlet, s.Cold.Inlet, 0, 0
	}
	ntu = ua / p.Min
	eff = Effectiveness(s, ntu, p)
	q = eff * p.Min * InletDifference(s)
	hotOut = s.Hot.Inlet - q/p.Hot
	coldOut = s.Cold.Inlet + q/p.Cold
	return
}
