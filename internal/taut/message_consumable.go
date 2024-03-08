package taut

import "github.com/slack-go/slack"

var ignorableMessages = []string{
	"bot_message",
	"channel_archive",
	"channel_join",
	"channel_leave",
	"channel_purpose",
	"channel_topic",
	"channel_unarchive",
	"file_comment",
	"file_mention",
	"file_share",
	"group_archive",
	"group_purpose",
	"group_topic",
	"group_unarchive",
	"message_deleted",
	"message_replied",
}

func isSpecialSubTypeMessage(msg slack.Msg) bool {
	if msg.Type != "message" {
		return false
	}

	return msg.SubType == "channel_topic" || msg.SubType == "channel_purpose"
}

func isConsumableMessage(msg slack.Msg) bool {
	if msg.Type == "message" && msg.SubType == "" && msg.Text == "" && len(msg.Files) > 0 {
		return false
	}

	// If a message has no SubType then it should be available to the user (in general)
	if msg.SubType == "" {
		return true
	}

	consumable := true

	for _, sub := range ignorableMessages {
		if msg.SubType == sub {
			consumable = false
			break
		}
	}

	return consumable
}
