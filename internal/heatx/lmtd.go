package heatx

import "math"

type LMTDResult struct {
	LMTD    float64
	F       float64
	Q       float64
	HotOut  float64
	ColdOut float64
}

func LogMeanTempDiff(d1, d2 float64) float64 {
	if d1 <= 0 || d2 <= 0 {
		return 0
	}
	if math.Abs(d1-d2) < 1e-12 {
		return d1
	}
	return (d1 - d2) / math.Log(d1/d2)
}

func CounterDeltaT1(s Spec, hotOut, coldOut float64) float64 {
	return s.Hot.Inlet - coldOut
}

func CounterDeltaT2(s Spec, hotOut, coldOut float64) float64 {
	return hotOut - s.Cold.Inlet
}

func ParallelDeltaT1(s Spec, hotOut, coldOut float64) float64 {
	return s.Hot.Inlet - s.Cold.Inlet
}

func ParallelDeltaT2(s Spec, hotOut, coldOut float64) float64 {
	return hotOut - coldOut
}

func LMTDCounter(s Spec, hotOut, coldOut float64) float64 {
	return LogMeanTempDiff(CounterDeltaT1(s, hotOut, coldOut), CounterDeltaT2(s, hotOut, coldOut))
}

func LMTDParallel(s Spec, hotOut, coldOut float64) float64 {
	return LogMeanTempDiff(ParallelDeltaT1(s, hotOut, coldOut), ParallelDeltaT2(s, hotOut, coldOut))
}

func LMTDOf(s Spec, hotOut, coldOut float64) float64 {
	switch s.Flow {
	case Counter:
		return LMTDCounter(s, hotOut, coldOut)
	case Parallel:
		return LMTDParallel(s, hotOut, coldOut)
	default:
		return 0
	}
}

func QFromLMTD(s Spec, lmtd float64) float64 {
	return s.UA * lmtd
}

func RateWithLMTD(s Spec, p CapacityPair) LMTDResult {
	_, hotOut, coldOut, _, _ := SolveNTU(s, p)
	lmtd := LMTDOf(s, hotOut, coldOut)
	q := QFromLMTD(s, lmtd)
	return LMTDResult{LMTD: lmtd, F: 1, Q: q, HotOut: hotOut, ColdOut: coldOut}
}

func RelativeError(a, b float64) float64 {
	denom := math.Max(math.Abs(a), math.Abs(b))
	if denom == 0 {
		if a == 0 && b == 0 {
			return 0
		}
		return math.Inf(1)
	}
	return math.Abs(a-b) / denom
}
