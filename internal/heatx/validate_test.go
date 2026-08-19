package heatx

import (
	"math"
	"testing"
)

func TestRejectZeroUA(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 80},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 20},
		UA:   0,
		Flow: Counter,
	}
	if _, err := Rate(s); err != ErrInvalidUA {
		t.Fatalf("want ErrInvalidUA, got %v", err)
	}
}

func TestRejectNegativeFlow(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: -1, SpecHeat: 4186, Inlet: 80},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 20},
		UA:   5000,
		Flow: Counter,
	}
	if _, err := Rate(s); err != ErrInvalidMassFlow {
		t.Fatalf("want ErrInvalidMassFlow, got %v", err)
	}
}

func TestRejectMissingInlet(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: math.NaN()},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 20},
		UA:   5000,
		Flow: Counter,
	}
	if _, err := Rate(s); err != ErrMissingInlet {
		t.Fatalf("want ErrMissingInlet, got %v", err)
	}
}

func TestRejectIdenticalStreams(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 60},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 60},
		UA:   5000,
		Flow: Counter,
	}
	if _, err := Rate(s); err != ErrSameFluidTwice {
		t.Fatalf("want ErrSameFluidTwice, got %v", err)
	}
}

func TestRejectHotNotHot(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 20},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 80},
		UA:   5000,
		Flow: Counter,
	}
	if _, err := Rate(s); err != ErrHotNotHot {
		t.Fatalf("want ErrHotNotHot, got %v", err)
	}
}

func TestParallelCrossRejected(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 0.5, SpecHeat: 4186, Inlet: 80},
		Cold: Stream{MassFlow: 4.0, SpecHeat: 4186, Inlet: 20},
		UA:   40000,
		Flow: Parallel,
	}
	if !ParallelOutletsCrossed(55, 60) {
		t.Fatal("parallel outlets 55/60 should be flagged as crossed")
	}
	if err := CheckParallelCross(s, 55, 60); err != ErrParallelCross {
		t.Fatalf("want ErrParallelCross, got %v", err)
	}
	if err := CheckParallelCross(s, 60, 55); err != nil {
		t.Fatalf("non-crossed outlets should pass, got %v", err)
	}
}

func TestUnreachableTarget(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 60},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 20},
		UA:   5000,
		Flow: Counter,
	}
	p := PairCapacities(s.Hot, s.Cold)
	tooMuch := MaxPossibleQ(s, p) * 2
	if _, err := RateTargetOutlets(s, 10, 100); err == nil {
		t.Fatal("want error for infeasible target outlets")
	}
	_ = tooMuch
}

func TestEffectivenessLimits(t *testing.T) {
	if eff := EffectivenessCounter(5, 1.0); eff != 5.0/6.0 {
		t.Fatalf("counter Cr=1 limit: got %v", eff)
	}
	if eff := EffectivenessCounter(0, 0.5); eff != 0 {
		t.Fatalf("NTU=0 effectiveness should be 0, got %v", eff)
	}
}
