package accord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

// DiscordHandlerFn is the signature for discord handler functions
// type DiscordHandlerFn func(s *discordgo.Session, ev any)
type DiscordHandlerFn interface{}

// DiscordConnection is the interface to be based on for RealDiscordConnection
// and for a mock connection
type DiscordConnection interface {
	AddHandler(handler DiscordHandlerFn)

	Start(token string) error
	Connect() error
	Disconnect() error

	GetSession() *discordgo.Session

	IsConnected() bool
}

// Connect opens a session with discord
func (d *Driver) Connect(token string) {
	err := d.Con.Start(token)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	d.registerHandlers()

	err = d.Con.Connect()

	if err != nil {
		log.Fatalf("Failed to open session: %v", err)
	}

	go d.outgoingEventLoop()
}

// Disconnect attempts to disconnect from discord
func (d *Driver) Disconnect() error {
	return d.Con.Disconnect()
}
