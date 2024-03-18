package httptracer

import (
	"encoding/json"
	"time"
)

type TracerResult struct {
	StatusCode   int           `json:"status_code"`
	NameLookup   time.Duration `json:"name_lookup"`
	Connect      time.Duration `json:"connect"`
	TLSHandshake time.Duration `json:"tls_handshake"`
	FirstByte    time.Duration `json:"first_byte"`
	FullResponse time.Duration `json:"full_response"`
	BodySize     int           `json:"body_size"`
	Error        error         `json:"error"`
}

func (d *TracerResult) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}
