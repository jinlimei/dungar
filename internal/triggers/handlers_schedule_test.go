package triggers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleScheduleEvent(t *testing.T) {
	svc := initMockServices()

	scheduleHandler(svc, 0, time.Unix(0, 0))
	assert.NotNil(t, scheduleGroup)
}

func TestInitScheduleGroup(t *testing.T) {
	svc := initMockServices()

	scheduleGroup = nil

	scheduleHandler(svc, 0, time.Unix(0, 0))
	assert.NotNil(t, scheduleGroup)
}

func brokenTestSchedulingNotSpammingThePlanet(t *testing.T) {
	// FIXME fix this test
	//svc := initMockServices()
	//
	//commentator := &core2.ScheduleEvHandler{
	//	Order: 1,
	//	Name:  "Commentator",
	//	Func: func(_ *core2.Service) []*core2.ScheduledMessage {
	//		return []*core2.ScheduledMessage{
	//			{"butts", false, "hello, world", time.Now()},
	//		}
	//	},
	//	RunTick: 6 * time.Hour,
	//	LastRan: 0,
	//}
	//
	//scheduleGroup = core2.NewScheduleTriggerGroup()
	//scheduleGroup.SetScheduleEvs([]*core2.ScheduleEvHandler{commentator})
	//
	//assert.Equal(t, int64(0), commentator.LastRan)
	//msgs := mockTrigger.scheduleHandler(svc, 0, time.Unix(0, 0))
	//assert.Equal(t, 1, len(msgs))
	//assert.NotEqual(t, int64(0), commentator.LastRan)
	//msgs = mockTrigger.scheduleHandler(svc, 0, time.Unix(0, 0))
	//assert.Equal(t, 0, len(msgs))
}
