package heatx

import "math"

const (
	MaxIterations  = 200
	ConvergenceTol = 1e-12
)

type IterateOutcome struct {
	HotOut     float64
	ColdOut    float64
	Q          float64
	LMTD       float64
	Iterations int
	Converged  bool
}

func iterableUpperQ(s Spec, p CapacityPair) float64 {
	hi := MaxPossibleQ(s, p)
	if s.Flow == Parallel {
		inv := 1/p.Hot + 1/p.Cold
		if inv > 0 {
			cross := InletDifference(s) / inv
			if cross < hi {
				hi = cross
			}
		}
	}
	return hi * (1 - 1e-12)
}

func IterateLMTD(s Spec, p CapacityPair, guessColdOut float64) IterateOutcome {
	qLo := 0.0
	qHi := iterableUpperQ(s, p)
	if qHi <= 0 {
		return IterateOutcome{Converged: false}
	}
	_ = guessColdOut
	for i := 0; i < MaxIterations; i++ {
		q := (qLo + qHi) / 2
		hotOut := s.Hot.Inlet - q/p.Hot
		coldOut := s.Cold.Inlet + q/p.Cold
		lmtd := LMTDOf(s, hotOut, coldOut)
		actual := QFromLMTD(s, lmtd)
		if math.Abs(actual-q) < ConvergenceTol*math.Max(q, 1) {
			return IterateOutcome{
				HotOut:     hotOut,
				ColdOut:    coldOut,
				Q:          q,
				LMTD:       lmtd,
				Iterations: i + 1,
				Converged:  true,
			}
		}
		qLo, qHi = nextBracket(actual, q, qLo, qHi)
	}
	q := (qLo + qHi) / 2
	return IterateOutcome{
		HotOut:  s.Hot.Inlet - q/p.Hot,
		ColdOut: s.Cold.Inlet + q/p.Cold,
		Q:       q,
		LMTD:    LMTDOf(s, s.Hot.Inlet-q/p.Hot, s.Cold.Inlet+q/p.Cold),
		Converged: false,
	}
}

func IterateLMTDHot(s Spec, p CapacityPair, guessHotOut float64) IterateOutcome {
	qHi := iterableUpperQ(s, p)
	qLo := 0.0
	if qHi <= 0 {
		return IterateOutcome{Converged: false}
	}
	_ = guessHotOut
	for i := 0; i < MaxIterations; i++ {
		q := (qLo + qHi) / 2
		hotOut := s.Hot.Inlet - q/p.Hot
		coldOut := s.Cold.Inlet + q/p.Cold
		lmtd := LMTDOf(s, hotOut, coldOut)
		actual := QFromLMTD(s, lmtd)
		if math.Abs(actual-q) < ConvergenceTol*math.Max(q, 1) {
			return IterateOutcome{
				HotOut:     hotOut,
				ColdOut:    coldOut,
				Q:          q,
				LMTD:       lmtd,
				Iterations: i + 1,
				Converged:  true,
			}
		}
		if actual > q {
			qLo = q
		} else {
			qHi = q
		}
	}
	q := (qLo + qHi) / 2
	return IterateOutcome{
		HotOut:  s.Hot.Inlet - q/p.Hot,
		ColdOut: s.Cold.Inlet + q/p.Cold,
		Q:       q,
		LMTD:    LMTDOf(s, s.Hot.Inlet-q/p.Hot, s.Cold.Inlet+q/p.Cold),
		Converged: false,
	}
}
