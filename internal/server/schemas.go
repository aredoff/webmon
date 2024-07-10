package server

import (
	"encoding/json"
	"fmt"
	"slices"
)

type traceReport struct {
	StatusCode       int     `json:"status_code"`
	NameLookupTime   float64 `json:"time_lookup"`
	ConnectTime      float64 `json:"time_connect"`
	TLSHandshakeTime float64 `json:"time_tls_handshake"`
	TransportTime    float64 `json:"time_transport"`
	FullTime         float64 `json:"time_full"`
	BodySize         int     `json:"body_size"`
	Error            int     `json:"error"`

	WarningCode int    `json:"warning_code"`
	WarningText string `json:"warning_text"`
}

func (d *traceReport) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (d *traceReport) AnalyzeWarning(excodes []int) {
	d.WarningCode = 0
	d.WarningText = ""
	if d.Error == 1 {
		d.WarningCode = 1
		d.WarningText = "Load site error"
		return
	}

	if d.StatusCode >= 400 {
		if !slices.Contains[[]int, int](excodes, d.StatusCode) {
			d.WarningCode = 1
			d.WarningText = fmt.Sprintf("Status code: %d", d.StatusCode)
		}
	}
}
