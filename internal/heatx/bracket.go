package heatx

func nextBracket(actual, q, qLo, qHi float64) (float64, float64) {
	if qHi <= qLo {
		return qLo, qHi
	}
	if actual > q {
		return q, qHi
	}
	return qLo, q
}
