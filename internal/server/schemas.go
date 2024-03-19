package server

import (
	"encoding/json"
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
}

func (d *traceReport) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}
