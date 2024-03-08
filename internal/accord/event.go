package accord

import (
	"log"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) handleEventResponseEnvelope(env *core2.EventResponseEnvelope) {
	if env == nil || len(env.Responses) == 0 {
		log.Printf("env or env.Responses is empty/nil/whatever")
		return
	}

	if d.outgoing != nil {
		log.Printf("d.outgoing is not nil")
		conv := make([]*core2.Response, 0, len(env.Responses))

		for _, rsp := range env.Responses {
			conv = append(conv, rsp.ToResponse())
		}

		d.outgoing.AddEnvelope(&core2.ResponseEnvelope{
			Message: &core2.IncomingMessage{
				ID:        env.Event.ID,
				ServerID:  env.Event.ServerID,
				ChannelID: env.Event.ChannelID,
				UserID:    env.Event.UserID,
				Contents:  env.Event.Text,
			},
			Responses: conv,
		})
	}
}
