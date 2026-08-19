package heatx

import "math"

func Kelvin(t float64) float64 {
	return t + 273.15
}

func Celsius(k float64) float64 {
	return k - 273.15
}

func WattToKW(w float64) float64 {
	return w / 1000
}

func PaToKPa(pa float64) float64 {
	return pa / 1000
}

func DegToRad(d float64) float64 {
	return d * math.Pi / 180
}

func RadToDeg(r float64) float64 {
	return r * 180 / math.Pi
}

func DensityFromTemp(t float64) float64 {
	return 1000 - 0.2*(t-20)
}

func CpWater(t float64) float64 {
	return 4186 + 0.9*(t-20)
}
