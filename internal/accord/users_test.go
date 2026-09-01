package accord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestDisplayDiscordUserNamePrefersGlobalName(t *testing.T) {
	name := displayDiscordUserName(&discordgo.User{
		Username:   "Dungarmatic",
		GlobalName: "Furrymatic",
	})

	assert.Equal(t, "Furrymatic", name)
}

func TestGetUserNameUsesBotUserFallback(t *testing.T) {
	driver := initMockDriver()
	driver.SetBotUser(&core2.BotUser{
		ID:    "454311604847378442",
		Name:  "Furrymatic",
		IsBot: true,
	})

	assert.Equal(t, "Furrymatic", driver.GetUserName("454311604847378442", "mock"))
}
