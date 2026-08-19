package heatx

import "math"

type OverallUA struct {
	Clean  float64
	Fouled float64
	U      float64
	Area   float64
}

func UFromUA(ua, area float64) float64 {
	if area <= 0 {
		return 0
	}
	return ua / area
}

func UAClean(area, u float64) float64 {
	return area * u
}

func ApplyFouling(uaClean float64, rf Fouling) float64 {
	if uaClean <= 0 {
		return 0
	}
	if rf.Total() <= 0 {
		return uaClean
	}
	inv := 1/uaClean + rf.Total()
	return 1 / inv
}

func ComputeOverallUA(area, uClean float64, rf Fouling) OverallUA {
	clean := UAClean(area, uClean)
	fouled := ApplyFouling(clean, rf)
	return OverallUA{
		Clean:  clean,
		Fouled: fouled,
		U:      UFromUA(fouled, area),
		Area:   area,
	}
}

func WithDirectUA(ua float64) Spec {
	return Spec{UA: ua, Flow: Counter}
}

func FoulingRaisesResistance(clean, fouled float64) bool {
	return fouled <= clean+1e-15
}

func RequireArea(ua, u float64) float64 {
	if u <= 0 {
		return math.Inf(1)
	}
	return ua / u
}
