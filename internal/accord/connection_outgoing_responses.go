package accord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
	"time"
)

// GetOutgoingResponses retrieves the OutgoingResponses and ensures that outgoing is
// not nil (to avoid nil dereference shenanigans)
func (d *Driver) GetOutgoingResponses() *core2.OutgoingResponses {
	if d.outgoing == nil {
		d.outgoing = core2.NewOutgoingResponses()
	}

	return d.outgoing
}

func (d *Driver) handleOutgoingResponses(outgoing *core2.OutgoingResponses) {
	envelopes := outgoing.RetrieveEnvelopes()

	for _, envelope := range envelopes {
		d.handleEnvelope(envelope)
	} // end of for _, envelope := range envelopes
}

func (d *Driver) handleEnvelope(envelope *core2.ResponseEnvelope) {
	for _, response := range envelope.Responses {
		if envelope.Message.ChannelID == "" {
			db.LogIssue(
				"invalid_message",
				"Encountered Bad Outgoing Message",
				spew.Sdump(response, envelope.Message),
			)

			time.Sleep(1 * time.Second)
			continue
		}

		if response.IsValidOutgoing(d, envelope.Message) {
			d.handleResponse(envelope, response)

			time.Sleep(1 * time.Second)
		} else {
			db.LogIssue(
				"invalid_message",
				"Attempted to send invalid message",
				spew.Sdump(response, envelope.Message),
			)

			time.Sleep(1 * time.Second)
		} // end of else
	} // end of for _, response := range envelope.Responses
}

func (d *Driver) handleResponse(envelope *core2.ResponseEnvelope, response *core2.Response) {
	var (
		formatted string

		sess *discordgo.Session
		err  error
	)

	formatted = response.Format(d, envelope.Message)

	log.Printf(
		"[outgoingEventLoop] channelID='%s', subChannelID='%s', responseType='%s', message='%s'",
		envelope.Message.ChannelID,
		envelope.Message.SubChannelID,
		response.ResponseType.String(),
		formatted,
	)

	sess = d.Con.GetSession()

	switch response.ResponseType {
	case core2.ResponseTypeBasic, core2.ResponseTypeDirectMessage:

		// We need to conform to discord's rather complex support for emoji
		formatted = d.translateMessageEmoji(formatted, envelope.Message.ServerID)

		if response.MakeSubChannelIfPossible {
			_, err = sess.ChannelMessageSendReply(
				envelope.Message.ChannelID,
				formatted,
				&discordgo.MessageReference{
					MessageID: envelope.Message.ID,
					ChannelID: envelope.Message.ChannelID,
					GuildID:   envelope.Message.ServerID,
				},
			)
		} else {
			_, err = sess.ChannelMessageSend(
				envelope.Message.ChannelID,
				formatted,
			)
		}

		if err != nil {
			db.LogIssue(
				"message_send_fail",
				"Failed to send message",
				fmt.Sprintf(
					"Response:\n%s\n\nEnvelope Message:\n%s\n",
					spew.Sdump(response),
					spew.Sdump(envelope.Message),
				),
			)
		} else {
			d.core.SetLastSentMessage()
		}

	case core2.ResponseTypeReaction:
		formatted = d.translateReactionEmoji(formatted, envelope.Message.ServerID)

		err = sess.MessageReactionAdd(
			envelope.Message.ChannelID,
			envelope.Message.ID,
			formatted,
		)

		log.Printf(
			"MessageReactionAdd(channel-id='%s',message-id='%s',formatted='%s'",
			envelope.Message.ChannelID,
			envelope.Message.ID,
			formatted,
		)

		if err != nil {
			db.LogIssue(
				"reaction_add_failed",
				"Failed to Add Reaction",
				fmt.Sprintf(
					"Reaction: '%s'\nChannelID: '%s'\nSubChannelID: '%s'\nError: %v",
					formatted,
					envelope.Message.ChannelID,
					envelope.Message.SubChannelID,
					err,
				),
			)
		} else {
			d.core.SetLastSentMessage()
		}
	} // end switch response.ResponseType
}
