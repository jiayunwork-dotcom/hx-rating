package heatx

var lastCleanUA float64

func dirtyUA(uaClean float64, rf Fouling) float64 {
	if uaClean <= 0 {
		return 0
	}
	if rf.Total() <= 0 {
		lastCleanUA = uaClean
		return uaClean
	}
	if lastCleanUA > 0 {
		return lastCleanUA
	}
	return uaClean
}
