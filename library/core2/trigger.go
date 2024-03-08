package core2

import "time"

// ResponseHandlerType is our type for determining what
// action to take after a ResponseEvFunc or ScheduleResponseEvFunc runs
type ResponseHandlerType uint8

const (
	// Filter takes the output of our ResponseEvFunc or ScheduleResponseEvFunc
	// and makes it the new output for the entire service. Any backlog which
	// are cancelled are excluded. Some can also be added, but this operates over
	// the whole.
	Filter ResponseHandlerType = iota

	// Adder takes the output of our ResponseEvFunc or ScheduleResponseEvFunc
	// and appends it to the list of existing outgoing responses. This only does
	// appending, and will not modify the existing backlog/slice
	Adder
)

type (
	// MessageEvFunc is our function for building responses from an IncomingMessage
	MessageEvFunc func(svc *Service, msg *IncomingMessage) []*Response

	// ResponseEvFunc is our function for taking in a slice of Response and
	// outputting another slice of Response. When combined with ResponseHandlerType,
	// we determine what to do with the output Response slice.
	ResponseEvFunc func(svc *Service, m []*Response) []*Response

	// ScheduleEvFunc when setLastRan outputs a slice of ScheduledMessage's or nil.
	ScheduleEvFunc func(svc *Service) []*ScheduledMessage

	// ScheduleResponseEvFunc takes the input ScheduledMessage's and builds an
	// output of them. Similarly to ResponseEvFunc, this is combined with
	// ResponseHandlerType to determine what to do with the result slice
	ScheduleResponseEvFunc func(svc *Service, sm []*ScheduledMessage) []*ScheduledMessage
)

// MessageEvHandler is our structure for a specific message handler. The order
// provides an operating order, the Name for debugging/updating/handling, and
// the MessageEvFunc itself for doing the work
type MessageEvHandler struct {
	Order int
	Name  string
	Func  MessageEvFunc
}

// ResponseEvHandler is our structure for ResponseEvFunc's which needs a
// ResponseHandlerType and order.
type ResponseEvHandler struct {
	Order int
	Type  ResponseHandlerType
	Name  string
	Func  ResponseEvFunc
}

// ScheduleEvHandler is our structure for scheduled messaging, and allows us
// identification on the RunTick and LastRan, as well as the Order
type ScheduleEvHandler struct {
	Order   int
	Name    string
	Func    ScheduleEvFunc
	RunTick time.Duration
	LastRan int64
}

func (seh *ScheduleEvHandler) setLastRan() {
	seh.LastRan = time.Now().Unix()
}

// ScheduleResponseEvHandler is similar to ResponseEvHandler but for the output
// of a ScheduleEvHandler.
type ScheduleResponseEvHandler struct {
	Order int
	Type  ResponseHandlerType
	Name  string
	Func  ScheduleResponseEvFunc
}
