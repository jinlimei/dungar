package triggers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func withFreshManicMinuteState(t *testing.T) {
	t.Helper()

	savedState := manicMinuteState
	savedPicker := manicMinuteTriggerWordPicker
	savedChance := masterChanceList["manicMinuteHandler--start"]

	manicMinuteState = &manicMinuteManager{}

	t.Cleanup(func() {
		manicMinuteState = savedState
		manicMinuteTriggerWordPicker = savedPicker
		masterChanceList["manicMinuteHandler--start"] = savedChance
	})
}

func sequentialManicMinuteWordPicker(words ...string) func() string {
	idx := 0

	return func() string {
		if len(words) == 0 {
			return ""
		}

		if idx >= len(words) {
			return words[len(words)-1]
		}

		word := words[idx]
		idx++
		return word
	}
}

func TestMessageContainsTriggerWord(t *testing.T) {
	assert.True(t, messageContainsTriggerWord("I do skill things", "skill"))
	assert.True(t, messageContainsTriggerWord("I do skill things", "SKILL"))
	assert.True(t, messageContainsTriggerWord("I want to go home", "to"))
	assert.False(t, messageContainsTriggerWord("I brought a skillet", "skill"))
}

func TestManicMinuteHandlerStartsAndSchedules(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)
	withFreshManicMinuteState(t)

	svc := initMockServices()
	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill")
	masterChanceList["manicMinuteHandler--start"] = 0.0
	mockDriver.SetUser("bob", "Bob")

	InitializeManicMinute()

	msg := makeMessage("I do skill things", "bob", "butts")
	msg.ServerID = "guild"
	rsp := manicMinuteHandler(svc, msg)

	assert.Len(t, rsp, 2)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.Equal(t, `Bob triggered a dungarmatic manic minute by saying "skill".`, rsp[0].Contents)
	assert.NotEmpty(t, rsp[1].Contents)
	assert.True(t, manicMinuteState.isActive())

	scheduled := manicMinuteScheduler(svc)
	assert.Len(t, scheduled, 1)
	assert.Equal(t, "butts", scheduled[0].ChannelID)
	assert.NotEmpty(t, scheduled[0].Contents)
}

func TestManicMinuteHandlerStopsAndRotatesTriggerWord(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)
	withFreshManicMinuteState(t)

	svc := initMockServices()
	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill", "banana")
	mockDriver.AddUserRole("bob", "admin")

	InitializeManicMinute()

	startMsg := makeMessage("I do skill things", "bob", "butts")
	startMsg.ServerID = "guild"
	rsp := manicMinuteHandler(svc, startMsg)
	assert.Len(t, rsp, 2)
	assert.True(t, manicMinuteState.isActive())
	assert.Equal(t, "skill", manicMinuteState.currentTriggerWord())

	stopMsg := makeMessage("@Dungar shut up!", "bob", "butts")
	stopMsg.ServerID = "guild"
	rsp = manicMinuteHandler(svc, stopMsg)

	assert.Len(t, rsp, 2)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.Equal(t, "This concludes dungarmatic's manic minute.", rsp[0].Contents)
	assert.Equal(t, "New trigger word set.", rsp[1].Contents)
	assert.False(t, manicMinuteState.isActive())
	assert.Equal(t, "banana", manicMinuteState.currentTriggerWord())
}

func TestManicMinuteDirectedStopRequiresAdmin(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)
	withFreshManicMinuteState(t)

	svc := initMockServices()
	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill", "banana")

	InitializeManicMinute()

	startMsg := makeMessage("I do skill things", "bob", "butts")
	startMsg.ServerID = "guild"
	rsp := manicMinuteHandler(svc, startMsg)
	assert.Len(t, rsp, 2)
	assert.True(t, manicMinuteState.isActive())

	stopMsg := makeMessage("@Dungar shut up!", "bob", "butts")
	stopMsg.ServerID = "guild"
	rsp = manicMinuteHandler(svc, stopMsg)

	assert.True(t, isEmptyRsp(rsp))
	assert.True(t, manicMinuteState.isActive())
	assert.Equal(t, "skill", manicMinuteState.currentTriggerWord())
}

func TestManicMinuteSchedulerEndsAndRotatesTriggerWord(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)
	withFreshManicMinuteState(t)

	svc := initMockServices()
	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill", "banana")

	InitializeManicMinute()

	startMsg := makeMessage("I do skill things", "bob", "butts")
	startMsg.ServerID = "guild"
	rsp := manicMinuteHandler(svc, startMsg)
	assert.Len(t, rsp, 2)
	assert.True(t, manicMinuteState.isActive())

	manicMinuteState.lock.Lock()
	manicMinuteState.endsAt = time.Now().Add(-time.Second)
	manicMinuteState.lock.Unlock()

	scheduled := manicMinuteScheduler(svc)

	assert.Len(t, scheduled, 2)
	assert.Equal(t, "This concludes dungarmatic's manic minute.", scheduled[0].Contents)
	assert.Equal(t, "New trigger word set.", scheduled[1].Contents)
	assert.False(t, manicMinuteState.isActive())
	assert.Equal(t, "banana", manicMinuteState.currentTriggerWord())
}

func TestManicMinuteAdminHandlerRequiresAdminRole(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)
	withFreshManicMinuteState(t)

	svc := initMockServices()
	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill")

	msg := makeMessage("!manic status", "bob", "butts")
	msg.ServerID = "guild"

	rsp := manicMinuteAdminHandler(svc, msg)
	assert.True(t, isEmptyRsp(rsp))

	mockDriver.AddUserRole("bob", "admin")
	rsp = manicMinuteAdminHandler(svc, msg)
	assert.Len(t, rsp, 1)
	assert.Contains(t, rsp[0].Contents, "MANIC MINUTE")
}

func TestManicMinuteAdminRotate(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)
	withFreshManicMinuteState(t)

	svc := initMockServices()
	mockDriver.AddUserRole("bob", "admin")
	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill", "banana")

	InitializeManicMinute()

	msg := makeMessage("!manic rotate", "bob", "butts")
	msg.ServerID = "guild"
	rsp := manicMinuteAdminHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.Contains(t, rsp[0].Contents, "banana")
	assert.Equal(t, "banana", manicMinuteState.currentTriggerWord())
}

func TestManicMinuteStatusShowsGuaranteedTriggerChance(t *testing.T) {
	withFreshManicMinuteState(t)

	manicMinuteTriggerWordPicker = sequentialManicMinuteWordPicker("skill", "banana")
	masterChanceList["manicMinuteHandler--start"] = 0.20

	InitializeManicMinute()

	status := manicMinuteState.statusSnapshot()
	assert.Equal(t, "skill", status.TriggerWord)
	assert.InDelta(t, 1.0, status.StartChance, 0.0001)
	assert.Equal(t, 0, status.StartMisses)

	_, rotated := manicMinuteState.rotate("test-rotate")
	assert.Equal(t, "banana", rotated)

	status = manicMinuteState.statusSnapshot()
	assert.Equal(t, "banana", status.TriggerWord)
	assert.InDelta(t, 1.0, status.StartChance, 0.0001)
	assert.Equal(t, 0, status.StartMisses)
}

func TestPickManicMinuteCooldownDurationFixed(t *testing.T) {
	assert.Equal(t, 30*time.Minute, pickManicMinuteCooldownDuration())
}
