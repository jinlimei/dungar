package core2

// IncomingEventHandler is our handler for non-message IncomingEvent's which _could_ be
// optionally responded to (i.e. return a ResponseEnvelope)
type IncomingEventHandler func(svc *Service, event *IncomingEvent) *EventResponseEnvelope

// HandleIncomingEvent is core's func for working with IncomingEvent's and passing them
// off to the triggers associated
func (s *Service) HandleIncomingEvent(ev *IncomingEvent) *EventResponseEnvelope {
	return s.evHandler(s, ev)
}
