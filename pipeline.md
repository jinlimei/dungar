# Messaging Pipeline


1. Messaging Driver Handles an incoming event, and converts it to `*core2.IncomingMessage`

2. Messaging Driver calls `core2.Service.HandleIncomingMessage(*core2.IncomingMessage)`

3. Core Service then does some internal necessary things and then runs `incMsgHandler` which is assigned via `core2.Service.SetIncomingMessageHandler`

4. At time of writing, the IncomingMessageHandler which ~ does things ~ is in the triggers service:

```go
func (s *Service) RegisterHandlers() {
	s.core.SetIncomingEventHandler(s.eventHandler)
	s.core.SetIncomingMessageHandler(s.incomingMessageHandler)
	s.core.SetScheduleHandler(s.scheduleHandler)
}
```

