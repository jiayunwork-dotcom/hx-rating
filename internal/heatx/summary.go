package heatx

type Summary struct {
	Flow       FlowType
	Q          float64
	HotIn      float64
	HotOut     float64
	ColdIn     float64
	ColdOut    float64
	Eff        float64
	NTU        float64
	LMTD       float64
	UA         float64
	U          float64
	Area       float64
	Feasible   bool
	Crossed    bool
	Residual   float64
	NTULMTDDiff float64
}

func Summarize(o RateOutcome) Summary {
	return Summary{
		Flow:        o.Spec.Flow,
		Q:           o.QNtu,
		HotIn:       o.Spec.Hot.Inlet,
		HotOut:      o.HotOut,
		ColdIn:      o.Spec.Cold.Inlet,
		ColdOut:     o.ColdOut,
		Eff:         o.Eff,
		NTU:         o.NTU,
		LMTD:        o.LMTD,
		UA:          o.Spec.UA,
		U:           o.U,
		Area:        o.Spec.Area,
		Feasible:    o.Feasible,
		Crossed:     o.Crossed,
		Residual:    o.Residual,
		NTULMTDDiff: o.RelDiff,
	}
}
