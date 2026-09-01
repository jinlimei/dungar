package dcli

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/accord"
	"gitlab.int.magneato.site/dungar/prototype/internal/triggers"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
	"sort"
)

// DiscordRunner sets up the necessary things to run discord bbygirl
func DiscordRunner() core2.ProtocolDriver {
	con := accord.NewRealDiscordConnection()

	accordDriver, err := accord.New(con)

	if err != nil {
		log.Fatalf("Failed to create new discord driver: %v", err)
	}

	coreSvc := core2.New(accordDriver)

	triggers.RegisterHandlers(coreSvc)

	allowedOutputChannels := utils.DiscordAllowedOutputChannelIDs()
	if len(allowedOutputChannels) == 0 {
		log.Printf("Discord channel allowlist (outgoing + learning): unrestricted")
	} else {
		channelIDs := make([]string, 0, len(allowedOutputChannels))

		for channelID := range allowedOutputChannels {
			channelIDs = append(channelIDs, channelID)
		}

		sort.Strings(channelIDs)
		log.Printf("Discord channel allowlist (outgoing + learning): %v", channelIDs)
	}

	log.Printf("Discord learning window: last %d days", utils.DiscordLearningLookbackDays)
	triggers.PreloadMarkovOnStartup()
	triggers.InitializeManicMinute()

	accordDriver.Connect(utils.DiscordAccessToken())

	return accordDriver
}
