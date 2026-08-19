package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"hx-rating/internal/heatx"
)

type rateOutcome struct {
	Q             float64 `json:"q"`
	HotOut        float64 `json:"hot_out"`
	ColdOut       float64 `json:"cold_out"`
	Effectiveness float64 `json:"effectiveness"`
	NTU           float64 `json:"ntu"`
	LMTD          float64 `json:"lmtd"`
	Feasible      bool    `json:"feasible"`
	EnergyResidual float64 `json:"energy_residual"`
	NTULMTDDiff   float64 `json:"ntu_lmtd_diff"`
	Crossed       bool    `json:"temperature_cross"`
}

func adaptRate(o heatx.RateOutcome) rateOutcome {
	return rateOutcome{
		Q:             o.QNtu,
		HotOut:        o.HotOut,
		ColdOut:       o.ColdOut,
		Effectiveness: o.Eff,
		NTU:           o.NTU,
		LMTD:          o.LMTD,
		Feasible:      o.Feasible,
		EnergyResidual: o.Residual,
		NTULMTDDiff:   o.RelDiff,
		Crossed:       o.Crossed,
	}
}

func formatRate(o heatx.RateOutcome, format string) string {
	if format == "json" {
		b, _ := json.MarshalIndent(adaptRate(o), "", "  ")
		return string(b)
	}
	return fmt.Sprintf("Q              %.6f W\n", o.QNtu) +
		fmt.Sprintf("hot_out        %.6f C\n", o.HotOut) +
		fmt.Sprintf("cold_out       %.6f C\n", o.ColdOut) +
		fmt.Sprintf("effectiveness  %.6f\n", o.Eff) +
		fmt.Sprintf("NTU            %.6f\n", o.NTU) +
		fmt.Sprintf("LMTD           %.6f K\n", o.LMTD) +
		fmt.Sprintf("feasible       %v\n", o.Feasible) +
		fmt.Sprintf("energy_resid   %.3e\n", o.Residual) +
		fmt.Sprintf("ntu_lmtd_diff  %.3e\n", o.RelDiff)
}

func writeRateJSON(w io.Writer, o heatx.RateOutcome) {
	fmt.Fprintln(w, formatRate(o, "json"))
}
