package dcli

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"gitlab.int.magneato.site/dungar/dungar/internal/db"
	"gitlab.int.magneato.site/dungar/dungar/internal/random"
	"gitlab.int.magneato.site/dungar/dungar/internal/utils"
	"gitlab.int.magneato.site/dungar/dungar/library/core2"
)

// ProtocolRunner is our one-stop-shop for initiating _some_ sort of protocol
// connection, be it Slack or Discord (rather than having a dedicated one for either)
func ProtocolRunner() {
	log.Printf("Starting Protocol Runner...")
	utils.LoadSettingsAndSecrets()

	log.Printf("Setting up RNG")
	random.UseTimeBasedSeed()

	log.Printf("Connecting to Database")
	db.ConnectToDatabase()

	var driver core2.ProtocolDriver

	log.Printf("Protocol Mode: %s", utils.ProtocolMode())

	switch utils.ProtocolMode() {
	case "discord":
		driver = DiscordRunner()
	default:
		log.Fatalf("FATAL: Unknown Protocol Mode '%s'!", utils.ProtocolMode())
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	<-signals

	if err := driver.Disconnect(); err != nil {
		log.Printf("ERROR: Attempted to disconnect but received error: %v", err)
	}

	os.Exit(0)
}
