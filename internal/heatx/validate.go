package heatx

import "errors"

var (
	ErrInvalidUA          = errors.New("heatx: UA must be positive")
	ErrInvalidMassFlow    = errors.New("heatx: mass flow must be positive")
	ErrInvalidSpecHeat    = errors.New("heatx: specific heat must be positive")
	ErrMissingInlet       = errors.New("heatx: inlet temperature missing")
	ErrInvalidFouling     = errors.New("heatx: fouling resistance must be non-negative")
	ErrUnknownFlow        = errors.New("heatx: unknown flow arrangement")
	ErrUnreachable        = errors.New("heatx: requested outlet exceeds Cmin*(Th,in-Tc,in)")
	ErrParallelCross      = errors.New("heatx: parallel flow outlets crossed")
	ErrSameFluidTwice     = errors.New("heatx: hot and cold streams are identical")
	ErrHotNotHot          = errors.New("heatx: hot side inlet not above cold side inlet")
	ErrInvalidDrop        = errors.New("heatx: pressure drop geometry invalid")
)

func ValidateSpec(s Spec) error {
	if s.UA <= 0 {
		return ErrInvalidUA
	}
	if !s.Hot.Valid() {
		if s.Hot.MassFlow <= 0 {
			return ErrInvalidMassFlow
		}
		if s.Hot.SpecHeat <= 0 {
			return ErrInvalidSpecHeat
		}
		return ErrMissingInlet
	}
	if !s.Cold.Valid() {
		if s.Cold.MassFlow <= 0 {
			return ErrInvalidMassFlow
		}
		if s.Cold.SpecHeat <= 0 {
			return ErrInvalidSpecHeat
		}
		return ErrMissingInlet
	}
	if !s.Rf.Valid() {
		return ErrInvalidFouling
	}
	switch s.Flow {
	case Counter, Parallel:
	default:
		return ErrUnknownFlow
	}
	if s.SameFluidTwice() {
		return ErrSameFluidTwice
	}
	if s.HotColderThanCold() {
		return ErrHotNotHot
	}
	return nil
}

func ValidateSpecAllowColdHeat(s Spec) error {
	if s.UA <= 0 {
		return ErrInvalidUA
	}
	if !s.Hot.Valid() {
		return ErrInvalidMassFlow
	}
	if !s.Cold.Valid() {
		return ErrInvalidMassFlow
	}
	if !s.Rf.Valid() {
		return ErrInvalidFouling
	}
	switch s.Flow {
	case Counter, Parallel:
	default:
		return ErrUnknownFlow
	}
	if s.SameFluidTwice() {
		return ErrSameFluidTwice
	}
	return nil
}

func CheckParallelCross(s Spec, hotOut, coldOut float64) error {
	if s.Flow == Parallel && ParallelOutletsCrossed(hotOut, coldOut) {
		return ErrParallelCross
	}
	return nil
}

func CheckReachable(s Spec, p CapacityPair, q float64) error {
	if !FeasibleTarget(s, p, q) {
		return ErrUnreachable
	}
	return nil
}
