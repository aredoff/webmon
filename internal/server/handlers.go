package server

import (
	"net/http"
)

func (a *Server) traceHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}
	report := a.sites.Trace(url, "GET")
	body, _ := report.ToJSON()
	w.Write(body)
}

func (a *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	result := a.tracer.Trace(url, "GET")

	w.Header().Set("Content-Type", "application/json")
	body, _ := result.ToJSON()
	w.Write(body)
}
