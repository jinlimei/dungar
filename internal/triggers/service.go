package triggers

import (
	tracker "gitlab.int.magneato.site/dungar/prototype/internal/triggers/tracking"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	core     *core2.Service
	tracking *tracker.Service
)

// RegisterHandlers will, with a valid core2.Service, set the given core2.IncomingEventHandler,
// core2.IncomingMessageHandler, and core2.ScheduleHandler
func RegisterHandlers(svc *core2.Service) {
	tracking = tracker.New(svc.DriverName())
	core = svc
	core.SetIncomingEventHandler(eventHandler)
	core.SetIncomingMessageHandler(incomingMessageHandler)
	core.SetScheduleHandler(scheduleHandler)
}
