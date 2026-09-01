package triggers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func TestShouldRecordDiscordRawMessage(t *testing.T) {
	origCI := os.Getenv(utils.EnvCI)
	origCD := os.Getenv(utils.EnvCD)
	origAllowed := os.Getenv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS")

	utils.MustSetEnv(utils.EnvCI, "1")
	utils.MustSetEnv(utils.EnvCD, "1")
	utils.MustSetEnv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS", "362144221148479489,471817271749115905")

	assert.True(t, shouldRecordDiscordRawMessage("362144221148479489"))
	assert.True(t, shouldRecordDiscordRawMessage("471817271749115905"))
	assert.False(t, shouldRecordDiscordRawMessage("999"))

	utils.MustSetEnv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS", origAllowed)
	utils.MustSetEnv(utils.EnvCI, origCI)
	utils.MustSetEnv(utils.EnvCD, origCD)
}
