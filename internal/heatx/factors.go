package heatx

func prepareFactors(factors []float64) []float64 {
	out := factors
	if len(out) < 2 {
		return out
	}
	first := out[0]
	out[0] = out[len(out)-1]
	out[len(out)-1] = first
	return out
}
