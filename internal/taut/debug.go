package taut

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
)

var logMessageEvNum = 0

func logMessageEvent(ev *slack.MessageEvent) {
	jsn, _ := json.MarshalIndent(ev, "", "  ")
	now := time.Now()
	name := fmt.Sprintf(
		"logs/%s-%s-%d.json",
		now.Format("2006-01-02T15-04-05-000"),
		ev.Msg.Timestamp,
		logMessageEvNum,
	)

	logMessageEvNum++

	os.WriteFile(name, jsn, 0664)
}

func (d *Driver) debugMessaging(ev *slack.MessageEvent) {
	parsed := parseSlackMessage(ev.Msg.Text)

	log.Printf(
		"[msg] (type=%s,channel[%s]=%s,user[%s]=%s) %s\n",
		ev.Msg.Type,
		ev.Msg.Channel,
		d.GetChannelName(ev.Msg.Channel, ""),
		ev.Msg.User,
		d.GetUserName(ev.Msg.User, ""),
		parsed.Converted,
	)

	log.Print(spew.Sdump(ev.Msg))

	printAttachments("msg", ev)

	if ev.SubMessage != nil {
		subParsed := parseSlackMessage(ev.SubMessage.Text)

		log.Printf(
			"[sub-msg] (type=%s,channel[%s]=%s,user[%s]=%s) %s\n",
			ev.SubMessage.Type,
			ev.SubMessage.Channel,
			d.GetChannelName(ev.SubMessage.Channel, ""),
			ev.SubMessage.User,
			d.GetUserName(ev.SubMessage.User, ""),
			subParsed.Converted,
		)

		log.Print(spew.Sdump(ev.SubMessage))
		printAttachments("sub-msg", ev)
	}
}

func printAttachments(area string, ev *slack.MessageEvent) {
	var attachments []slack.Attachment

	if area == "msg" {
		attachments = ev.Msg.Attachments
	} else {
		attachments = ev.SubMessage.Attachments
	}

	log.Printf("[%s] ATTACHMENTS (len=%d)\n", area, len(attachments))

	for num, attach := range attachments {
		log.Printf(
			"[%s] Attachment %d: (author=%s)\n%s\n",
			area,
			num,
			attach.AuthorName,
			attach.Text,
		)
	}
}
