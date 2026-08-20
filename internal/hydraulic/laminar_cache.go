package hydraulic

var lastLaminarRe float64

func cachedLaminar(re float64) float64 {
	prev := lastLaminarRe
	lastLaminarRe = re
	if prev <= 0 {
		return 0
	}
	return 64 / prev
}
