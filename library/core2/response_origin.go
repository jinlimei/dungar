package core2

// ResponseOrigin is where a response has originated from
type ResponseOrigin uint8

const (
	// OriginUnknown the Response never had its ResponseOrigin set
	OriginUnknown ResponseOrigin = iota
	// OriginMessage is a response from an IncomingMessage
	OriginMessage
	// OriginEvent is a response from an IncomingEvent
	OriginEvent
	// OriginScheduled is a response that was scheduled
	OriginScheduled
)
