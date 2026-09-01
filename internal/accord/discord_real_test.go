package accord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestDiscordGatewayIntents(t *testing.T) {
	intents := discordGatewayIntents()

	assert.NotZero(t, intents&discordgo.IntentsGuilds)
	assert.NotZero(t, intents&discordgo.IntentsGuildMembers)
	assert.NotZero(t, intents&discordgo.IntentsGuildEmojis)
	assert.NotZero(t, intents&discordgo.IntentsGuildMessages)
	assert.NotZero(t, intents&discordgo.IntentsGuildMessageReactions)
	assert.NotZero(t, intents&discordgo.IntentsMessageContent)
}
