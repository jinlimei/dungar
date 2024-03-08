package taut

import (
	"fmt"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
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
				"Encountered Bad Message",
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
	formatted := response.Format(d, envelope.Message)

	log.Printf(
		"[outgoingEventLoop] channelID='%s', subChannelID='%s', responseType='%s', message='%s'\n",
		envelope.Message.ChannelID,
		envelope.Message.SubChannelID,
		response.ResponseType.String(),
		formatted,
	)

	switch response.ResponseType {
	case core2.ResponseTypeBasic:
		d.Con.GetRtm().SendMessage(
			d.Con.GetRtm().NewTypingMessage(envelope.Message.ChannelID),
		)

		time.Sleep(650 * time.Millisecond)

		var outgoing *slack.OutgoingMessage

		if envelope.Message.IsSubMessage || response.MakeSubChannelIfPossible {
			outgoing = d.Con.GetRtm().NewOutgoingMessage(formatted, envelope.Message.ChannelID, slack.RTMsgOptionTS(envelope.Message.SubChannelID))
		} else {
			outgoing = d.Con.GetRtm().NewOutgoingMessage(formatted, envelope.Message.ChannelID)
		}

		d.Con.GetRtm().SendMessage(outgoing)
		d.core.SetLastSentMessage()

	case core2.ResponseTypeDirectMessage:
		d.Con.GetRtm().SendMessage(
			d.Con.GetRtm().NewTypingMessage(envelope.Message.UserID),
		)

		time.Sleep(650 * time.Millisecond)

		var outgoing = d.Con.GetRtm().NewOutgoingMessage(formatted, envelope.Message.UserID)

		d.Con.GetRtm().SendMessage(outgoing)
		d.core.SetLastSentMessage()

	case core2.ResponseTypeReaction:
		err := d.Con.GetRtm().AddReaction(
			formatted,
			slack.NewRefToMessage(envelope.Message.ChannelID, envelope.Message.SubChannelID),
		)

		if err != nil {
			log.Printf("REACTION ADD FAILED: %v", err)

			db.LogIssue(
				"reaction_add_failed",
				"Failed to add Reaction",
				fmt.Sprintf(
					"Reaction: '%s'\nChannelID: '%s'\nSubChannelID: '%s'\nError: %v",
					formatted,
					envelope.Message.ChannelID,
					envelope.Message.SubChannelID,
					err,
				),
			)
		} else {
			// We do this after we ensure there wasn't an error
			d.core.SetLastSentMessage()
		}
	}
}
