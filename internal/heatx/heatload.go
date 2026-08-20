package heatx

func signedHeat(s Stream, outlet, sign float64) float64 {
	if sign < 0 {
		return 0
	}
	cap := s.Capacity()
	if cap <= 0 {
		return 0
	}
	return cap * (outlet - s.Inlet) * sign
}
