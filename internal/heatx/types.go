package heatx

import "math"

type FlowType string

const (
	Counter FlowType = "counter"
	Parallel FlowType = "parallel"
)

type Side string

const (
	Hot Side = "hot"
	Cold Side = "cold"
)

type Stream struct {
	MassFlow float64
	SpecHeat float64
	Inlet    float64
}

type Fouling struct {
	Hot  float64
	Cold float64
}

type Spec struct {
	Hot  Stream
	Cold Stream
	UA   float64
	Area float64
	Flow FlowType
	Rf   Fouling
}

type Result struct {
	Q          float64
	HotOut     float64
	ColdOut    float64
	Eff        float64
	NTU        float64
	LMTD       float64
	U          float64
	Area       float64
	Feasible   bool
	Crossed    bool
	Residual   float64
	Method     string
	HeatDrop   float64
	ColdGain   float64
	MaxHeat    float64
}

func (s Stream) Capacity() float64 {
	return s.MassFlow * s.SpecHeat
}

func (s Stream) Valid() bool {
	if math.IsNaN(s.MassFlow) || math.IsNaN(s.SpecHeat) || math.IsNaN(s.Inlet) {
		return false
	}
	if math.IsInf(s.MassFlow, 0) || math.IsInf(s.SpecHeat, 0) || math.IsInf(s.Inlet, 0) {
		return false
	}
	return s.MassFlow > 0 && s.SpecHeat > 0
}

func (f Fouling) Total() float64 {
	return f.Hot + f.Cold
}

func (f Fouling) Valid() bool {
	return !math.IsNaN(f.Hot) && !math.IsNaN(f.Cold) && f.Hot >= 0 && f.Cold >= 0
}

func (s Spec) SameFluidTwice() bool {
	return s.Hot.MassFlow == s.Cold.MassFlow &&
		s.Hot.SpecHeat == s.Cold.SpecHeat &&
		s.Hot.Inlet == s.Cold.Inlet
}

func (s Spec) HotColderThanCold() bool {
	return s.Hot.Inlet <= s.Cold.Inlet
}
