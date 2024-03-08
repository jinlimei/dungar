package core2

import (
	"log"
	"sort"
)

// MessageTriggerGroup is our group handling for backlog
type MessageTriggerGroup struct {
	messageHandlers  []MessageEvHandler
	responseHandlers []ResponseEvHandler
	responses        []*Response
	message          *IncomingMessage
	consumed         bool
}

// Reset resets the full state of the MessageTriggerGroup
func (mtg *MessageTriggerGroup) Reset() {
	mtg.responses = make([]*Response, 0)
	mtg.message = nil
	mtg.consumed = false
}

// SetMessageTriggers sets the list of MessageEvHandler's for the group
func (mtg *MessageTriggerGroup) SetMessageTriggers(triggers []MessageEvHandler) {
	mtg.messageHandlers = triggers

	sort.SliceStable(mtg.messageHandlers, func(i, j int) bool {
		return mtg.messageHandlers[i].Order < mtg.messageHandlers[j].Order
	})
}

// AddMessageTrigger adds a single MessageEvHandler to the group
func (mtg *MessageTriggerGroup) AddMessageTrigger(trigger MessageEvHandler) {
	mtg.messageHandlers = append(mtg.messageHandlers, trigger)

	sort.SliceStable(mtg.messageHandlers, func(i, j int) bool {
		return mtg.messageHandlers[i].Order < mtg.messageHandlers[j].Order
	})
}

// SetResponseTriggers sets a list of ResponseEvHandlers to the group
func (mtg *MessageTriggerGroup) SetResponseTriggers(triggers []ResponseEvHandler) {
	mtg.responseHandlers = triggers

	sort.SliceStable(mtg.responseHandlers, func(i, j int) bool {
		return mtg.responseHandlers[i].Order < mtg.responseHandlers[j].Order
	})
}

// AddResponseTrigger adds a single ResponseEvHandler to the group
func (mtg *MessageTriggerGroup) AddResponseTrigger(trigger ResponseEvHandler) {
	mtg.responseHandlers = append(mtg.responseHandlers, trigger)

	sort.SliceStable(mtg.responseHandlers, func(i, j int) bool {
		return mtg.responseHandlers[i].Order < mtg.responseHandlers[j].Order
	})
}

// Responses returns our slice of Response
func (mtg *MessageTriggerGroup) Responses() []*Response {
	return mtg.responses
}

// Process takes the incoming IncomingMessage and runs it against our
// list of MessageEvHandler's and then finally against our list of
// ResponseEvHandler's before filling `mtg.responses` and `mtg.consumed`
func (mtg *MessageTriggerGroup) Process(svc *Service, msg *IncomingMessage) []*Response {
	mtg.Reset()

	mtg.message = msg

	responses := mtg.processWithHandlers(svc)

	if len(responses) == 0 {
		//log.Printf("After processing, have %d responses\n", len(responses))
		mtg.responses = nil
		return nil
	}

	responses = mtg.processFilters(svc, responses)

	//log.Printf("After filtering, have %d responses\n", len(responses))

	if isAnyCancelled(responses) {
		//log.Printf("At least 1 response is flagged as cancelled, so skipping.\n")
		mtg.responses = nil
		return nil
	}

	//log.Printf("All said and done, %d responses", len(responses))

	mtg.responses = responses

	return responses
}

func (mtg *MessageTriggerGroup) processFilters(svc *Service, rs []*Response) []*Response {
	// No need to go through all the effort of
	// response handlers if there aren't any
	if len(mtg.responseHandlers) == 0 {
		return rs
	}

	var (
		output = rs
		result []*Response
	)

	for _, handler := range mtg.responseHandlers {
		result = handler.Func(svc, output)

		switch handler.Type {
		case Filter:
			output = result
		case Adder:
			output = append(output, result...)
		}
	}

	return output
}

func (mtg *MessageTriggerGroup) processWithHandlers(svc *Service) []*Response {
	var (
		output    = make([]*Response, 0)
		responses []*Response
		consumed  = false
		handled   = false
	)

	for _, handler := range mtg.messageHandlers {
		responses = mtg.runMessageEv(svc, handler)
		consumed = didAnyConsume(responses)
		handled = isAnyHandled(responses)

		// lil log spammy
		log.Printf("Handler (order=%d) '%s' consumed=%v, handled=%v\n",
			handler.Order, handler.Name, consumed, handled)

		// No responses have 'consumed' flagged and no responses had 'handled'
		// flagged (and by extension: cancelled), so we can just move on.
		if !consumed && !handled {
			continue
		}

		// We were actually handling at least one message.
		if handled {
			log.Printf("[processWithHandlers] Handler (order=%d) '%s' consumed=%v, handled=%v\n",
				handler.Order, handler.Name, consumed, handled)
			mtg.dbgMessaging()

			output = append(output, responses...)
		}

		// Even if we did not handle any backlog insofar as a response is
		// generated, we do have the message consumed.
		if consumed {
			break
		}
	}

	return output
}

func (mtg *MessageTriggerGroup) dbgMessaging() {
	var uID string

	if mtg.message.UserID != "" {
		uID = mtg.message.UserID
	}

	log.Printf(
		"[processWithHandlers] Handler for msg UserID='%s',Channel='%s'\n",
		uID,
		mtg.message.ChannelID,
	)
}

func (mtg *MessageTriggerGroup) runMessageEv(svc *Service, handler MessageEvHandler) []*Response {
	var (
		responses = handler.Func(svc, mtg.message)
		outgoing  = make([]*Response, 0)
	)

	//log.Printf("Ran runMessageEv on %s, with respones=%d\n",
	//	handler.Name, len(responses))

	// Filter out nil's
	for _, rsp := range responses {
		if rsp == nil {
			continue
		}

		outgoing = append(outgoing, rsp)
	}

	return outgoing
}
