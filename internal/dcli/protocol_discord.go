package dcli

import (
	"gitlab.int.magneato.site/dungar/dungar/internal/accord"
	"gitlab.int.magneato.site/dungar/dungar/internal/triggers"
	"gitlab.int.magneato.site/dungar/dungar/internal/utils"
	"gitlab.int.magneato.site/dungar/dungar/library/core2"
	"log"
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

	accordDriver.Connect(utils.DiscordAccessToken())

	return accordDriver
}
