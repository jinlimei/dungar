package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// MessageHandler to make sure our bot _only_ responds in the testing channel.
func testingHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	channel, err := svc.GetChannel(msg.ChannelID, msg.ServerID)

	if err != nil {
		return core2.EmptyRsp()
	}

	if channel.Archetype != "testing" {
		return core2.CancelledRsp()
	}

	return core2.EmptyRsp()
}
