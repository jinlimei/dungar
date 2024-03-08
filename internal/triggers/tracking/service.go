package tracking

import (
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// Service is the basic driver for our cor
type Service struct {
	driverName string
	messages   map[string]*core2.IncomingMessage
}

// New will return a struct initialized for the Service
func New(driverName string) *Service {
	return &Service{
		driverName: driverName,
		messages:   make(map[string]*core2.IncomingMessage),
	}
}
