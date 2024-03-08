package accord

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"log"
	"time"
)

func (d *Driver) outgoingEventLoop() {
	var nilSessionWait = 0

	for {
		if d.Con.GetSession() == nil {
			if nilSessionWait > 5 {
				db.LogIssue("nil_con_wait", "NilConWait", "nilcon>5")
				panic("nil_con_wait hit")
			}

			nilSessionWait++
			log.Printf("Received nilSessionWait(cnt=%d)", nilSessionWait)
			time.Sleep(5 * time.Second)
			continue
		}

		if d.botUser == nil || d.botUser.ID == "" {
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
	}
}
