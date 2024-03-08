package triggers

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	lastSentMessage = time.Now()

	hasSaidTopWorks = false
	hasSaidSubWorks = false

	hasSaidEmojiWorks = false
)

func debugHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	//channel, _ := svc.GetChannel(msg.ChannelID)
	//
	//if channel.Name != "dungar-test" {
	//	return core2.EmptyRsp()
	//}

	now := time.Now()
	now = now.Add(-5 * time.Second)

	if lastSentMessage.After(now) {
		log.Printf("lastSentMessage is after 'now' (now - 5seconds), skipping!")
		return core2.EmptyRsp()
	}

	if strings.Contains(msg.Contents, "dungar") {
		if strings.Contains(msg.Contents, "react") {
			return core2.SingleReactionRsp("rowletlove")

		} else if strings.Contains(msg.Contents, "emoji") && !hasSaidEmojiWorks {
			hasSaidEmojiWorks = true
			return core2.MakeSingleRsp("yeah here's an emoji: :rowletlove: oh maybe another one :pensive:")

		} else if msg.IsSubMessage && !hasSaidSubWorks {
			hasSaidSubWorks = true
			return core2.MakeSingleRsp("yeah it works in a sub-message neat")

		} else if !hasSaidTopWorks {
			hasSaidTopWorks = true
			return core2.MakeSingleRsp("yeah it works bby")
		}
	}

	if len(msg.Attachments) > 0 {
		// idk what this shit is
		outgoing := make([]string, 0)
		log.Printf("debugHandler, attachments: %d", len(msg.Attachments))

		if len(msg.Attachments) > 0 {
			for idx, attach := range msg.Attachments {
				outgoing = append(outgoing, fmt.Sprintf("Attachment %d\n%s", idx+1, attach.String()))
			}
		}

		//log.Printf("Outgoing Text: %+v", outgoing)
		return core2.MakeSingleRsp(strings.Join(outgoing, "\n"))
	}

	return core2.EmptyRsp()
}
