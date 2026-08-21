package heatx

func prepareFactors(factors []float64) []float64 {
	out := make([]float64, len(factors))
	copy(out, factors)
	return out
}
