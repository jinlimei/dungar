package taut

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) incomingEventLoop() {
	for msg := range d.Con.GetRtm().IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectingEvent:
			log.Printf("Connecting: '%s'", spew.Sdump(ev))
		case *slack.ConnectedEvent:
			d.SetBotUser(&core2.BotUser{
				ID:    ev.Info.User.ID,
				Name:  strings.ToLower(ev.Info.User.Name),
				IsBot: true,
			})

			tpl := `
Connected, Yay!
Count: %d
URL: %s
User: ID=%s,Name=%s
Team: ID=%s,Name=%s,Domain=%s
`

			log.Printf(
				tpl,
				ev.ConnectionCount,
				ev.Info.URL,
				ev.Info.User.ID,
				ev.Info.User.Name,
				ev.Info.Team.ID,
				ev.Info.Team.Name,
				ev.Info.Team.Domain,
			)

			d.teamID = ev.Info.Team.ID

		case *slack.MessageEvent:
			d.handleMessageEvent(ev)
		case *slack.ReactionAddedEvent:
			d.handleReactionAdded(ev)
		case *slack.ReactionRemovedEvent:
			d.handleReactionRemoved(ev)
		case *slack.RTMError:
			log.Println("Error")
			log.Println(spew.Sdump(ev))
		case *slack.AckMessage:
			log.Printf("slack rtm acknowledged our message as (true=good) %v\n", ev.Ok)

			if !ev.Ok {
				// TODO sentry messaging
				//utils.SentryMessage(
				//	fmt.Sprintf("Slack failed to acknowledge our message: %v?", spew.Sdump(ev)),
				//	nil,
				//)

				log.Print(spew.Sdump(ev.Error, ev.RTMResponse))
			}
		case *slack.PinAddedEvent:
			d.handlePinAdded(ev)
		case *slack.PinRemovedEvent:
			d.handlePinRemoved(ev)
		case *slack.MemberJoinedChannelEvent:
			d.handleMemberJoinedChannel(ev)
		case *slack.MemberLeftChannelEvent:
			d.handleMemberLeftChannel(ev)
		case *slack.ConnectionErrorEvent:
			if ev.Attempt >= 5 {
				msg := fmt.Sprintf("Could not connect after 5 attempts: %v\n", ev.Error())

				//utils.SentryMessage(msg, nil)

				log.Println(msg)
				log.Println("Waiting 60s before panic and shutdown.")

				time.Sleep(1 * time.Minute)
				panic(msg)
			}

			log.Printf(
				"Connection Error Event (attempt=%d): %+v\n",
				ev.Attempt,
				ev.Error(),
			)

			log.Printf(spew.Sdump(ev))
			time.Sleep(29 * time.Second)
		default:
			jsn, _ := json.Marshal(ev)

			fmt.Printf("OTHER EVENT LOGGING:\n%s\nJSON:\n%s\n", spew.Sdump(ev), string(jsn))
		}
	}
}
