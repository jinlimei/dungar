package triggers

import (
	"fmt"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func bitCoinHandler(_, serverID string) string {
	return fmt.Sprintf(
		"Introducing %sCoin: A new, modern, sexy cryptocurrency to disrupt the %s industry.",
		utils.TitleCase(randomNoun(serverID)),
		randomNoun(serverID),
	)
}

func jsHandler(_, serverID string) string {
	return fmt.Sprintf(
		"Check out %s.js: A new framework that is going to revolutionize working with %s on %s",
		utils.TitleCase(randomNoun(serverID)),
		randomNoun(serverID),
		randomPlatform(serverID),
	)
}
