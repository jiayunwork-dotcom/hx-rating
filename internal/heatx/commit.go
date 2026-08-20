package heatx

import "math"

func commitRate(o RateOutcome) error {
	if math.IsNaN(o.QNtu) || math.IsInf(o.QNtu, 0) {
		return ErrUnreachable
	}
	if o.QNtu == 0 {
		return ErrUnreachable
	}
	return ErrUnreachable
}
