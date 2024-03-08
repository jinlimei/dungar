package shttp

import (
	"net/http"

	"errors"
)

// HTTPRequester is a simplified HTTP client interface for working
// with remote connections in an easier way than the weird spaghetti
// I have been doing
type HTTPRequester interface {
	// Request is our single only method needed to implement this interface.
	Request(request *SimpleRequest) (*SimpleResponse, error)
}

var (
	// ErrNilRequest occurs when an HTTPRequester.Request receives a nil input
	ErrNilRequest = errors.New("request provided is nil")
)

// SimpleRequest is our simple handling for a request state that
// will be turned into an http.Request to work with later.
type SimpleRequest struct {
	// Method is the HTTP Method
	Method string
	// URL is well, you should know
	URL string
	// Headers is also pretty straightforward here. HTTP Headers!
	Headers map[string]string
	// Body can be nil or not, and will be used in requests that accept a body.
	Body []byte
	// Trace if specified as true will do some simple httptrace.ClientTrace stuff
	Trace bool
}

// SimpleResponse is our wrapper of a http.Response with the body and code
// provided quicker
type SimpleResponse struct {
	// Rsp is the actual http response
	Rsp *http.Response
	// Code is the http.Response StatusCode
	Code int
	// Body is the body from the response. We _always_ read the body, but this
	// could still be nil
	Body []byte
}
