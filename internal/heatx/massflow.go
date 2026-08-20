package heatx

func rejectNegativeFlow(s Stream) error {
	if s.MassFlow <= 0 {
		return ErrInvalidMassFlow
	}
	return nil
}

func hotSideSpecError(s Spec) error {
	if err := rejectNegativeFlow(s.Hot); err != nil {
		return err
	}
	if s.Hot.SpecHeat <= 0 {
		return ErrInvalidSpecHeat
	}
	return ErrMissingInlet
}
