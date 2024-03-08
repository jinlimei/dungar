package triggers

import (
	"log"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

var (
	mockService *core2.Service
	mockDriver  *core2.MockProtocolDriver
)

func initMockServices() *core2.Service {
	mockService, mockDriver = core2.NewMockService()
	core = mockService

	mockDriver.SetChannel("politics", "politics", "public", core2.ChannelPublic)
	mockDriver.SetChannel("dungar-test", "dungar-test", "testing", core2.ChannelPublic)
	mockDriver.SetChannel("butts", "butts", "public", core2.ChannelPublic)
	mockDriver.SetChannel("bar", "bar", "public", core2.ChannelPublic)

	mockDriver.SetBotUser(core2.BotUser{
		ID:      "dungar",
		Name:    "Dungar",
		IsBot:   true,
		IsAdmin: false,
	})

	mockDriver.Users["dungar"] = core2.User{
		ID:      "dungar",
		Name:    "Dungar",
		IsBot:   true,
		IsAdmin: false,
	}

	mockDriver.SetUser("foo", "Foo")
	mockDriver.SetUser("butts", "Buttes")
	mockDriver.SetUser("kinkouin", "Kinkouin")

	return mockService
}

func pickWeightedChoice(choices []weightedChoice) string {
	totalWeight := 0.0

	for _, choice := range choices {
		totalWeight += choice.Chance
	}

	randomWeight := random.Float64() * totalWeight
	accruedWeight := 0.0

	for _, choice := range choices {
		accruedWeight += choice.Chance

		if accruedWeight >= randomWeight {
			return choice.Response
		}
	}

	// this should never happen
	log.Fatal("a thing that should never happen")
	return choices[0].Response
}

func pickRandomWord(str string, normalize bool) string {
	words := utils.StringToWords(str, normalize)

	return random.PickString(words)
}

func haltingErr(loc string, err error) {
	utils.HaltingError("triggers haltingErr"+loc, err)
}
