package tracking

import (
	"errors"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// StoreMessage will store a raw core2.IncomingMessage by its message ID
func (s *Service) StoreMessage(message *core2.IncomingMessage) error {
	s.messages[message.ID] = message
	return nil
}

// RetrieveMessage attempts to retrieve a given message by its ID
func (s *Service) RetrieveMessage(messageID string) (*core2.IncomingMessage, error) {
	msg, ok := s.messages[messageID]

	if !ok {
		return nil, errors.New("message not stored")
	}

	return msg, nil
}
