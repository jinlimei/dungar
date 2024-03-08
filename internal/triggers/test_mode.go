package triggers

import "gitlab.int.magneato.site/dungar/prototype/library/core2"

func testModeHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	channel, _ := svc.GetChannel(msg.ChannelID, msg.ServerID)

	if channel.Name == "dungar-test" {
		return core2.EmptyRsp()
	}

	return core2.CancelledRsp()
}
