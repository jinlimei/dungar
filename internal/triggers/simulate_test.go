package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestSimulateUserHandlerDirectedMention(t *testing.T) {
	svc := initMockServices()

	savedBuilder := userSimulationBuilder
	userSimulationBuilder = func(target core2.User, serverID string) (string, int, error) {
		assert.Equal(t, "peter", target.ID)
		assert.Equal(t, "Peter", target.Name)
		assert.Equal(t, "arena", serverID)
		return "peter says hello there", 3, nil
	}
	t.Cleanup(func() {
		userSimulationBuilder = savedBuilder
	})

	mockDriver.SetBotUser(core2.BotUser{
		ID:      "dungar",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	})
	mockDriver.Users["dungar"] = core2.User{
		ID:      "dungar",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	}
	mockDriver.SetUser("peter", "Peter")

	botID := "dungar"
	targetID := "peter"
	msg := &core2.IncomingMessage{
		UserID:    "bob",
		ServerID:  "arena",
		ChannelID: "arena",
		Contents:  "@Dungarmatic simulate @Peter",
		ParsedContents: &core2.ParsedMessage{
			Tokens: []core2.MessageToken{
				{Token: "<@dungar>", Type: core2.TokenUserID, Value: &botID},
				{Token: " ", Type: core2.TokenSpace},
				{Token: "simulate", Type: core2.TokenWord},
				{Token: " ", Type: core2.TokenSpace},
				{Token: "<@peter>", Type: core2.TokenUserID, Value: &targetID},
			},
		},
	}

	rsp := simulateUserHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.Equal(t, "peter says hello there", rsp[0].Contents)
}

func TestSimulateUserHandlerUsage(t *testing.T) {
	svc := initMockServices()
	msg := makeMessage("@Dungar simulate", "bob", "arena")
	msg.ServerID = "arena"

	rsp := simulateUserHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.Contains(t, rsp[0].Contents, "simulate @user")
}

func TestSimulateUserHandlerSelfUsesDungarmaticVoice(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)

	svc := initMockServices()

	mockDriver.SetBotUser(core2.BotUser{
		ID:      "dungar",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	})
	mockDriver.Users["dungar"] = core2.User{
		ID:      "dungar",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	}

	botID := "dungar"
	targetID := "dungar"
	msg := &core2.IncomingMessage{
		UserID:    "bob",
		ServerID:  "arena",
		ChannelID: "arena",
		Contents:  "@Dungarmatic simulate @Dungarmatic",
		ParsedContents: &core2.ParsedMessage{
			Tokens: []core2.MessageToken{
				{Token: "<@dungar>", Type: core2.TokenUserID, Value: &botID},
				{Token: " ", Type: core2.TokenSpace},
				{Token: "simulate", Type: core2.TokenWord},
				{Token: " ", Type: core2.TokenSpace},
				{Token: "<@dungar>", Type: core2.TokenUserID, Value: &targetID},
			},
		},
	}

	rsp := simulateUserHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.NotEmpty(t, rsp[0].Contents)
	assert.NotContains(t, rsp[0].Contents, "don't know enough")
}
