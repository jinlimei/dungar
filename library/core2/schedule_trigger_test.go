package core2

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func makeScheduleGroup() *ScheduleTriggerGroup {
	return &ScheduleTriggerGroup{
		schedules: []*ScheduleEvHandler{},
		responses: []ScheduleResponseEvHandler{},
		messages:  nil,
	}
}

func schedTimeMsg(_ *Service) []*ScheduledMessage {
	now := time.Now()
	m := &ScheduledMessage{
		ChannelID: "ABCDEF",
		Cancelled: false,
		Contents:  fmt.Sprintf("The time now is: %v\n", now.Format(time.RFC822Z)),
		SentAt:    now,
	}

	return []*ScheduledMessage{m}
}

func schedNothing1(_ *Service) []*ScheduledMessage {
	return nil
}

func schedNothing2(_ *Service) []*ScheduledMessage {
	return []*ScheduledMessage{}
}

func TestScheduleProcessNoOps(t *testing.T) {
	now := time.Now()
	svc := &Service{}

	grp := makeScheduleGroup()
	grp.SetScheduleEvs([]*ScheduleEvHandler{
		{0, "SchedNothing1", schedNothing1, 5 * time.Second, 0},
		{1, "SchedNothing2", schedNothing2, 15 * time.Second, 0},
	})

	// We run. Now!
	grp.Process(svc)

	msgs := grp.Messages()
	assert.NotNil(t, msgs)
	assert.Len(t, msgs, 0)

	for idx := range grp.schedules {
		grp.schedules[idx].LastRan = now.Unix() - 601
	}

	// We run. Now!
	grp.Process(svc)

	msgs = grp.Messages()
	assert.NotNil(t, msgs)
	assert.Len(t, msgs, 0)
}

func TestScheduleMixedOps(t *testing.T) {
	now := time.Now()
	svc := &Service{}

	grp := makeScheduleGroup()
	grp.SetScheduleEvs([]*ScheduleEvHandler{
		{0, "SchedNothing1", schedNothing1, 5 * time.Minute, 0},
		{1, "SchedTimeMsg", schedTimeMsg, 15 * time.Minute, 0},
	})

	// We run. Now!
	grp.Process(svc)

	msgs := grp.Messages()
	assert.NotNil(t, msgs)
	assert.Len(t, msgs, 1)

	for idx := range grp.schedules {
		grp.schedules[idx].LastRan = now.Unix() - 601
	}

	// We run. Now!
	grp.Process(svc)

	msgs = grp.Messages()
	assert.NotNil(t, msgs)
	assert.Len(t, msgs, 0)

	for idx := range grp.schedules {
		grp.schedules[idx].LastRan = now.Unix() - 901
	}

	// We run. Now!
	grp.Process(svc)

	msgs = grp.Messages()
	assert.NotNil(t, msgs)
	assert.Len(t, msgs, 1)

}
