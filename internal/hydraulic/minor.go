package hydraulic

type Nozzle struct {
	Diameter float64
	Count    int
}

func NozzleVelocity(flow float64, n Nozzle) float64 {
	if n.Count <= 0 || n.Diameter <= 0 {
		return 0
	}
	area := CrossSection(n.Diameter) * float64(n.Count)
	if area <= 0 {
		return 0
	}
	return flow / area
}

func EntranceLoss(vel float64) float64 {
	return DynamicHead(vel, WaterDensity) * 0.5
}

func ExitLoss(vel float64) float64 {
	return DynamicHead(vel, WaterDensity)
}

func ElbowLoss(vel float64, elbows int) float64 {
	return DynamicHead(vel, WaterDensity) * float64(elbows) * 0.9
}

func TotalMinorLoss(vel float64, elbows, bends int) float64 {
	return EntranceLoss(vel) + ExitLoss(vel) + ElbowLoss(vel, elbows+bends)
}

func SumPressure(r ...Result) float64 {
	sum := 0.0
	for _, v := range r {
		sum += v.DeltaP
	}
	return sum
}
