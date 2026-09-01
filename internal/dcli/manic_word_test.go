package dcli

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
)

func TestNormalizeManicMinuteProtocolMode(t *testing.T) {
	mode, err := normalizeManicMinuteProtocolMode("discord")
	assert.NoError(t, err)
	assert.Equal(t, "discord", mode)

	mode, err = normalizeManicMinuteProtocolMode("Slack")
	assert.NoError(t, err)
	assert.Equal(t, "slack", mode)
}

func TestNormalizeManicMinuteProtocolModeUnknown(t *testing.T) {
	mode, err := normalizeManicMinuteProtocolMode("xmpp")
	assert.Error(t, err)
	assert.Equal(t, "", mode)
}

func TestFormatManicMinuteWordWithCooldown(t *testing.T) {
	now := time.Date(2026, time.June, 9, 16, 30, 0, 0, time.UTC)
	state := &db.ManicMinuteRuntimeState{
		TriggerWord:   "week",
		HasCooldown:   true,
		CooldownUntil: now.Add(17*time.Minute + 30*time.Second),
	}

	got := formatManicMinuteWord(state, now)
	assert.Equal(t, "Current manic-word: week (cooldown remaining: 17m30s)", got)
}

func TestFormatManicMinuteWordWithoutCooldown(t *testing.T) {
	now := time.Date(2026, time.June, 9, 16, 30, 0, 0, time.UTC)
	state := &db.ManicMinuteRuntimeState{
		TriggerWord: "week",
	}

	got := formatManicMinuteWord(state, now)
	assert.Equal(t, "Current manic-word: week (cooldown remaining: none)", got)
}

func TestSanitizeManicMinuteCLITriggerWord(t *testing.T) {
	assert.Equal(t, "bacon", sanitizeManicMinuteCLITriggerWord(" Bacon "))
	assert.Equal(t, "skill_issue", sanitizeManicMinuteCLITriggerWord("skill_issue"))
	assert.Equal(t, "to", sanitizeManicMinuteCLITriggerWord("to"))
	assert.Equal(t, "", sanitizeManicMinuteCLITriggerWord("x"))
	assert.Equal(t, "", sanitizeManicMinuteCLITriggerWord("https://example.com"))
	assert.Equal(t, "", sanitizeManicMinuteCLITriggerWord("@dungar"))
}

func TestBuildManicMinuteCLIAutomaticCandidateWords(t *testing.T) {
	got := buildManicMinuteCLIAutomaticCandidateWords("to Bacon week bacon <@123> https://example.com skill-issue")
	assert.Equal(t, []string{"bacon", "week", "skill-issue"}, got)
}
