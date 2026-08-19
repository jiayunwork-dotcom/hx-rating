package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunRateExampleFile(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunRate([]string{"../../example/counter.json"}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "Q") {
		t.Fatalf("output missing Q: %s", out.String())
	}
	if !strings.Contains(out.String(), "feasible") {
		t.Fatalf("output missing feasibility: %s", out.String())
	}
}

func TestRunRateStdin(t *testing.T) {
	input := `{
	  "hot": { "mass_flow": 2.5, "specific_heat": 4186, "inlet_temp": 90 },
	  "cold": { "mass_flow": 2.6, "specific_heat": 4186, "inlet_temp": 30 },
	  "ua": 8000,
	  "flow": "counter"
	}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errBuf.String())
	}
	if !strings.Contains(out.String(), "feasible       true") {
		t.Fatalf("expected feasible true, got:\n%s", out.String())
	}
}

func TestRunRateRejectsBadInput(t *testing.T) {
	input := `{"hot": {"mass_flow": 0, "specific_heat": 4186, "inlet_temp": 80},
	          "cold": {"mass_flow": 1, "specific_heat": 4186, "inlet_temp": 20},
	          "ua": 5000, "flow": "counter"}`
	old := stdinOverride
	stdinOverride = strings.NewReader(input)
	defer func() { stdinOverride = old }()

	var out, errBuf bytes.Buffer
	code := RunRate([]string{}, &out, &errBuf)
	if code == 0 {
		t.Fatal("zero mass flow should exit non-zero")
	}
	if !strings.Contains(errBuf.String(), "error") {
		t.Fatalf("stderr should contain error, got %q", errBuf.String())
	}
}

func TestRunRateUnknownFile(t *testing.T) {
	var out, errBuf bytes.Buffer
	code := RunRate([]string{"-f", "../../example/nope.json"}, &out, &errBuf)
	if code == 0 {
		t.Fatal("missing file should exit non-zero")
	}
}
