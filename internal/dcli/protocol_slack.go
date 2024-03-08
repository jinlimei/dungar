package dcli

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/taut"
	"gitlab.int.magneato.site/dungar/prototype/internal/triggers"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// SlackRunner is our entrypoint for running Slack connections with Dungar
func SlackRunner() core2.ProtocolDriver {
	con := taut.NewRealSlackConnection()

	tautDriver := taut.New(con)

	coreSvc := core2.New(tautDriver)

	triggers.RegisterHandlers(coreSvc)

	tautDriver.Connect(utils.SlackAccessToken())

	return tautDriver
}
