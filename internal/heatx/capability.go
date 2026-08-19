package heatx

import "math"

type Capability struct {
	HotOut   float64
	ColdOut  float64
	Q        float64
	Eff      float64
	UA       float64
	NTU      float64
	Feasible bool
}

func RequiredUAForOutlets(s Spec, hotOut, coldOut float64) (Capability, error) {
	design, err := DesignForOutlets(s, hotOut, coldOut)
	if err != nil {
		return Capability{}, err
	}
	return Capability{
		HotOut:   hotOut,
		ColdOut:  coldOut,
		Q:        design.Q,
		Eff:      design.Eff,
		UA:       design.RequiredUA,
		NTU:      design.NTU,
		Feasible: design.Feasible,
	}, nil
}

func RequiredUAForQ(s Spec, targetQ float64) (Capability, error) {
	if err := ValidateSpec(s); err != nil {
		return Capability{}, err
	}
	p := PairCapacities(s.Hot, s.Cold)
	if !FeasibleTarget(s, p, targetQ) {
		return Capability{}, ErrUnreachable
	}
	dMax := InletDifference(s)
	eff := targetQ / (p.Min * dMax)
	ntu := NTUFromEffectiveness(s, eff, p)
	hotOut := s.Hot.Inlet - targetQ/p.Hot
	coldOut := s.Cold.Inlet + targetQ/p.Cold
	return Capability{
		HotOut:   hotOut,
		ColdOut:  coldOut,
		Q:        targetQ,
		Eff:      eff,
		UA:       ntu * p.Min,
		NTU:      ntu,
		Feasible: ntu > 0 && !math.IsNaN(ntu) && !math.IsInf(ntu, 0),
	}, nil
}

func AreaForUA(ua, u float64) float64 {
	if u <= 0 {
		return math.Inf(1)
	}
	return ua / u
}

func UAForArea(area, u float64) float64 {
	return area * u
}

func EffectiveUA(clean, rfTotal float64) float64 {
	if clean <= 0 {
		return 0
	}
	if rfTotal <= 0 {
		return clean
	}
	return 1 / (1/clean + rfTotal)
}

func MarginOf(capability Capability, availableUA float64) float64 {
	if capability.UA <= 0 {
		return 0
	}
	return (availableUA - capability.UA) / capability.UA
}
