package triggers

import (
	"time"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var scheduleGroup *core2.ScheduleTriggerGroup

// scheduleHandler will let us schedule up messages to send out if nothing is going on ~
func scheduleHandler(svc *core2.Service, currentOutgoingCount uint, lastSentMessage time.Time) core2.ScheduledMessages {
	if scheduleGroup == nil {
		initScheduleGroup()
	}

	scheduleGroup.Process(svc)

	return scheduleGroup.Messages()
}

func initScheduleGroup() {
	scheduleGroup = core2.NewScheduleTriggerGroup()

	scheduleGroup.SetScheduleEvs([]*core2.ScheduleEvHandler{
		{0, "OutputQueue", outputQueueScheduler, 1 * time.Second, 0},
		{1, "Commentator", commentateScheduler, 6 * time.Hour, 0},
	})
}
