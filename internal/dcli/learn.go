package dcli

import (
	"fmt"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// Learn is a protocol-aware CLI entrypoint for loading Markov data from the
// configured raw message source. This is a standalone process and does not
// mutate the in-memory Markov state of any already-running dungar process.
func Learn() {
	utils.LoadSettingsAndSecrets()
	db.ConnectToDatabase()

	mode, err := normalizeLearnProtocolMode(utils.ProtocolMode())
	if err != nil {
		log.Fatal(err)
	}

	switch mode {
	case "discord":
		learnDiscord()
	case "slack":
		log.Fatal("CLI learn path for Slack is not implemented yet")
	default:
		log.Fatalf("Unhandled learn mode '%s'", mode)
	}
}

func normalizeLearnProtocolMode(mode string) (string, error) {
	switch {
	case strings.EqualFold(mode, "discord"):
		return "discord", nil
	case strings.EqualFold(mode, "slack"):
		return "slack", nil
	default:
		return "", fmt.Errorf("unknown protocol mode for learn: %s", mode)
	}
}

func learnDiscord() {
	markov := markov3.MakeMarkov("cli-learn")
	learned := markov.LearnFromRawMessages()

	log.Printf(
		"CLI learn complete: mode=discord learned=%d words=%d fragments=%d fragment_words=%d",
		learned,
		len(markov.RevWords),
		len(markov.Fragments),
		len(markov.FragmentWords),
	)
	log.Printf(
		"CLI learn note: this standalone command does not update any already-running dungar process; use !m3 or !markov load-m3 to refresh the live bot state",
	)
}
