package core2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistory(t *testing.T) {
	assert.True(t, true)

	driver := &MockProtocolDriver{}
	svc := New(driver)

	svc.SetIncomingMessageHandler(func(msg *IncomingMessage) *ResponseEnvelope {
		return &ResponseEnvelope{
			Message:   msg,
			Responses: EmptyRsp(),
		}
	})

	assert.True(t, len(svc.GetPreviousMessages(1)) == 0)

	svc.HandleIncomingMessage(makeIncomingMessage(
		"aaaaAAAAAA",
		MessageTypeBasic,
	))

	assert.True(t, len(svc.GetPreviousMessages(1)) == 1)
}
