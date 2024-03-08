package core2

import (
	"sort"
	"time"
)

// ScheduleTriggerGroup is our grouping for all the behaviors necessary
// for a functioning trigger group
type ScheduleTriggerGroup struct {
	schedules []*ScheduleEvHandler
	responses []ScheduleResponseEvHandler
	messages  []*ScheduledMessage
}

// NewScheduleTriggerGroup initializes internal values to work
// with the scheduler.
func NewScheduleTriggerGroup() *ScheduleTriggerGroup {
	return &ScheduleTriggerGroup{
		schedules: make([]*ScheduleEvHandler, 0),
		responses: make([]ScheduleResponseEvHandler, 0),
		messages:  make([]*ScheduledMessage, 0),
	}
}

// Reset resets the backlog of the group
func (stg *ScheduleTriggerGroup) Reset() {
	stg.messages = make([]*ScheduledMessage, 0)
}

// Messages returns the backlog of the group
func (stg *ScheduleTriggerGroup) Messages() []*ScheduledMessage {
	return stg.messages
}

// SetScheduleEvs provides a list of ScheduleEvHandler's to use for the
// group
func (stg *ScheduleTriggerGroup) SetScheduleEvs(sts []*ScheduleEvHandler) {
	stg.schedules = sts

	sort.SliceStable(stg.schedules, func(i, j int) bool {
		return stg.schedules[i].Order < stg.schedules[j].Order
	})
}

// AddScheduleEv adds a single ScheduleEvHandler to the schedule listing
func (stg *ScheduleTriggerGroup) AddScheduleEv(st *ScheduleEvHandler) {
	stg.schedules = append(stg.schedules, st)

	sort.SliceStable(stg.schedules, func(i, j int) bool {
		return stg.schedules[i].Order < stg.schedules[j].Order
	})
}

// UpdateScheduleEv updates an existing ScheduleEvHandler in the group.
// Useful for changing the time/duration of a schedule handler
func (stg *ScheduleTriggerGroup) UpdateScheduleEv(st *ScheduleEvHandler) {
	for idx, handler := range stg.schedules {
		if handler.Name == st.Name {
			stg.schedules[idx] = handler
			break
		}
	}
}

// RemoveScheduleEv removes the ScheduleEvHandler from the group if it
// does in fact exist in the group
func (stg *ScheduleTriggerGroup) RemoveScheduleEv(name string) {
	output := make([]*ScheduleEvHandler, 0)

	for _, handler := range stg.schedules {
		if handler.Name != name {
			output = append(output, handler)
		}
	}

	stg.schedules = output

	sort.SliceStable(stg.schedules, func(i, j int) bool {
		return stg.schedules[i].Order < stg.schedules[j].Order
	})
}

// SetResponseEvs sets all the ScheduleResponseEvHandler's for a group
func (stg *ScheduleTriggerGroup) SetResponseEvs(sts []ScheduleResponseEvHandler) {
	stg.responses = sts

	sort.SliceStable(stg.responses, func(i, j int) bool {
		return stg.responses[i].Order < stg.responses[j].Order
	})
}

// AddResponseEv adds a single ScheduleResponseEvHandler to the group
func (stg *ScheduleTriggerGroup) AddResponseEv(srt ScheduleResponseEvHandler) {
	stg.responses = append(stg.responses, srt)

	sort.SliceStable(stg.responses, func(i, j int) bool {
		return stg.responses[i].Order < stg.responses[j].Order
	})
}

// Process runs all the ScheduleEvHandler's for the group which fit
// the time needs to output the messaging. All ScheduledMessage's are
// then setLastRan through the list of ScheduleResponseEvHandler's
func (stg *ScheduleTriggerGroup) Process(svc *Service) {
	stg.Reset()

	var (
		messages = make([]*ScheduledMessage, 0)
		now      = time.Now().UTC().Unix()

		output []*ScheduledMessage
	)

	for _, handler := range stg.schedules {
		//log.Printf("handler '%v': runTick='%v', lastRun='%v', now='%v'\n",
		//	handler.Name, handler.RunTick.Seconds(), handler.LastRan, now)

		if (now - handler.LastRan) >= int64(handler.RunTick.Seconds()) {
			handler.LastRan = now
			output = stg.runSchedule(svc, handler)

			if output == nil || len(output) == 0 {
				//log.Printf("handle '%v': had no output\n", handler.Name)
				continue
			}

			//log.Printf("handle '%v': had %d output\n", handler.Name, len(output))
			messages = append(messages, output...)
		}
	} // end for

	anyCancelled := false
	for _, rsp := range messages {
		if rsp.Cancelled {
			anyCancelled = true
			break
		}
	}

	if anyCancelled {
		stg.messages = []*ScheduledMessage{}
	} else {
		stg.messages = stg.processFilters(svc, messages)
	}
}

// processFilters will take a set of scheduled messages and then filter them with a specific ResponseHandlerType
// for Filter response handling it will convert an output message
// for Adder it will add a new additional message
func (stg *ScheduleTriggerGroup) processFilters(svc *Service, msgs []*ScheduledMessage) []*ScheduledMessage {
	if len(stg.responses) == 0 {
		return msgs
	}

	var (
		output = msgs
		result []*ScheduledMessage
	)

	for _, handler := range stg.responses {
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

func (stg *ScheduleTriggerGroup) runSchedule(svc *Service, handler *ScheduleEvHandler) []*ScheduledMessage {
	var (
		messages = make([]*ScheduledMessage, 0)
		results  = handler.Func(svc)
	)

	handler.setLastRan()

	if results == nil || len(results) == 0 {
		return nil
	}

	for _, msg := range results {
		if msg == nil {
			continue
		}

		if msg.IsEmpty() {
			continue
		}

		messages = append(messages, msg)
	}

	return messages
}
