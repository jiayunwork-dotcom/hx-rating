# hx-rating

Shell-and-tube heat exchanger rating tool. Given mass flows, specific heats and
inlet temperatures on both sides, plus an overall conductance UA (or U and area),
it computes outlet temperatures, heat transferred, effectiveness, NTU and LMTD for
counter-current and parallel (co-current) flow, and reports whether the duty is
feasible and whether a temperature cross occurs.

Two methods are used and cross-checked: the NTU-effectiveness closed-form
solution and an LMTD calculation over the same inlets/outlets; their relative
difference is reported. Energy is closed: heat released by the hot stream equals
heat gained by the cold stream within tolerance (residual reported). Optional
fouling resistances degrade the clean UA (`1/U = 1/U_clean + R_f`). An optional
tube-side pressure drop is estimated from velocity and Darcy friction.

## Build

```bash
go build .
go test ./...
```

## Usage

```bash
go run . rate example/counter.json
cat example/counter.json | go run . rate
go run . rate -f example/dirty.json
```

Input is JSON on stdin or a file:

```json
{
  "hot":  { "mass_flow": 2.5, "specific_heat": 4186, "inlet_temp": 90 },
  "cold": { "mass_flow": 2.6, "specific_heat": 4186, "inlet_temp": 30 },
  "ua": 8000,
  "flow": "counter",
  "fouling": { "hot": 0.00002, "cold": 0.00001 },
  "pressure_drop": { "diameter": 0.02, "length": 6.0, "roughness": 0.00015, "density": 998 }
}
```

`ua` may be replaced by `u` plus `area`. `flow` is `"counter"` or `"parallel"`.

Output lists Q, both outlet temperatures, effectiveness, NTU, LMTD, feasibility,
the energy residual and the NTU-vs-LMTD relative difference. Invalid input (zero
UA, non-positive flow or specific heat, missing inlet temperature, identical
streams, parallel outlets crossing, duty above the `Cmin*(Th,in-Tc,in)` limit)
is rejected with a message on stderr and a non-zero exit code.

## Limits

Duty can never exceed `Cmin * (Th,in - Tc,in)`; a target outlet beyond that is
reported infeasible rather than computed. Parallel flow with cold outlet above
hot outlet is rejected. Fouling resistances must be non-negative.
