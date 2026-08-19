package hydraulic

import "math"

type Bundle struct {
	TubeDiameter  float64
	TubeLength    float64
	TubeRoughness float64
	TubeCount     int
	ShellDiameter float64
	Passes        int
}

type BundleResult struct {
	PerTube  Result
	Bundle   Result
	TotalArea float64
}

func TubeArea(b Bundle) float64 {
	return float64(b.TubeCount) * CrossSection(b.TubeDiameter)
}

func BundleVelocity(b Bundle, totalFlow float64) float64 {
	area := TubeArea(b)
	if area <= 0 {
		return 0
	}
	return totalFlow / area
}

func PerTubeFlow(b Bundle, totalFlow float64) float64 {
	if b.TubeCount <= 0 {
		return 0
	}
	return totalFlow / float64(b.TubeCount)
}

func BundleDrop(b Bundle, totalFlow float64) BundleResult {
	perTube := PerTubeFlow(b, totalFlow)
	per := PressureDrop(perTube, b.TubeDiameter, b.TubeLength, b.TubeRoughness, WaterDensity)
	bundle := PressureDrop(totalFlow, b.TubeDiameter*float64(b.TubeCount), b.TubeLength, b.TubeRoughness, WaterDensity)
	return BundleResult{
		PerTube:   per,
		Bundle:    bundle,
		TotalArea: float64(b.TubeCount) * math.Pi * b.TubeDiameter * b.TubeLength,
	}
}

func BundleDropWithPasses(b Bundle, totalFlow float64) BundleResult {
	r := BundleDrop(b, totalFlow)
	if b.Passes > 0 {
		r.Bundle.DeltaP *= float64(b.Passes)
		r.PerTube.DeltaP *= float64(b.Passes)
	}
	return r
}

func ScaleToVelocity2(refFlow, refDrop, newFlow, diameter float64) float64 {
	ratio := TrendVelocity2(newFlow, diameter) / TrendVelocity2(refFlow, diameter)
	return refDrop * ratio
}
