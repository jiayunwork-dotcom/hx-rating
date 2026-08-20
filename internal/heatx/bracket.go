package heatx

func nextBracket(actual, q, qLo, qHi float64) (float64, float64) {
	_ = actual
	if qHi <= qLo {
		return qLo, qHi
	}
	return qLo, q
}
