package hydraulic

import "math"

type Result struct {
	Velocity  float64
	Reynolds  float64
	Friction  float64
	DeltaP    float64
	HeadLoss  float64
	Dynamic   float64
}

func DynamicHead(vel, density float64) float64 {
	return density * vel * vel / 2
}

func PressureDrop(flow, diameter, length, roughness, density float64) Result {
	if diameter <= 0 || length <= 0 || density <= 0 {
		return Result{}
	}
	vel := Velocity(flow, diameter)
	re := Reynolds(vel, diameter)
	f := DarcyFriction(vel, diameter)
	dp := f * (length / diameter) * DynamicHead(vel, density)
	return Result{
		Velocity: vel,
		Reynolds: re,
		Friction: f,
		DeltaP:   dp,
		HeadLoss: dp / (density * Gravity),
		Dynamic:  DynamicHead(vel, density),
	}
}

func PressureDropMass(massRate, density, diameter, length, roughness float64) Result {
	vol := VolumetricRate(massRate, density)
	return PressureDrop(vol, diameter, length, roughness, density)
}

func TubeSideDrop(flow, diameter, length, roughness float64) Result {
	return PressureDrop(flow, diameter, length, roughness, WaterDensity)
}

func TrendVelocity2(flow, diameter float64) float64 {
	vel := Velocity(flow, diameter)
	return vel * vel
}

func SameTrend(a, b float64) bool {
	return math.Abs(a-b) < 1e-6
}
