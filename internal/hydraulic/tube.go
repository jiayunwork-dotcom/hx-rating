package hydraulic

import "math"

type Tube struct {
	ID        float64
	OD        float64
	Length    float64
	Roughness float64
	Count     int
}

func TubeInnerArea(t Tube) float64 {
	return float64(t.Count) * math.Pi * t.ID * t.Length
}

func TubeOuterArea(t Tube) float64 {
	return float64(t.Count) * math.Pi * t.OD * t.Length
}

func TubeVolume(t Tube) float64 {
	return float64(t.Count) * CrossSection(t.ID) * t.Length
}

func BundleResidenceTime(t Tube, flow float64) float64 {
	vol := TubeVolume(t)
	if vol <= 0 {
		return 0
	}
	return vol / flow
}

func TubeSideVelocityCheck(t Tube, flow, maxVel float64) (float64, bool) {
	vel := Velocity(flow, t.ID)
	return vel, vel <= maxVel
}

func HeatTransferAreaForDuty(duty, lmtd, u float64) float64 {
	if lmtd <= 0 || u <= 0 {
		return 0
	}
	return duty / (u * lmtd)
}

func TubesForArea(area, od, length float64) int {
	if od <= 0 || length <= 0 {
		return 0
	}
	per := math.Pi * od * length
	if per <= 0 {
		return 0
	}
	return int(math.Ceil(area / per))
}

func LayoutDiameter(tubeCount int, pitch, od float64) float64 {
	if tubeCount <= 0 {
		return 0
	}
	n := math.Ceil(math.Sqrt(float64(tubeCount)))
	return n*pitch + od
}
