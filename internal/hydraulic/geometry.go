package hydraulic

import "math"

const (
	WaterDensity   = 998.0
	WaterViscosity = 0.0008
	Gravity        = 9.81
)

type Geometry struct {
	Diameter  float64
	Length    float64
	Roughness float64
}

func CrossSection(diameter float64) float64 {
	return math.Pi * diameter * diameter / 4
}

func WettedPerimeter(diameter float64) float64 {
	return math.Pi * diameter
}

func HydraulicDiameter(diameter float64) float64 {
	if diameter <= 0 {
		return 0
	}
	return 4 * CrossSection(diameter) / WettedPerimeter(diameter)
}

func FlowArea(diameter float64) float64 {
	return CrossSection(diameter)
}

func Velocity(flow, diameter float64) float64 {
	area := CrossSection(diameter)
	if area <= 0 {
		return 0
	}
	return flow / area
}

func MassRate(flow, density float64) float64 {
	return flow * density
}

func VolumetricRate(mass, density float64) float64 {
	if density <= 0 {
		return 0
	}
	return mass / density
}
