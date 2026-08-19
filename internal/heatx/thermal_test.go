package heatx

import (
	"math"
	"testing"
)

func TestParallelEnergyClosed(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1.8, SpecHeat: 4186, Inlet: 85},
		Cold: Stream{MassFlow: 3.0, SpecHeat: 4186, Inlet: 25},
		UA:   5000,
		Flow: Parallel,
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
}

func TestCounterBetterThanParallel(t *testing.T) {
	counter := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	parallel := counter
	parallel.Flow = Parallel
	co, err := Rate(counter)
	if err != nil {
		t.Fatalf("Rate counter: %v", err)
	}
	po, err := Rate(parallel)
	if err != nil {
		t.Fatalf("Rate parallel: %v", err)
	}
	if co.QNtu <= po.QNtu {
		t.Fatalf("counter Q=%.6f should exceed parallel Q=%.6f", co.QNtu, po.QNtu)
	}
}

func TestLargeUAApproachesLimit(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 30},
		UA:   1e9,
		Flow: Counter,
	}
	o, err := Rate(s)
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	maxQ := MaxPossibleQ(s, o.Pair)
	if o.QNtu > maxQ*1.001 {
		t.Fatalf("Q=%.6f exceeds physical max %.6f", o.QNtu, maxQ)
	}
	if o.HotOut < s.Cold.Inlet-1e-3 || o.ColdOut > s.Hot.Inlet+1e-3 {
		t.Fatalf("outlets crossed the inlet span: hot_out=%.4f cold_out=%.4f", o.HotOut, o.ColdOut)
	}
}

func TestFoulingUnitConsistency(t *testing.T) {
	clean := ComputeOverallUA(20, 500, Fouling{})
	if math.Abs(clean.Clean-10000) > 1e-9 {
		t.Fatalf("clean UA=%.6f want 10000", clean.Clean)
	}
	fouled := ApplyFouling(10000, Fouling{Hot: 0.00001, Cold: 0.00001})
	if fouled >= 10000 {
		t.Fatalf("fouled UA=%.6f should be below clean", fouled)
	}
	inv := 1/fouled - 1.0/10000
	if math.Abs(inv-0.00002) > 1e-12 {
		t.Fatalf("resistance increment %.9f want 0.00002", inv)
	}
}

func TestDesignInfeasible(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 60},
		Cold: Stream{MassFlow: 1, SpecHeat: 4186, Inlet: 20},
		UA:   5000,
		Flow: Counter,
	}
	if _, err := DesignForOutlets(s, 5, 90); err == nil {
		t.Fatal("design with cold_out above hot_in should be rejected")
	}
	if _, err := DesignForOutlets(s, 60, 30); err == nil {
		t.Fatal("design with hot_out at inlet should be rejected")
	}
}

func TestDesignConsistentWithRate(t *testing.T) {
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
	d, err := DesignForOutlets(s, o.HotOut, o.ColdOut)
	if err != nil {
		t.Fatalf("DesignForOutlets: %v", err)
	}
	if math.Abs(d.RequiredUA-s.UA)/s.UA > 0.02 {
		t.Fatalf("required UA=%.6f not close to spec UA=%.6f", d.RequiredUA, s.UA)
	}
}
