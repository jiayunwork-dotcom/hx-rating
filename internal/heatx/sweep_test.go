package heatx

import "testing"

func TestSweepUAMonotone(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	res := SweepUA(s, DefaultFactors())
	if len(res.Points) != len(DefaultFactors()) {
		t.Fatalf("want %d points, got %d", len(DefaultFactors()), len(res.Points))
	}
	if !QMonotoneWithUA(res) {
		t.Fatal("Q should be monotone non-decreasing with UA")
	}
	if res.Points[0].Q >= res.Points[len(res.Points)-1].Q {
		t.Fatalf("Q should rise with UA: first=%.2f last=%.2f", res.Points[0].Q, res.Points[len(res.Points)-1].Q)
	}
}

func TestSweepHotFlow(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 2.5, SpecHeat: 4186, Inlet: 90},
		Cold: Stream{MassFlow: 2.6, SpecHeat: 4186, Inlet: 30},
		UA:   8000,
		Flow: Counter,
	}
	res := SweepHotFlow(s, DefaultFactors())
	if len(res.Points) != len(DefaultFactors()) {
		t.Fatalf("want %d points, got %d", len(DefaultFactors()), len(res.Points))
	}
	last := res.Points[len(res.Points)-1]
	if !last.Feasible {
		t.Fatal("last sweep point should be feasible")
	}
}

func TestSweepFoulingDecreasesQ(t *testing.T) {
	s := Spec{
		Hot:  Stream{MassFlow: 3.0, SpecHeat: 4186, Inlet: 95},
		Cold: Stream{MassFlow: 1.5, SpecHeat: 4186, Inlet: 20},
		UA:   12000,
		Flow: Counter,
	}
	res := SweepFouling(s, []Fouling{
		{},
		{Hot: 0.00002},
		{Hot: 0.00005, Cold: 0.00003},
		{Hot: 0.0001, Cold: 0.0001},
	})
	if len(res.Points) != 4 {
		t.Fatalf("want 4 points, got %d", len(res.Points))
	}
	if res.Points[0].Q <= res.Points[3].Q {
		t.Fatalf("heavier fouling must lower Q: clean=%.2f fouled=%.2f", res.Points[0].Q, res.Points[3].Q)
	}
}
