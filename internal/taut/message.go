package taut

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) handleMessageEvent(ev *slack.MessageEvent) {
	// Skip out on messages that come from the bot from the channel
	// We would prefer it if Dungar did not respond to herself.
	if d.isFromBotUser(ev) {
		return
	}

	//logMessageEvent(ev)

	if isSpecialSubTypeMessage(ev.Msg) {
		d.handleSpecialSubTypeMessage(ev.Msg)
		return
	}

	// ignore messages, hidden or non-consumable
	if !isConsumableMessage(ev.Msg) {
		clientMsgID := ev.Msg.ClientMsgID

		if ev.SubMessage != nil {
			clientMsgID = ev.SubMessage.ClientMsgID
		}

		log.Printf(
			"Message skipped %s -- type='%s' sub-type='%s' text-len='%d' hidden='%v', consumable='false'",
			clientMsgID,
			ev.Msg.Type,
			ev.Msg.SubType,
			len(ev.Msg.Text),
			ev.Msg.Hidden,
		)

		return
	}

	msg, err := d.convertSlackMessage(ev)

	if err != nil {
		jsn, _ := json.Marshal(ev)

		db.LogIssue(
			"convert_slack_message",
			"Failed to convert Slack Message",
			fmt.Sprintf(
				"Error: %v\n\nMessage: %s\n",
				err,
				string(jsn),
			),
		)

		return
	}

	log.Printf("Handle message %s (UserID='%s',ChannelID='%s',SubChannelID='%s',Type='%s')\n",
		msg.ID, msg.UserID, msg.ChannelID, msg.SubChannelID, msg.Type)

	log.Printf("Handle Message Contents: %s", msg.Contents)

	envelope := d.core.HandleIncomingMessage(msg)

	if envelope == nil || envelope.Responses == nil || len(envelope.Responses) == 0 {
		return
	}

	log.Printf(
		"Handled incoming message (responses=%d)\n",
		len(envelope.Responses),
	)

	d.GetOutgoingResponses().AddEnvelope(envelope)
}

func (d *Driver) handleSpecialSubTypeMessage(msg slack.Msg) {
	switch msg.SubType {
	case "channel_topic":
		d.core.HandleIncomingEvent(&core2.IncomingEvent{
			ID:        msg.ClientMsgID,
			UserID:    msg.User,
			ChannelID: msg.Channel,
			ServerID:  d.teamID,
			Text:      msg.Text,
			Type:      core2.EventChannelTopic,
		})
	case "channel_purpose":
		d.core.HandleIncomingEvent(&core2.IncomingEvent{
			ID:        msg.ClientMsgID,
			UserID:    msg.User,
			ChannelID: msg.Channel,
			ServerID:  d.teamID,
			Text:      msg.Text,
			Type:      core2.EventChannelPurpose,
		})
	}
}
