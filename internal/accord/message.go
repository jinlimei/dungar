package accord

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) handleMessageUpdateEvent(s *discordgo.Session, ev *discordgo.MessageUpdate) {
	if ev == nil || isMessageFromSelf(s, ev.Message) {
		return
	}

	var (
		msg        *core2.IncomingMessage
		parsed     *core2.ParsedMessage
		translated string
	)

	if ev.Content == "" && ev.Embeds != nil && len(ev.Embeds) > 0 {
		// We've got an embedYay
		// Let's see if it's a message already handled
		converted := d.convertEmbeds(ev.Embeds)

		authorID := ""

		if ev.Author != nil {
			authorID = ev.Author.ID
		}

		msg = &core2.IncomingMessage{
			ID:          ev.ID,
			UserID:      authorID,
			ServerID:    ev.GuildID,
			ChannelID:   ev.ChannelID,
			Type:        core2.MessageTypeChanged,
			Attachments: converted,
		}
	} else if ev.Content != "" && (ev.Embeds == nil || len(ev.Embeds) == 0) {
		parsed = parseDiscordMessage(ev.Message)
		translated = d.translateParsedMessage(ev.GuildID, parsed)

		msg = &core2.IncomingMessage{
			ID:            ev.ID,
			UserID:        ev.Author.ID,
			ServerID:      ev.GuildID,
			ChannelID:     ev.ChannelID,
			Contents:      translated,
			LowerContents: strings.ToLower(translated),
			Type:          core2.MessageTypeChanged,
		}
	} else {
		logEvent("message_update_empty", ev.Timestamp, ev)
		return
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
	logEvent("message_delete", ev.Timestamp, ev)

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
