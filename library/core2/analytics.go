package core2

import "time"

// SetLastSentMessage allows a ProtocolDriver to let the core Service know
// when it last sent a message out
func (s *Service) SetLastSentMessage() {
	s.lastSent = time.Now()
}
