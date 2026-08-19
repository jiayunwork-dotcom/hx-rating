package hydraulic

import (
	"math"
	"testing"
)

func TestVelocityScalesWithFlow(t *testing.T) {
	d := 0.02
	v1 := Velocity(0.001, d)
	v2 := Velocity(0.002, d)
	if math.Abs(v2/v1-2.0) > 1e-9 {
		t.Fatalf("doubling flow should double velocity: v1=%.6f v2=%.6f", v1, v2)
	}
}

func TestDropRisesWithVelocity2(t *testing.T) {
	r1 := PressureDrop(0.001, 0.02, 6, 0.00015, 998)
	r2 := PressureDrop(0.002, 0.02, 6, 0.00015, 998)
	ratio := r2.DeltaP / r1.DeltaP
	if r1.DeltaP <= 0 || r2.DeltaP <= r1.DeltaP {
		t.Fatalf("deltaP should grow with flow: r1=%.6f r2=%.6f", r1.DeltaP, r2.DeltaP)
	}
	// friction follows the velocity trend: doubling flow must more than double drop
	if ratio < 1.5 {
		t.Fatalf("doubling flow should raise drop by at least 1.5x, got %.3f", ratio)
	}
	if ratio > 6 {
		t.Fatalf("doubling flow should not raise drop more than 6x, got %.3f", ratio)
	}
}

func TestLaminarFriction(t *testing.T) {
	if f := LaminarFriction(1000); math.Abs(f-0.064) > 1e-9 {
		t.Fatalf("64/Re=0.064 got %v", f)
	}
}

func TestZeroDiameterSafe(t *testing.T) {
	r := PressureDrop(0.001, 0, 6, 0.00015, 998)
	if r.Velocity != 0 || r.DeltaP != 0 {
		t.Fatalf("zero diameter should be safe, got %+v", r)
	}
}

func TestBundleDrop(t *testing.T) {
	b := Bundle{
		TubeDiameter:  0.02,
		TubeLength:    6,
		TubeRoughness: 0.00015,
		TubeCount:     10,
		Passes:        2,
	}
	res := BundleDrop(b, 0.01)
	if res.PerTube.DeltaP <= 0 {
		t.Fatalf("per-tube deltaP should be positive, got %+v", res.PerTube)
	}
	withPasses := BundleDropWithPasses(b, 0.01)
	if withPasses.Bundle.DeltaP <= res.Bundle.DeltaP {
		t.Fatal("extra pass should raise bundle deltaP")
	}
}

func TestMinorLossesPositive(t *testing.T) {
	loss := TotalMinorLoss(1.0, 2, 1)
	if loss <= 0 {
		t.Fatalf("minor losses should be positive, got %v", loss)
	}
	if EntranceLoss(1.0) >= ExitLoss(1.0) {
		t.Fatal("entrance loss should be below exit loss")
	}
}
