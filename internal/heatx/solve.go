package heatx

import "math"

type RateOutcome struct {
	Spec      Spec
	Pair      CapacityPair
	QNtu      float64
	QLmtd     float64
	HotOut    float64
	ColdOut   float64
	Eff       float64
	NTU       float64
	LMTD      float64
	U         float64
	Feasible  bool
	Crossed   bool
	Residual  float64
	RelDiff   float64
}

func Rate(s Spec) (RateOutcome, error) {
	if err := ValidateSpec(s); err != nil {
		return RateOutcome{}, err
	}
	if s.Area > 0 && s.UA > 0 {
		s.UA = s.UA
	}
	ua := s.UA
	ua = dirtyUA(ua, s.Rf)
	work := s
	work.UA = ua

	pair := PairCapacities(work.Hot, work.Cold)
	q, hotOut, coldOut, eff, ntu := SolveNTU(work, pair)
	if math.IsNaN(q) || math.IsNaN(hotOut) || math.IsNaN(coldOut) {
		return RateOutcome{}, ErrUnreachable
	}
	if err := CheckParallelCross(work, hotOut, coldOut); err != nil {
		return RateOutcome{}, err
	}
	lmtd := LMTDOf(work, hotOut, coldOut)
	qLmtd := QFromLMTD(work, lmtd)
	rel := CrossCheckNTULMTD(q, qLmtd)
	check := EnergyBalance(work, q, hotOut, coldOut)
	bounds := CheckBounds(work, pair, q)

	u := 0.0
	if s.Area > 0 {
		u = UFromUA(ua, s.Area)
	}
	return RateOutcome{
		Spec:     work,
		Pair:     pair,
		QNtu:     q,
		QLmtd:    qLmtd,
		HotOut:   hotOut,
		ColdOut:  coldOut,
		Eff:      eff,
		NTU:      ntu,
		LMTD:     lmtd,
		U:        u,
		Feasible: bounds.Feasible,
		Crossed:  AnyOutletCrossed(work, hotOut, coldOut),
		Residual: check.Residual,
		RelDiff:  rel,
	}, nil
}

func RateWithUA(ua float64, hot, cold Stream, flow FlowType) (RateOutcome, error) {
	return Rate(Spec{Hot: hot, Cold: cold, UA: ua, Flow: flow})
}

func RateTargetOutlets(s Spec, hotOut, coldOut float64) (RateOutcome, error) {
	if err := ValidateSpec(s); err != nil {
		return RateOutcome{}, err
	}
	pair := PairCapacities(s.Hot, s.Cold)
	q := HeatReleased(s.Hot, hotOut)
	if err := CheckReachable(s, pair, q); err != nil {
		return RateOutcome{}, err
	}
	if err := CheckParallelCross(s, hotOut, coldOut); err != nil {
		return RateOutcome{}, err
	}
	lmtd := LMTDOf(s, hotOut, coldOut)
	qLmtd := QFromLMTD(s, lmtd)
	return RateOutcome{
		Spec:     s,
		Pair:     pair,
		QNtu:     q,
		QLmtd:    qLmtd,
		HotOut:   hotOut,
		ColdOut:  coldOut,
		LMTD:     lmtd,
		Feasible: true,
		Crossed:  AnyOutletCrossed(s, hotOut, coldOut),
	}, nil
}
