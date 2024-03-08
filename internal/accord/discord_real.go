package accord

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

var (
	// ErrConnectNotCalled is to provide a checkable error when we attempt to do things before we ever connected
	ErrConnectNotCalled = errors.New("connect has not been called (no valid session)")
)

// RealDiscordConnection is the real implementation for the interface DiscordConnection
// and manages the interaction with the discordgo library
type RealDiscordConnection struct {
	session   *discordgo.Session
	connected bool
}

// NewRealDiscordConnection starts an instance of RealDiscordConnection
func NewRealDiscordConnection() *RealDiscordConnection {
	return &RealDiscordConnection{}
}

// IsConnected is whether or not Connect has been called
func (r *RealDiscordConnection) IsConnected() bool {
	return r.connected
}

// Start takes the token and builds a new session.
// Note that this *does not* start and open the session
func (r *RealDiscordConnection) Start(token string) error {
	dgo, err := discordgo.New("Bot " + token)

	if err != nil {
		return err
	}

	dgo.LogLevel = discordgo.LogInformational

	r.session = dgo
	r.connected = true

	return nil
}

// Connect opens a session
func (r *RealDiscordConnection) Connect() error {
	return r.session.Open()
}

// Disconnect closes a session
func (r *RealDiscordConnection) Disconnect() error {
	return r.session.Close()
}

// AddHandler is a wrapper around the AddHandler function in the discordgo.Session
func (r *RealDiscordConnection) AddHandler(handler DiscordHandlerFn) {
	//log.Printf("RealDiscordConnection AddHandler, session: %+v", r.session)
	r.session.AddHandler(handler)
}

// GetSession returns the session that was built when Start is called
func (r *RealDiscordConnection) GetSession() *discordgo.Session {
	return r.session
}
