package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"hx-rating/internal/heatx"
)

type InputFile struct {
	Hot   streamSpec `json:"hot"`
	Cold  streamSpec `json:"cold"`
	UA    float64    `json:"ua"`
	Area  float64    `json:"area"`
	U     float64    `json:"u"`
	Flow  string     `json:"flow"`
	Fouling foulingSpec `json:"fouling"`
	Drop  *dropSpec  `json:"pressure_drop"`
}

type streamSpec struct {
	MassFlow  float64 `json:"mass_flow"`
	SpecHeat  float64 `json:"specific_heat"`
	InletTemp float64 `json:"inlet_temp"`
}

type foulingSpec struct {
	Hot  float64 `json:"hot"`
	Cold float64 `json:"cold"`
}

type dropSpec struct {
	Diameter  float64 `json:"diameter"`
	Length    float64 `json:"length"`
	Roughness float64 `json:"roughness"`
	Density   float64 `json:"density"`
}

var ErrEmptyInput = errors.New("cli: empty input")
var ErrInvalidJSON = errors.New("cli: invalid json")

var stdinOverride io.Reader

func ParseInput(r io.Reader) (InputFile, error) {
	var in InputFile
	dec := json.NewDecoder(r)
	if err := dec.Decode(&in); err != nil {
		if err == io.EOF {
			return in, ErrEmptyInput
		}
		return in, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}
	return in, nil
}

func ReadFile(path string) (InputFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return InputFile{}, err
	}
	defer f.Close()
	return ParseInput(f)
}

func ReadStdin() (InputFile, error) {
	if stdinOverride != nil {
		return ParseInput(stdinOverride)
	}
	return ParseInput(os.Stdin)
}

func (in InputFile) ToSpec() heatx.Spec {
	ua := in.UA
	if ua <= 0 && in.U > 0 && in.Area > 0 {
		ua = in.U * in.Area
	}
	rf := heatx.Fouling{Hot: in.Fouling.Hot, Cold: in.Fouling.Cold}
	flow := heatx.Counter
	if in.Flow == "parallel" {
		flow = heatx.Parallel
	}
	s := heatx.Spec{
		Hot: heatx.Stream{
			MassFlow: in.Hot.MassFlow,
			SpecHeat: in.Hot.SpecHeat,
			Inlet:    in.Hot.InletTemp,
		},
		Cold: heatx.Stream{
			MassFlow: in.Cold.MassFlow,
			SpecHeat: in.Cold.SpecHeat,
			Inlet:    in.Cold.InletTemp,
		},
		UA:   ua,
		Area: in.Area,
		Flow: flow,
		Rf:   rf,
	}
	return s
}

func (in InputFile) Validate() error {
	if in.Hot.MassFlow == 0 && in.Hot.SpecHeat == 0 && in.Hot.InletTemp == 0 &&
		in.Cold.MassFlow == 0 && in.Cold.SpecHeat == 0 && in.Cold.InletTemp == 0 {
		return ErrEmptyInput
	}
	return nil
}
