package accord

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) handleMessageUpdateEvent(s *discordgo.Session, ev *discordgo.MessageUpdate) {
	if ev == nil || ev.Message == nil || isMessageFromSelf(s, ev.Message) {
		return
	}

	if !isConsumableDiscordUpdateMessage(ev.Message) {
		if db.DeleteRawDiscordMessage(ev.ID) {
			db.NotifyRawMessageMutation("discord-message-update-delete")
		}

		logEvent("message_update_empty", ev.Timestamp, ev)
		return
	}

	var (
		msg *core2.IncomingMessage
		err error
	)
	if ev.Author != nil {
		msg, err = d.convertMessageCreate(ev.Message)
		if err != nil {
			jsn, _ := json.Marshal(ev)

			db.LogIssue(
				"convert_discord_message_update",
				"Failed to convert discord message update",
				fmt.Sprintf("Error: %v\n\nMessage: %s\n", err, string(jsn)),
			)

			return
		}
	} else {
		parsed := parseDiscordMessage(ev.Message)
		translated := d.translateParsedMessage(ev.GuildID, parsed)

		msg = &core2.IncomingMessage{
			ID:             ev.ID,
			UserID:         "",
			ServerID:       ev.GuildID,
			ChannelID:      ev.ChannelID,
			Contents:       translated,
			LowerContents:  strings.ToLower(translated),
			ParsedContents: parsed,
		}

		channel, ok := d.guilds.getChannelByID(ev.ChannelID)
		if ok && channel.IsThread() {
			msg.IsSubMessage = true
			msg.SubChannelID = channel.ID
			msg.ChannelID = channel.ParentID
		}

		if ev.MessageReference != nil {
			msg.IsReply = true
			msg.ReplyToID = ev.MessageReference.MessageID
		}
	}

	msg.Type = core2.MessageTypeChanged

	if shouldSyncDiscordUpdateMessage(msg) {
		if db.SyncRawDiscordMessage(
			msg.ID,
			msg.ServerID,
			msg.ChannelID,
			msg.UserID,
			msg.Contents,
			ev.Timestamp.UTC(),
		) {
			db.NotifyRawMessageMutation("discord-message-update")
		}
	} else if db.DeleteRawDiscordMessage(msg.ID) {
		db.NotifyRawMessageMutation("discord-message-update-filtered")
	}

	d.core.HandleIncomingMessage(msg)
	logEvent("message_update", ev.Timestamp, ev)
}

func (d *Driver) handleMessageCreateEvent(s *discordgo.Session, ev *discordgo.MessageCreate) {
	logEvent("message_create", ev.Timestamp, ev)

	// The author is us, so it's one of our messages.
	if isMessageFromSelf(s, ev.Message) {
		return
	}

	if !isConsumableMessage(ev) {
		log.Printf(
			"Message skipped %s -- type='%d' -- contents='%s'",
			ev.ID,
			ev.Type,
			ev.Content,
		)

		return
	}

	msg, err := d.convertMessageCreate(ev.Message)

	if err != nil {
		jsn, _ := json.Marshal(ev)

		db.LogIssue(
			"convert_discord_message",
			"Failed to convert discord message",
			fmt.Sprintf("Error: %v\n\nMessage: %s\n", err, string(jsn)),
		)

		return
	}

	envelope := d.core.HandleIncomingMessage(msg)

	if envelope == nil || envelope.Responses == nil || len(envelope.Responses) == 0 {
		return
	}

	d.GetOutgoingResponses().AddEnvelope(envelope)
}

func (d *Driver) handleMessageDeleteEvent(_ *discordgo.Session, ev *discordgo.MessageDelete) {
	if ev == nil {
		return
	}

	logEvent("message_delete", ev.Timestamp, ev)

	if db.DeleteRawDiscordMessage(ev.ID) {
		db.NotifyRawMessageMutation("discord-message-delete")
	}

	var (
		authorID string
	)

	if ev.Author != nil {
		authorID = ev.Author.ID
	}

	d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.ID,
		UserID:    authorID,
		ChannelID: ev.ChannelID,
		Text:      ev.Content,
		Type:      core2.EventMessageDelete,
	})
}

func shouldSyncDiscordUpdateMessage(msg *core2.IncomingMessage) bool {
	if msg == nil {
		return false
	}

	if strings.TrimSpace(msg.Contents) == "" {
		return false
	}

	if !utils.IsDiscordLearningChannelAllowed(msg.ChannelID) {
		return false
	}

	if db.IsRawMessageUserBlocked("discord", msg.ServerID, msg.UserID) {
		return false
	}

	return true
}

func isConsumableDiscordUpdateMessage(msg *discordgo.Message) bool {
	if msg == nil {
		return false
	}

	if msg.WebhookID != "" {
		return false
	}

	if strings.TrimSpace(msg.Content) == "" {
		return false
	}

	switch msg.Type {
	case discordgo.MessageTypeDefault, discordgo.MessageTypeReply, discordgo.MessageTypeChatInputCommand:
		return true
	default:
		return false
	}
}
