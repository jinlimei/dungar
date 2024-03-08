package accord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) convertEmbeds(embeds []*discordgo.MessageEmbed) []core2.IncomingAttachment {
	converted := make([]core2.IncomingAttachment, 0)

	for _, embed := range embeds {
		var (
			imageURL    = ""
			serviceName = ""
			attachType  core2.AttachmentType
		)

		if embed.Thumbnail != nil {
			imageURL = embed.Thumbnail.URL
		}

		if embed.Provider != nil {
			serviceName = embed.Provider.Name
		}

		switch serviceName {
		case "twitter":
			attachType = core2.AttachmentTweet
		case "youtube":
			attachType = core2.AttachmentLinkUnfurl
		default:
			attachType = core2.AttachmentUnknown
		}

		authorName := ""
		authorURL := ""

		if embed.Author != nil {
			authorName = embed.Author.Name
			authorURL = embed.Author.URL
		}

		converted = append(converted, core2.IncomingAttachment{
			AuthorName:  authorName,
			AuthorURL:   authorURL,
			Title:       embed.Title,
			TitleLink:   embed.URL,
			ImageURL:    imageURL,
			ServiceName: serviceName,
			Fallback:    "",
			Type:        attachType,
		})
	}

	return converted
}

func (d *Driver) convertMessageCreate(ev *discordgo.Message) (*core2.IncomingMessage, error) {
	parsed := parseDiscordMessage(ev)
	converted := d.translateParsedMessage(ev.GuildID, parsed)

	msg := &core2.IncomingMessage{
		ID:             ev.ID,
		UserID:         ev.Author.ID,
		ChannelID:      ev.ChannelID,
		ServerID:       ev.GuildID,
		Contents:       converted,
		LowerContents:  strings.ToLower(converted),
		ParsedContents: parsed,
	}

	channel, ok := d.guilds.getChannelByID(ev.ChannelID)

	if !ok {
		return nil, core2.ErrChannelNotFound
	}

	// We're a thread! A THREAD!
	if channel.IsThread() {
		msg.IsSubMessage = true
		msg.SubChannelID = channel.ID
		msg.ChannelID = channel.ParentID
	}

	if ev.MessageReference != nil {
		msg.IsReply = true
		msg.ReplyToID = ev.MessageReference.MessageID
	}

	msg.Type = core2.MessageTypeBasic

	return msg, nil
}
