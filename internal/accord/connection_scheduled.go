package accord

import "gitlab.int.magneato.site/dungar/prototype/library/core2"

func (d *Driver) runScheduledTasks(outgoing *core2.OutgoingResponses) {
	var (
		outgoingCnt = uint(outgoing.Length())
		scheduled   = d.core.HandleSchedule(outgoingCnt)
	)

	if scheduled != nil && len(scheduled) > 0 {
		envelopes := scheduled.ToResponseEnvelopes()
		for _, envelope := range envelopes {
			outgoing.AddEnvelope(envelope)
		}
	}
}
