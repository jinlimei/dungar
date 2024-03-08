package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestTestingHandler(t *testing.T) {
	svc := initMockServices()

	screamingInDungarTest := &core2.IncomingMessage{
		ChannelID: "dungar-test",
		Contents:  "AAAAAAAAAAAAA",
	}

	randomCommandInDungarTest := &core2.IncomingMessage{
		ChannelID: "dungar-test",
		Contents:  "!random",
	}

	screamingInAnotherChannel := &core2.IncomingMessage{
		ChannelID: "butts",
		Contents:  "BBBBBBBBBBBB",
	}

	rsp := testingHandler(svc, screamingInDungarTest)
	assert.NotNil(t, rsp)
	assert.Len(t, rsp, 1)
	assert.False(t, rsp[0].ConsumedMessage, "ConsumedMessage on Response should be false for `screamingInDungarTest`")
	assert.False(t, rsp[0].HandledMessage, "HandledMessage on Response should be false for `screamingInDungarTest`")

	rsp = testingHandler(svc, randomCommandInDungarTest)
	assert.NotNil(t, rsp)
	assert.Len(t, rsp, 1)
	assert.False(t, rsp[0].ConsumedMessage)

	rsp = testingHandler(svc, screamingInAnotherChannel)
	assert.NotNil(t, rsp)
	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].IsCancelled(), "IsCancelled() on Response should be true for `screamingInAnotherChannel`")
	assert.True(t, rsp[0].IsHandled(), "IsHandled() on Response should be true for `screamingInAnotherChannel`")
}
