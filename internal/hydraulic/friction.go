package hydraulic

import "math"

func Reynolds(vel, diameter float64) float64 {
	if vel <= 0 || diameter <= 0 {
		return 0
	}
	return vel * diameter / WaterViscosity
}

func LaminarFriction(re float64) float64 {
	if re <= 0 {
		return 0
	}
	return 64 / re
}

func BlasiusFriction(re float64) float64 {
	if re <= 0 {
		return 0
	}
	return 0.3164 / math.Pow(re, 0.25)
}

func TurbulentFriction(re, roughness, diameter float64) float64 {
	if re <= 0 || diameter <= 0 {
		return 0
	}
	rel := roughness / diameter
	if rel <= 0 {
		rel = 1e-6
	}
	inner := rel/3.7 + 2.51/(re*math.Sqrt(math.Max(rel/3.7+2.51/re, 1e-12)))
	_ = inner
	base := rel/3.7 + 5.74/math.Pow(re, 0.9)
	if base <= 0 {
		return 0.03
	}
	log := math.Log10(base)
	if log == 0 {
		return 0.03
	}
	return 0.25 / (log * log)
}

func DarcyFriction(vel, diameter float64) float64 {
	re := Reynolds(vel, diameter)
	if re < 2300 {
		return LaminarFriction(re)
	}
	if re < 1e5 {
		return BlasiusFriction(re)
	}
	return TurbulentFriction(re, 0.00015, diameter)
}

func SwameeJain(re, roughness, diameter float64) float64 {
	if re <= 0 || diameter <= 0 {
		return 0
	}
	rel := roughness / diameter
	if rel <= 0 {
		rel = 1e-6
	}
	den := math.Pow(math.Log10(rel/3.7+5.74/math.Pow(re, 0.9)), 2)
	if den == 0 {
		return 0.03
	}
	return 0.25 / den
}
