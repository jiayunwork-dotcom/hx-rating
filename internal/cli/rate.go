package cli

import (
	"fmt"
	"io"
	"math"

	"hx-rating/internal/heatx"
	"hx-rating/internal/hydraulic"
)

type RateOptions struct {
	Path   string
	Format string
}

func RunRate(args []string, stdout, stderr io.Writer) int {
	var path, format string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--file":
			if i+1 < len(args) {
				i++
				path = args[i]
			}
		case "-o", "--output":
			if i+1 < len(args) {
				i++
				format = args[i]
			}
		case "-h", "--help":
			fmt.Fprintln(stdout, rateUsage)
			return 0
		default:
			if path == "" {
				path = args[i]
			}
		}
	}

	var in InputFile
	var err error
	if path != "" {
		in, err = ReadFile(path)
	} else {
		in, err = ReadStdin()
	}
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}
	if err := in.Validate(); err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	spec := in.ToSpec()
	outcome, err := heatx.Rate(spec)
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, formatRate(outcome, format))
	if in.Drop != nil && in.Drop.Diameter > 0 {
		writeDrop(stdout, in)
	}
	return 0
}

const rateUsage = `usage: hx-rating rate [-f <file>] [-o json|text]

rate a heat exchanger from JSON on stdin or from a file.
fields: hot{cold}, cold, ua (or u+area), flow=counter|parallel, fouling{hot,cold}
prints Q, outlet temperatures, effectiveness, NTU, LMTD and feasibility.`

func writeDrop(w io.Writer, in InputFile) {
	drop := in.Drop
	flow := in.Hot.MassFlow / 998
	res := hydraulic.PressureDrop(flow, drop.Diameter, drop.Length, drop.Roughness, 998)
	fmt.Fprintf(w, "tube_velocity  %.4f m/s\n", res.Velocity)
	fmt.Fprintf(w, "tube_deltaP    %.3f Pa\n", res.DeltaP)
	if math.IsNaN(res.DeltaP) || math.IsInf(res.DeltaP, 0) {
		fmt.Fprintf(w, "tube_deltaP    0.000 Pa\n")
	}
}
