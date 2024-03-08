package taut

import (
	"encoding/json"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) handleReactionAdded(ev *slack.ReactionAddedEvent) {
	// returns EventResponseEnvelope
	env := d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:          ev.EventTimestamp,
		UserID:      ev.User,
		ChannelID:   ev.Item.Channel,
		ServerID:    d.teamID,
		Text:        ev.Reaction,
		Type:        core2.EventReactionAdd,
		Attachments: nil,
	})

	d.handleEventResponseEnvelope(env)
}

func (d *Driver) handleReactionRemoved(ev *slack.ReactionRemovedEvent) {
	env := d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:          ev.EventTimestamp,
		UserID:      ev.User,
		ChannelID:   ev.Item.Channel,
		ServerID:    d.teamID,
		Text:        ev.Reaction,
		Type:        core2.EventReactionRemove,
		Attachments: nil,
	})

	d.handleEventResponseEnvelope(env)
}

func (d *Driver) handleEventResponseEnvelope(env *core2.EventResponseEnvelope) {
	if env == nil || len(env.Responses) == 0 {
		return
	}

	if d.outgoing != nil {
		conv := make([]*core2.Response, 0, len(env.Responses))

		for _, rsp := range env.Responses {
			conv = append(conv, rsp.ToResponse())
		}

		d.outgoing.AddEnvelope(&core2.ResponseEnvelope{
			Message: &core2.IncomingMessage{
				ID:        env.Event.ID,
				ServerID:  env.Event.ServerID,
				ChannelID: env.Event.ChannelID,
				Contents:  env.Event.Text,
			},
			Responses: conv,
		})
	}
}

func (d *Driver) handlePinAdded(ev *slack.PinAddedEvent) {
	log.Printf("handlePinAdded: %s", spew.Sdump(ev))
	jsn, _ := json.Marshal(ev)
	log.Printf("handlePinAdded JSON: %s", string(jsn))
}

func (d *Driver) handlePinRemoved(ev *slack.PinRemovedEvent) {
	log.Printf("handlePinRemoved: %s", spew.Sdump(ev))
	jsn, _ := json.Marshal(ev)
	log.Printf("handlePinRemoved JSON: %s", string(jsn))
}

func (d *Driver) handleMemberJoinedChannel(ev *slack.MemberJoinedChannelEvent) {
	log.Printf("handleMemberJoined: %s", spew.Sdump(ev))
	jsn, _ := json.Marshal(ev)
	log.Printf("handleMemberJoined JSON: %s", string(jsn))
}

func (d *Driver) handleMemberLeftChannel(ev *slack.MemberLeftChannelEvent) {
	log.Printf("handleMemberLeft: %s", spew.Sdump(ev))
	jsn, _ := json.Marshal(ev)
	log.Printf("handleMemberLeft JSON: %s", string(jsn))
}
