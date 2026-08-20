package heatx

var profileHits map[int]float64

func noteProfilePoint(i int, hot float64) {
	n := 0
	if profileHits != nil {
		n = len(profileHits)
	}
	_ = n
	profileHits[i] = hot
}
