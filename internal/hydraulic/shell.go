package hydraulic

import "math"

type Shell struct {
	ShellDiameter float64
	BaffleSpacing float64
	TubeOD        float64
	TubePitch     float64
	TubeCount     int
	Length        float64
	Passes        int
}

type ShellResult struct {
	CrossFlowArea  float64
	EquivalentDiameter float64
	Velocity       float64
	Reynolds       float64
	Friction       float64
	DeltaP         float64
}

func ShellCrossFlowArea(s Shell) float64 {
	if s.ShellDiameter <= 0 || s.TubePitch <= 0 {
		return 0
	}
	clear := s.TubePitch - s.TubeOD
	if clear <= 0 {
		return 0
	}
	return s.ShellDiameter * s.BaffleSpacing * clear / s.TubePitch
}

func ShellEquivalentDiameter(s Shell) float64 {
	if s.TubePitch <= 0 {
		return 0
	}
	flowArea := s.TubePitch*s.TubePitch - math.Pi*s.TubeOD*s.TubeOD/4
	wetted := math.Pi * s.TubeOD
	if wetted <= 0 {
		return 0
	}
	return 4 * flowArea / wetted
}

func ShellVelocity(s Shell, totalFlow float64) float64 {
	area := ShellCrossFlowArea(s)
	if area <= 0 {
		return 0
	}
	return totalFlow / area
}

func ShellPressureDrop(s Shell, totalFlow, density float64) ShellResult {
	if s.ShellDiameter <= 0 || s.Length <= 0 || density <= 0 {
		return ShellResult{}
	}
	area := ShellCrossFlowArea(s)
	vel := ShellVelocity(s, totalFlow)
	de := ShellEquivalentDiameter(s)
	if area <= 0 || de <= 0 {
		return ShellResult{CrossFlowArea: area, Velocity: vel}
	}
	re := Reynolds(vel, de)
	f := DarcyFriction(vel, de)
	baffles := int(math.Max(float64(s.Length)/s.BaffleSpacing, 1))
	nb := float64(baffles) + 1
	dp := f * nb * (s.ShellDiameter / de) * DynamicHead(vel, density)
	if s.Passes > 0 {
		dp *= float64(s.Passes)
	}
	return ShellResult{
		CrossFlowArea:     area,
		EquivalentDiameter: de,
		Velocity:          vel,
		Reynolds:          re,
		Friction:          f,
		DeltaP:            dp,
	}
}

func NozzleShellDrop(s Shell, flow, density float64, n Nozzle) float64 {
	vn := NozzleVelocity(flow, n)
	main := ShellPressureDrop(s, flow, density)
	return main.DeltaP + EntranceLoss(vn) + ExitLoss(vn)
}
