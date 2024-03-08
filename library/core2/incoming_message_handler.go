package core2

import (
	"log"
	"time"
)

// IncomingMessageHandler is our standardized signature for message handlers
type IncomingMessageHandler func(msg *IncomingMessage) *ResponseEnvelope

// HandleIncomingMessage allows a ProtocolDriver to signal a message event to Service.
// This will also interact with the message queue for specific backlog.
func (s *Service) HandleIncomingMessage(msg *IncomingMessage) *ResponseEnvelope {
	s.lastRecv = time.Now()

	switch msg.Type {
	case MessageTypeBasic, MessageTypeMe, MessageTypeBroadcast:
		// Queue message!
		s.backlog = append(s.backlog, msg)
		// Pop the last message out
		if len(s.backlog) >= 100 {
			s.backlog = s.backlog[1:]
		}

	case MessageTypeChanged, MessageTypeDeleted:
		// A previous message has been changed or deleted, so we're going to go into the
		// queue and update that message.
		for k := len(s.backlog) - 1; k >= 0; k-- {
			if s.backlog[k].ID == msg.ID {
				s.backlog[k] = msg
				break
			}
		}

	case MessageTypeUnknown:
		log.Printf("WARN: Unknown Message Type Incoming: %+v", msg)
	}

	return s.incMsgHandler(msg)
}
