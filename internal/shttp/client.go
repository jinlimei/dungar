package shttp

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	tracer "net/http/httptrace"
	"time"
)

// New returns our default HTTP client
func New() *HTTP {
	return &HTTP{}
}

// HTTP is our default client
type HTTP struct {
	transport    *http.Transport
	client       *http.Client
	trace        *tracer.ClientTrace
	traceEnabled bool
}

// EnableTrace sets up tracing enabled for _all_ requests.
func (h *HTTP) EnableTrace() {
	h.traceEnabled = true

	if h.trace != nil {
		return
	}

	h.trace = &tracer.ClientTrace{
		GotConn: func(info tracer.GotConnInfo) {
			log.Printf("[http trace] GotConn: (reused: %v, was-idle: %v, idle-time: %v)\n",
				info.Reused, info.WasIdle, info.IdleTime)
		},
	}
}

// DisableTrace disables the all-request tracing. If SimpleRequest still requests
// tracing, we will still do tracing (this does not override SimpleRequest)
func (h *HTTP) DisableTrace() {
	h.traceEnabled = false
}

// Request is our HTTPRequester implementation
func (h *HTTP) Request(sReq *SimpleRequest) (*SimpleResponse, error) {
	if sReq == nil {
		return nil, ErrNilRequest
	}

	// Build & store our client!
	h.buildClient()

	req, err := http.NewRequest(sReq.Method, sReq.URL, bytes.NewReader(sReq.Body))

	if err != nil {
		return nil, err
	}

	for key, val := range sReq.Headers {
		req.Header.Add(key, val)
	}

	origTrace := h.traceEnabled
	// If we're enabling trace due to the request we want to return it
	// to the original setting.
	if sReq.Trace {
		h.EnableTrace()
		h.traceEnabled = true
	}

	// Setup trace if enabled and with a trace func provided.
	if h.traceEnabled && h.trace != nil {
		req = req.WithContext(tracer.WithClientTrace(req.Context(), h.trace))
	}

	// Return traceEnabled to its original setting
	if sReq.Trace {
		h.traceEnabled = origTrace
	}

	rsp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	err = rsp.Body.Close()
	if err != nil {
		log.Printf("Failed to close body of HTTP request: %v\n", err)
	}

	return &SimpleResponse{
		Rsp:  rsp,
		Code: rsp.StatusCode,
		Body: body,
	}, nil
}

func (h *HTTP) buildClient() {
	if h.client != nil {
		return
	}

	h.transport = &http.Transport{
		MaxIdleConns:        30,
		MaxIdleConnsPerHost: 30,
		MaxConnsPerHost:     30,
		IdleConnTimeout:     30 * time.Second,
		TLSHandshakeTimeout: 300 * time.Second,
	}

	h.client = &http.Client{
		Transport: h.transport,
		Timeout:   30 * time.Second,
	}
}
