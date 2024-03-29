package httptracer

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

const (
	maxRedirectsDefault   = 10
	requestTimeoutDefault = 5 * time.Second
)

type HttpTracer interface {
	Trace(url, method string) *TracerResult
	SetTimeout(d time.Duration)
	SetHeaders(key, value string)
}

func New() HttpTracer {
	return &tracer{
		client: &http.Client{
			Timeout: requestTimeoutDefault,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxRedirectsDefault {
					return errors.New("stopped after 10 redirects")
				}
				return nil
			},
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		headers: map[string]string{},
	}
}

type tracer struct {
	client  *http.Client
	headers map[string]string
}

func (t *tracer) SetTimeout(d time.Duration) {
	t.client.Timeout = d
}

func (t *tracer) SetHeaders(key, value string) {
	t.headers[key] = value
}

func (t *tracer) SetMaxRedirects(n int) {
	t.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= n {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
}

func (t *tracer) Trace(siteURL, method string) *TracerResult {
	_, err := url.ParseRequestURI(siteURL)
	if err != nil {
		return &TracerResult{Error: err}
	}

	defer t.client.CloseIdleConnections()

	req, _ := http.NewRequest(strings.ToUpper(method), siteURL, nil)

	var startTime, connectStartTime, nameLookupStartTime, tlsHandshakeStartTime time.Time
	httpStatData := TracerResult{}

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { nameLookupStartTime = time.Now() },
		DNSDone:  func(ddi httptrace.DNSDoneInfo) { httpStatData.NameLookup = time.Since(nameLookupStartTime) },

		TLSHandshakeStart: func() { tlsHandshakeStartTime = time.Now() },
		TLSHandshakeDone:  func(cs tls.ConnectionState, err error) { httpStatData.TLSHandshake = time.Since(tlsHandshakeStartTime) },

		ConnectStart: func(network, addr string) { connectStartTime = time.Now() },
		ConnectDone:  func(network, addr string, err error) { httpStatData.Connect = time.Since(connectStartTime) },

		GotFirstResponseByte: func() { httpStatData.FirstByte = time.Since(startTime) },
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	for key, value := range t.headers {
		req.Header.Add(key, value)
	}

	startTime = time.Now()
	resp, err := t.client.Do(req)
	if err != nil {
		httpStatData.Error = err
		return &httpStatData
	}
	defer resp.Body.Close()

	httpStatData.StatusCode = resp.StatusCode

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		httpStatData.Error = err
		return &httpStatData
	}

	httpStatData.FullResponse = time.Since(startTime)
	httpStatData.BodySize = binary.Size(bodyBytes)

	return &httpStatData
}
