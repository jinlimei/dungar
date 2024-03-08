package core2

// SubMessageEvFunc is for MessageEvFunc's which do not need a Service parameter
// and instead just need the IncomingMessage parameter
type SubMessageEvFunc func(msg *IncomingMessage) []*Response

// WrapNoServiceFunc will take an incoming SubMessageEvFunc and translate it to a
// MessageEvFunc via a callback
func WrapNoServiceFunc(fn SubMessageEvFunc) MessageEvFunc {
	return func(_ *Service, msg *IncomingMessage) []*Response {
		return fn(msg)
	}
}
