package triggers

import (
	"fmt"
	"log"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var lastReactResponse = time.Unix(0, 0)

func handleReactions(svc *core2.Service, ev *core2.IncomingEvent) []*core2.EventResponse {
	log.Printf("REACTION BBY: %+v", ev)

	now := time.Now()

	if now.Sub(lastReactResponse).Minutes() <= 1.0 {
		return nil
	}

	if svc.GetBotUser().ID != ev.UserID {
		return nil
	}

	lastReactResponse = now

	if ev.Type == core2.EventReactionAdd {
		log.Printf("Responding to Add")
		return []*core2.EventResponse{
			{
				PrefixUsername: true,
				ConsumedEvent:  true,
				Contents:       fmt.Sprintf("w0w thx beb %s", ev.Text),
			},
		}
	}

	log.Printf("Responding to Remove")
	return []*core2.EventResponse{
		{
			PrefixUsername: true,
			ConsumedEvent:  true,
			Contents:       "w0w rude",
		},
	}
}
