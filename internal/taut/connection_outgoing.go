package taut

import (
	"log"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func (d *Driver) outgoingEventLoop() {
	var nilConWait = 0

	for {
		if d.Con.GetRtm() == nil {
			if nilConWait > 5 {
				db.LogIssue("nil_con_wait", "NilConWait", "nilConWait hit")
				panic("nil_con_wait hit")
			}

			nilConWait++
			log.Printf("Received nilConWait(cnt=%d)\n", nilConWait)
			time.Sleep(5 * time.Second)
			continue
		}

		// We're not connected yet
		if d.GetBotUser().ID == "" {
			log.Printf("Waiting to connect before handling outgoingEventLoop")
			time.Sleep(5 * time.Second)
			continue
		}

		outgoing := d.GetOutgoingResponses()

		if utils.IsSilentRunning() {
			if outgoing.HasEnvelopes() {
				log.Printf("silent_running=true stopped %d envelopes from being processed.", outgoing.Length())
			}
			continue
		}

		if outgoing == nil {
			db.LogIssue("outgo_is_nil", "Outgoing Envelope is <nil>", "")
			continue
		}

		d.runScheduledTasks(outgoing)

		if !outgoing.HasEnvelopes() {
			//log.Printf("No outgoing envelopes for incoming message.\n")
			time.Sleep(1 * time.Second)
			continue
		}

		d.handleOutgoingResponses(outgoing)
	} // end of for (infinite loop)
}
