package utils

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMustUseEnvVars(t *testing.T) {
	origCI := os.Getenv(EnvCI)
	origCD := os.Getenv(EnvCD)

	MustSetEnv(EnvCI, "")
	MustSetEnv(EnvCD, "")
	assert.False(t, MustUseEnvVars())

	MustSetEnv(EnvCI, "1")
	MustSetEnv(EnvCD, "1")
	assert.True(t, MustUseEnvVars())

	MustSetEnv(EnvCI, origCI)
	MustSetEnv(EnvCD, origCD)
}

func TestDiscordAllowedOutputChannelIDs(t *testing.T) {
	origCI := os.Getenv(EnvCI)
	origCD := os.Getenv(EnvCD)
	origAllowed := os.Getenv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS")

	MustSetEnv(EnvCI, "1")
	MustSetEnv(EnvCD, "1")
	MustSetEnv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS", "362144221148479489, 471817271749115905 ,,")

	allowed := DiscordAllowedOutputChannelIDs()

	_, hasGeneral := allowed["362144221148479489"]
	_, hasBaconbar := allowed["471817271749115905"]

	assert.Len(t, allowed, 2)
	assert.True(t, hasGeneral)
	assert.True(t, hasBaconbar)

	MustSetEnv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS", origAllowed)
	MustSetEnv(EnvCI, origCI)
	MustSetEnv(EnvCD, origCD)
}

func TestDiscordAllowedLearningChannelIDs(t *testing.T) {
	origCI := os.Getenv(EnvCI)
	origCD := os.Getenv(EnvCD)
	origAllowed := os.Getenv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS")

	MustSetEnv(EnvCI, "1")
	MustSetEnv(EnvCD, "1")
	MustSetEnv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS", "362144221148479489")

	allowed := DiscordAllowedLearningChannelIDs()
	_, hasGeneral := allowed["362144221148479489"]

	assert.Len(t, allowed, 1)
	assert.True(t, hasGeneral)

	MustSetEnv("DUNGAR_DISCORD_ALLOWED_OUTPUT_CHANNEL_IDS", origAllowed)
	MustSetEnv(EnvCI, origCI)
	MustSetEnv(EnvCD, origCD)
}
