package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func rawMessageRecorder(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	switch svc.DriverName() {
	case "slack":
		db.LegacyRecordRawMessage(msg.Contents, "slack")
	case "discord":
		db.RecordRawDiscordMessage(
			msg.ID,
			msg.ServerID,
			msg.ChannelID,
			msg.Contents,
		)
	}

	return core2.EmptyRsp()
}
