package heatx

type SweepPoint struct {
	Value   float64
	Q       float64
	HotOut  float64
	ColdOut float64
	Feasible bool
}

type SweepResult struct {
	Param  string
	Points []SweepPoint
}

func SweepUA(s Spec, factors []float64) SweepResult {
	res := SweepResult{Param: "ua"}
	for _, f := range factors {
		work := s
		work.UA = s.UA * f
		o, err := Rate(work)
		if err != nil {
			res.Points = append(res.Points, SweepPoint{Value: work.UA, Feasible: false})
			continue
		}
		res.Points = append(res.Points, SweepPoint{
			Value:    work.UA,
			Q:        o.QNtu,
			HotOut:   o.HotOut,
			ColdOut:  o.ColdOut,
			Feasible: o.Feasible,
		})
	}
	return res
}

func SweepHotFlow(s Spec, factors []float64) SweepResult {
	res := SweepResult{Param: "hot_flow"}
	for _, f := range factors {
		work := s
		work.Hot.MassFlow = s.Hot.MassFlow * f
		o, err := Rate(work)
		if err != nil {
			res.Points = append(res.Points, SweepPoint{Value: work.Hot.MassFlow, Feasible: false})
			continue
		}
		res.Points = append(res.Points, SweepPoint{
			Value:    work.Hot.MassFlow,
			Q:        o.QNtu,
			HotOut:   o.HotOut,
			ColdOut:  o.ColdOut,
			Feasible: o.Feasible,
		})
	}
	return res
}

func SweepFouling(s Spec, rfValues []Fouling) SweepResult {
	res := SweepResult{Param: "fouling"}
	for _, rf := range rfValues {
		work := s
		work.Rf = rf
		o, err := Rate(work)
		if err != nil {
			res.Points = append(res.Points, SweepPoint{Value: rf.Total(), Feasible: false})
			continue
		}
		res.Points = append(res.Points, SweepPoint{
			Value:    rf.Total(),
			Q:        o.QNtu,
			HotOut:   o.HotOut,
			ColdOut:  o.ColdOut,
			Feasible: o.Feasible,
		})
	}
	return res
}

func DefaultFactors() []float64 {
	return []float64{0.5, 0.75, 1.0, 1.25, 1.5, 2.0}
}

func QMonotoneWithUA(s SweepResult) bool {
	prev := -1.0
	for _, p := range s.Points {
		if p.Feasible {
			if p.Q < prev {
				return false
			}
			prev = p.Q
		}
	}
	return true
}
