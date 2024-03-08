package accord

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func makeDiscordMessage(str string) *discordgo.Message {
	return &discordgo.Message{
		Content: str,
	}
}

func TestBasicMessageParsing(t *testing.T) {
	driver := initMockDriver()
	parsed := parseDiscordMessage(makeDiscordMessage("\u003c@569251658773692538\u003e a or b?"))
	parsed.Converted = driver.translateParsedMessage("mock", parsed)

	assert.Equal(t, 7, len(parsed.Tokens))
	assert.Equal(t, 1, len(parsed.IDTokens()))
	assert.Equal(t, "@Dungar a or b?", parsed.Converted)
}

func TestAllKnownBracketThings(t *testing.T) {
	driver := initMockDriver()
	parsed := parseDiscordMessage(makeDiscordMessage(
		"here's <@62657919375646720> or maybe (<@&574140080751378432>) <@569251658773692538>? how about <@&766488443248836609> (should be role) or @everyone (should be everyone) or @here",
	))
	parsed.Converted = driver.translateParsedMessage("mock", parsed)

	assert.Equal(t, 40, len(parsed.Tokens))
	assert.Equal(t, 4, len(parsed.IDTokens()))
	assert.Equal(t, "here's @jinli.mei or maybe (&Dungar) @Dungar? how about &example (should be role) or @everyone (should be everyone) or @here", parsed.Converted)
}

func TestWeirdKryptnMessage(t *testing.T) {
	driver := initMockDriver()
	parsed := parseDiscordMessage(makeDiscordMessage(
		"\u003c:nice:1112827243341815849\u003e  https://i.arq.dev/C9RMaAU.png",
	))

	parsed.Converted = driver.translateParsedMessage("mock", parsed)

	spew.Dump(parsed)
}
