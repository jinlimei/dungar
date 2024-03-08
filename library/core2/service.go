package core2

import "time"

// Service is the wrapper for core2's architecture
type Service struct {
	driver        ProtocolDriver
	incMsgHandler IncomingMessageHandler
	evHandler     IncomingEventHandler
	scheduler     ScheduleHandler

	backlog  []*IncomingMessage
	lastSent time.Time
	lastRecv time.Time
}

// New provides a Service with a valid ProtocolDriver and setup for Service.backlog
// This will also execute ProtocolDriver.SetCore on the incoming driver
func New(driver ProtocolDriver) *Service {
	svc := &Service{
		driver:   driver,
		backlog:  make([]*IncomingMessage, 0, 100),
		lastSent: time.Unix(0, 0),
		lastRecv: time.Unix(0, 0),
	}

	driver.SetCore(svc)

	return svc
}

// SetIncomingMessageHandler will take the given IncomingMessageHandler and set it
func (s *Service) SetIncomingMessageHandler(incMsgHandler IncomingMessageHandler) {
	s.incMsgHandler = incMsgHandler
}

// SetIncomingEventHandler will register the IncomingEventHandler
func (s *Service) SetIncomingEventHandler(evHandler IncomingEventHandler) {
	s.evHandler = evHandler
}

// SetScheduleHandler will register the ScheduleHandler
func (s *Service) SetScheduleHandler(scheduler ScheduleHandler) {
	s.scheduler = scheduler
}

// Driver returns the specific ProtocolDriver of this instance
func (s *Service) Driver() ProtocolDriver {
	return s.driver
}

// DriverName returns the expected driver name for the set ProtocolDriver
func (s *Service) DriverName() string {
	return s.driver.DriverName()
}
