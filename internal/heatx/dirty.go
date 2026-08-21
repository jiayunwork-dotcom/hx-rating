package heatx

func dirtyUA(uaClean float64, rf Fouling) float64 {
	return ApplyFouling(uaClean, rf)
}
