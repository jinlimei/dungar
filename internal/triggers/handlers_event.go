package triggers

import (
	"log"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func eventHandler(svc *core2.Service, ev *core2.IncomingEvent) *core2.EventResponseEnvelope {
	if ev.Type.IsChannelEvent() {
		trackChannelFromEvent(svc, ev)
	}

	if ev.Type == core2.EventReactionAdd || ev.Type == core2.EventReactionRemove {
		responses := handleReactions(svc, ev)

		if len(responses) > 0 {
			return &core2.EventResponseEnvelope{
				Event:     ev,
				Responses: responses,
			}
		}
	}

	return nil
}

func trackChannelFromEvent(svc *core2.Service, ev *core2.IncomingEvent) {
	channel, err := svc.GetChannel(ev.ChannelID, ev.ServerID)

	if err != nil {
		log.Printf(
			"ERROR: failed to retrieve channel '%s' / server '%s': %v",
			ev.ChannelID,
			ev.ServerID,
			err,
		)

		return
	}

	err = tracking.StoreChannel(&channel)

	if err != nil {
		log.Printf(
			"ERROR: failed to store channel '%s' / server '%s': %v",
			ev.ChannelID,
			ev.ServerID,
			err,
		)
	}
}
