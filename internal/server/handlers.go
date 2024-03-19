package server

import (
	"encoding/json"
	"net/http"
)

func (a *Server) traceHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	treport := a.sites.Trace(url, "GET")

	reportFor := traceReport{
		StatusCode:       treport.StatusCode,
		NameLookupTime:   treport.NameLookup.Seconds(),
		ConnectTime:      treport.Connect.Seconds(),
		TLSHandshakeTime: treport.TLSHandshake.Seconds(),
		TransportTime:    (treport.FullResponse - treport.FirstByte).Seconds(),
		FullTime:         treport.FullResponse.Seconds(),
		BodySize:         treport.BodySize,
	}
	switch treport.Error {
	case nil:
		reportFor.Error = 0
	default:
		reportFor.Error = 1
	}

	body, err := reportFor.ToJSON()
	if err != nil {
		http.Error(w, "Json convert report error", http.StatusBadRequest)
		return
	}
	w.Write(body)
}

func (a *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	result := map[string]string{
		"status": "ok",
	}
	js, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
