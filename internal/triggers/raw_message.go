package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func rawMessageRecorder(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	switch svc.DriverName() {
	case "slack":
		db.LegacyRecordRawMessage(msg.Contents, "slack")
	case "discord":
		if !shouldRecordDiscordRawMessage(msg.ChannelID) {
			return core2.EmptyRsp()
		}

		if db.IsRawMessageUserBlocked("discord", msg.ServerID, msg.UserID) {
			return core2.EmptyRsp()
		}

		db.RecordRawDiscordMessage(
			msg.ID,
			msg.ServerID,
			msg.ChannelID,
			msg.UserID,
			msg.Contents,
		)
	}

	return core2.EmptyRsp()
}

func shouldRecordDiscordRawMessage(channelID string) bool {
	return utils.IsDiscordLearningChannelAllowed(channelID)
}
