package heatx

import "math"

type CorrectionFactors struct {
	P float64
	R float64
	F float64
}

func TemperatureRatioP(s Spec, hotOut, coldOut float64) float64 {
	den := s.Hot.Inlet - s.Cold.Inlet
	if den == 0 {
		return 0
	}
	return (coldOut - s.Cold.Inlet) / den
}

func CapacityRatioR(s Spec, hotOut, coldOut float64) float64 {
	num := s.Hot.Inlet - hotOut
	den := coldOut - s.Cold.Inlet
	if den == 0 {
		return 0
	}
	return num / den
}

func CorrectionFactorF1_2(s Spec, hotOut, coldOut float64) float64 {
	p := TemperatureRatioP(s, hotOut, coldOut)
	r := CapacityRatioR(s, hotOut, coldOut)
	if p >= 1 {
		return 0
	}
	if r == 1 {
		den := 1 - p + p*math.Log(2-p)
		if den == 0 {
			return 0
		}
		return math.Sqrt(2) * (1 - p) / (2 * den)
	}
	if r == 0 {
		return 1
	}
	term := math.Sqrt(r*r + 1)
	num := math.Log((1 - p*r) / (1 - p))
	den := (r - 1) * math.Log((2-p*(r+1-term))/(2-p*(r+1+term)))
	if den == 0 {
		return 0
	}
	return num * term / den
}

func LMTDFactor(s Spec, hotOut, coldOut float64) float64 {
	if s.Flow == Counter || s.Flow == Parallel {
		return 1
	}
	return CorrectionFactorF1_2(s, hotOut, coldOut)
}

func CorrectedLMTD(s Spec, hotOut, coldOut float64) float64 {
	lmtd := LMTDOf(s, hotOut, coldOut)
	f := LMTDFactor(s, hotOut, coldOut)
	return lmtd * f
}

func MeanTemperatureDifference(s Spec, hotOut, coldOut float64) float64 {
	return CorrectedLMTD(s, hotOut, coldOut)
}
