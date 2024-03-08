package taut

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) convertSlackMessage(ev *slack.MessageEvent) (*core2.IncomingMessage, error) {
	user, err := d.handleMessageUser(ev)

	if err != nil {
		return nil, err
	}

	var (
		parsed    = parseSlackMessage(ev.Msg.Text)
		converted = d.translateParsedMessage(parsed)
	)

	msg := &core2.IncomingMessage{
		ID:             ev.Msg.ClientMsgID,
		UserID:         user.ID,
		ChannelID:      ev.Msg.Channel,
		ServerID:       d.teamID,
		Contents:       converted,
		LowerContents:  strings.ToLower(converted),
		ParsedContents: parsed,
		Attachments:    d.convertSlackAttachments(ev.Attachments),
	}

	if ev.Msg.ThreadTimestamp != "" {
		msg.SubChannelID = ev.Msg.ThreadTimestamp
		msg.IsSubMessage = true
	} else {
		msg.SubChannelID = ev.Msg.Timestamp
		msg.IsSubMessage = false
	}

	switch ev.Msg.Type {
	case "message":
		msg.Type = core2.MessageTypeBasic
	}

	switch ev.Msg.SubType {
	case "message_changed":
		msg.Type = core2.MessageTypeChanged

		parsed = parseSlackMessage(ev.SubMessage.Text)
		converted = d.translateParsedMessage(parsed)

		msg.ID = ev.SubMessage.ClientMsgID
		msg.Contents = converted
		msg.LowerContents = strings.ToLower(converted)
		msg.ParsedContents = parsed
		msg.Attachments = d.convertSlackAttachments(ev.SubMessage.Attachments)
	}

	return msg, nil
}

func (d *Driver) convertSlackAttachments(attachments []slack.Attachment) []core2.IncomingAttachment {
	out := make([]core2.IncomingAttachment, 0)

	for _, attach := range attachments {
		var (
			attachType  core2.AttachmentType
			fallback    = attach.Fallback
			imageURL, _ = utils.CoalesceStr(attach.ImageURL, attach.ThumbURL)
		)

		switch strings.ToLower(attach.ServiceName) {
		case "twitter":
			attachType = core2.AttachmentTweet
			fallback = attach.Text + "\n - " + attach.AuthorName
			imageURL = attach.AuthorIcon
		case "youtube", "reddit":
			attachType = core2.AttachmentLinkUnfurl
			fallback = attach.Title
		case "":
			// No ServiceName could mean it's an image (or otherwise) url link
			if attach.Title == "" && attach.AuthorName == "" && attach.Text == "" && attach.ImageURL == attach.FromURL {
				attachType = core2.AttachmentImage
				fallback = ""
			}

		default:
			attachType = core2.AttachmentLinkUnfurl

			if attach.Fallback == fmt.Sprintf("%s: %s", attach.ServiceName, attach.Title) {
				fallback = attach.Text
			}

		}

		out = append(out, core2.IncomingAttachment{
			AuthorName:  attach.AuthorName,
			AuthorURL:   attach.AuthorLink,
			Title:       attach.Title,
			TitleLink:   attach.TitleLink,
			ImageURL:    imageURL,
			ServiceName: attach.ServiceName,
			Fallback:    fallback,
			Type:        attachType,
		})
	}

	return out
}
