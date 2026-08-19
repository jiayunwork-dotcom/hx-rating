package heatx

import (
	"math"
	"testing"
)

func TestCounterEnergyClosed(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	o, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	released := HeatReleased(s.Hot, o.HotOut)
	gained := HeatGained(s.Cold, o.ColdOut)
	if math.Abs(released-gained) > 1e-6*math.Max(released, gained) {
		t.Fatalf("energy not closed: released=%.9f gained=%.9f", released, gained)
	}
	if o.Residual != 0 && math.Abs(o.Residual) > 1e-6*math.Max(o.QNtu, 1) {
		t.Fatalf("residual %e too large", o.Residual)
	}
}

func TestFoulingReducesQ(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 3.0, SpecHeat: 4186, Inlet: 95},
		Cold: Stream{MassFlow: 1.5, SpecHeat: 4186, Inlet: 20},
		UA:   12000,
		Flow: Counter,
	}
	clean, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate clean: %v", err)
	}
	s.Rf = Fouling{Hot: 0.00005, Cold: 0.00003}
	dirty, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate dirty: %v", err)
	}
	if dirty.QNtu >= clean.QNtu {
		t.Fatalf("fouled Q=%.6f not below clean Q=%.6f", dirty.QNtu, clean.QNtu)
	}
	if dirty.ColdOut >= clean.ColdOut {
		t.Fatalf("fouled cold_out=%.6f not below clean cold_out=%.6f", dirty.ColdOut, clean.ColdOut)
	}
}

func TestNTULMTDAgree(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	o, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	if o.RelDiff > 1e-6 {
		t.Fatalf("NTU and LMTD disagree: qNtu=%.9f qLmtd=%.9f rel=%.3e", o.QNtu, o.QLmtd, o.RelDiff)
	}
}
