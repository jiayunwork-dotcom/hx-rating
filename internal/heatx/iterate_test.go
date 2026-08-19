package heatx

import (
	"math"
	"testing"
)

func TestIterateLMTDConverges(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	p := PairCapacities(s.Hot, s.Cold)
	o, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	it := IterateLMTD(s, p, 40)
	if !it.Converged {
		t.Fatalf("iteration did not converge: %+v", it)
	}
	if math.Abs(it.ColdOut-o.ColdOut) > 1e-6 {
		t.Fatalf("iterated cold_out=%.9f want %.9f", it.ColdOut, o.ColdOut)
	}
	if math.Abs(it.HotOut-o.HotOut) > 1e-6 {
		t.Fatalf("iterated hot_out=%.9f want %.9f", it.HotOut, o.HotOut)
	}
}

func TestIterateLMTDParallel(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1.8, SpecHeat: 4186, Inlet: 85},
		Cold: Stream{MassFlow: 3.0, SpecHeat: 4186, Inlet: 25},
		UA:   5000,
		Flow: Parallel,
	}
	p := PairCapacities(s.Hot, s.Cold)
	o, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	it := IterateLMTD(s, p, 30)
	if !it.Converged {
		t.Fatalf("parallel iteration did not converge: %+v", it)
	}
	if math.Abs(it.ColdOut-o.ColdOut) > 1e-6 {
		t.Fatalf("iterated cold_out=%.9f want %.9f", it.ColdOut, o.ColdOut)
	}
}

func TestTemperatureProfileMonotonic(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	p := PairCapacities(s.Hot, s.Cold)
	pts := TemperatureProfile(s, p, 11)
	if len(pts) != 11 {
		t.Fatalf("want 11 profile points, got %d", len(pts))
	}
	if pts[0].Hot != 90 || pts[0].Cold != 30 {
		t.Fatalf("inlet endpoints wrong: %+v", pts[0])
	}
	o, _ := Rate(s)
	last := pts[len(pts)-1]
	if math.Abs(last.Hot-o.HotOut) > 1e-6 || math.Abs(last.Cold-o.ColdOut) > 1e-6 {
		t.Fatalf("outlet endpoint wrong: %+v want hot=%.4f cold=%.4f", last, o.HotOut, o.ColdOut)
	}
}

func TestCorrectionFactorUnit(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	if f := LMTDFactor(s, 67, 52); f != 1 {
		t.Fatalf("counter flow F factor should be 1, got %v", f)
	}
	p := TemperatureRatioP(s, 67, 52)
	r := CapacityRatioR(s, 67, 52)
	if p <= 0 || p >= 1 {
		t.Fatalf("P out of range: %v", p)
	}
	if r <= 0 {
		t.Fatalf("R out of range: %v", r)
	}
}

func TestRequiredUAForQ(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	o, _ := Rate(s)
	cap, err := RequiredUAForQ(s, o.QNtu)
	if err != nil {
		t.Fatalf("RequiredUAForQ: %v", err)
	}
	if math.Abs(cap.UA-s.UA)/s.UA > 0.02 {
		t.Fatalf("required UA=%.6f not close to spec UA=%.6f", cap.UA, s.UA)
	}
	if _, err := RequiredUAForQ(s, MaxPossibleQ(s, o.Pair)*2); err != ErrUnreachable {
		t.Fatalf("want ErrUnreachable for oversized target Q, got %v", err)
	}
}
