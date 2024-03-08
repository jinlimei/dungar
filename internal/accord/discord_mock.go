package accord

import (
	"github.com/bwmarrin/discordgo"
)

// MockDiscordConnection provides a simulated way of working with discord connections
type MockDiscordConnection struct {
	handlers []DiscordHandlerFn
	Members  map[string]*discordgo.Member
	Channels map[string]*discordgo.Channel
}

// AddHandler adds a new handler function
func (m *MockDiscordConnection) AddHandler(handler DiscordHandlerFn) {
	if m.handlers == nil {
		m.handlers = make([]DiscordHandlerFn, 0)
	}

	m.handlers = append(m.handlers, handler)
}

// Start begins the mock connection (which does nothing and will never error)
func (m *MockDiscordConnection) Start(token string) error {
	return nil
}

// Connect "connects" the mock connection (does nothing and will never error)
func (m *MockDiscordConnection) Connect() error {
	return nil
}

// Disconnect "disconnects" the mock connection (does nothing and will never error)
func (m *MockDiscordConnection) Disconnect() error {
	return nil
}

// GetSession returns the mock session, always returns nil
func (m *MockDiscordConnection) GetSession() *discordgo.Session {
	return nil
}

// IsConnected checks if the mock connection is "connected" -- always returns false
func (m *MockDiscordConnection) IsConnected() bool {
	return false
}
