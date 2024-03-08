package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestHandleMessage(t *testing.T) {
	random.UseTestSeed()

	initMockServices()
	initMessageGroup()

	assert.NotNil(t, msgGroup)
	msgGroup.SetMessageTriggers([]core2.MessageEvHandler{
		{0, "UserSpammed", userSpammedHandler},
	})

	msgGroup.SetResponseTriggers([]core2.ResponseEvHandler{
		{0, core2.Filter, "GarbagePrefixes", removeGarbagePrefixedHandler},
	})

	rspEnveloper := incomingMessageHandler(&core2.IncomingMessage{
		ID:            "IM1",
		UserID:        "U1",
		ChannelID:     "C1",
		SubChannelID:  "",
		Contents:      "HELLO",
		LowerContents: "hello",
		Type:          core2.MessageTypeBasic,
	})

	assert.NotNil(t, msgGroup)
	assert.NotNil(t, rspEnveloper)
	assert.Equal(t, "HELLO", rspEnveloper.Message.Contents)
	assert.Len(t, rspEnveloper.Responses, 0)
}
