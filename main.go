package main

import (
	"fmt"
	"os"

	"hx-rating/internal/cli"
)

const usage = `hx-rating — shell-and-tube heat exchanger rating (NTU & LMTD)

usage:
  hx-rating rate [flags] [file]
    rate an exchanger from a JSON file or stdin (see example/counter.json)

flags:
  -f, --file <path>   read JSON from a file (default: stdin)
  -h, --help          show this help

input JSON fields:
  hot.cold:  { "mass_flow", "specific_heat", "inlet_temp" }
  ua         overall heat transfer conductance W/K (or u + area)
  flow       "counter" | "parallel"
  fouling    { "hot", "cold" } optional fouling resistances
  pressure_drop { "diameter", "length", "roughness", "density" } optional

examples:
  hx-rating rate example/counter.json
  cat example/counter.json | hx-rating rate`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	switch cmd {
	case "rate":
		os.Exit(cli.RunRate(args, os.Stdout, os.Stderr))
	case "-h", "--help", "help":
		fmt.Println(usage)
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s\n", cmd, usage)
		os.Exit(2)
	}
}
