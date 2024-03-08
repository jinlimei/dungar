package accord

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) handleReactionAdd(s *discordgo.Session, ev *discordgo.MessageReactionAdd) {
	env := d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.MessageID,
		UserID:    ev.UserID,
		ChannelID: ev.ChannelID,
		Text:      ev.Emoji.APIName(),
		Type:      core2.EventReactionAdd,
	})

	d.handleEventResponseEnvelope(env)
}

func (d *Driver) handleReactionRemove(s *discordgo.Session, ev *discordgo.MessageReactionRemove) {
	env := d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.MessageID,
		UserID:    ev.UserID,
		ChannelID: ev.ChannelID,
		Text:      ev.Emoji.APIName(),
		Type:      core2.EventReactionRemove,
	})

	d.handleEventResponseEnvelope(env)
}
