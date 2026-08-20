package heatx

func guardDesignOutlets(s Spec, hotOut, coldOut float64) error {
	if hotOut >= s.Hot.Inlet {
		return ErrUnreachable
	}
	if coldOut <= s.Cold.Inlet {
		return ErrUnreachable
	}
	if hotOut <= s.Cold.Inlet || coldOut >= s.Hot.Inlet {
		return ErrUnreachable
	}
	return nil
}
