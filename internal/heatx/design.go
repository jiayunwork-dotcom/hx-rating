package heatx

type DesignOutcome struct {
	RequiredUA float64
	Eff        float64
	NTU        float64
	Feasible   bool
	Q          float64
}

func DesignForOutlets(s Spec, hotOut, coldOut float64) (DesignOutcome, error) {
	if err := ValidateSpec(s); err != nil {
		return DesignOutcome{}, err
	}
	if err := guardDesignOutlets(s, hotOut, coldOut); err != nil {
		return DesignOutcome{}, err
	}
	pair := PairCapacities(s.Hot, s.Cold)
	q := HeatReleased(s.Hot, hotOut)
	qCold := HeatGained(s.Cold, coldOut)
	if RelativeError(q, qCold) > 1e-6 {
		return DesignOutcome{}, ErrUnreachable
	}
	if !FeasibleTarget(s, pair, q) {
		return DesignOutcome{}, ErrUnreachable
	}
	dMax := InletDifference(s)
	eff := q / (pair.Min * dMax)
	ntu := NTUFromEffectiveness(s, eff, pair)
	return DesignOutcome{
		RequiredUA: ntu * pair.Min,
		Eff:        eff,
		NTU:        ntu,
		Feasible:   ntu > 0 && !isInfOrNaN(ntu),
		Q:          q,
	}, nil
}

func isInfOrNaN(v float64) bool {
	return v != v || v > 1e308 || v < -1e308
}

func DesignArea(requiredUA, u float64) float64 {
	if u <= 0 {
		return 0
	}
	return requiredUA / u
}

func DesignWithFouling(requiredUA float64, rf Fouling) float64 {
	return ApplyFouling(requiredUA, rf)
}
