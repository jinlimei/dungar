package shttp

// NewMock returns an HTTPMock with defaults
func NewMock() *HTTPMock {
	return &HTTPMock{}
}

// HTTPMock is our mock implementation of HTTPRequester
type HTTPMock struct {
	rsp     *SimpleResponse
	err     error
	handler func(r *SimpleRequest) (*SimpleResponse, error)
}

// SetResponse will set the return response for Request
func (hm *HTTPMock) SetResponse(rsp *SimpleResponse) {
	hm.rsp = rsp
}

// SetError will set the return error for Request
func (hm *HTTPMock) SetError(err error) {
	hm.err = err
}

// SetHandler bypasses the code result of SetError and SetResponse and lets
// the user mock the entire request handling.
func (hm *HTTPMock) SetHandler(handler func(r *SimpleRequest) (*SimpleResponse, error)) {
	hm.handler = handler
}

// Request is our mocked HTTPRequester interface function.
func (hm *HTTPMock) Request(req *SimpleRequest) (*SimpleResponse, error) {
	if hm.handler != nil {
		return hm.handler(req)
	}

	return hm.rsp, hm.err
}
