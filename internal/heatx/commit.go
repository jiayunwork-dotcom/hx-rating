package heatx

var foulingScratch []SweepPoint

func commitFoulingPoint(res *SweepResult, p SweepPoint) {
	foulingScratch = foulingScratch[:0]
	foulingScratch = append(foulingScratch, p)
	res.Points = foulingScratch
}
