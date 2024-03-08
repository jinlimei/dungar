package taut

import (
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// MockSlackConnection is our mock methodology for having a slack connection
type MockSlackConnection struct {
	Users    map[string]*slack.User
	Channels map[string]*slack.Channel
}

// NewMockSlackConnection will make a new standardized MockSlackConnection
func NewMockSlackConnection() *MockSlackConnection {
	return &MockSlackConnection{
		Users:    make(map[string]*slack.User),
		Channels: make(map[string]*slack.Channel),
	}
}

// SetUser will create a slack.User
func (m *MockSlackConnection) SetUser(id string, name string) {
	m.Users[id] = &slack.User{
		ID:       id,
		Name:     name,
		RealName: name,
		Profile: slack.UserProfile{
			DisplayName:           name,
			DisplayNameNormalized: name,
		},
	}
}

// SetChannel will create a slack.Channel
func (m *MockSlackConnection) SetChannel(id string, name string) {
	m.Channels[id] = &slack.Channel{
		GroupConversation: slack.GroupConversation{
			Conversation: slack.Conversation{
				ID:                 id,
				NameNormalized:     name,
				IsGroup:            false,
				IsShared:           false,
				IsIM:               false,
				IsExtShared:        false,
				IsOrgShared:        false,
				IsPendingExtShared: false,
				IsPrivate:          false,
				IsMpIM:             false,
				User:               "",
			},
			Name:    name,
			Topic:   slack.Topic{},
			Purpose: slack.Purpose{},
		},
		IsChannel: false,
		IsGeneral: false,
		IsMember:  false,
	}
}

// Connect necessary for SlackConnection
func (m *MockSlackConnection) Connect(token string) error {
	return nil
}

// Disconnect necessary for SlackConnection
func (m *MockSlackConnection) Disconnect() error {
	return nil
}

// GetUserName necessary for SlackConnection
func (m *MockSlackConnection) GetUserName(id string) string {
	user, ok := m.Users[id]

	if ok {
		return user.Name
	}

	return ""
}

// GetUser necessary for SlackConnection
func (m *MockSlackConnection) GetUser(id string) (*slack.User, error) {
	user, ok := m.Users[id]

	if !ok {
		return nil, core2.ErrUserNotFound
	}

	return user, nil
}

// GetUsers necessary for SlackConnection
func (m *MockSlackConnection) GetUsers() map[string]*slack.User {
	return m.Users
}

// GetChannelName necessary for SlackConnection
func (m *MockSlackConnection) GetChannelName(id string) string {
	channel, ok := m.Channels[id]

	if ok {
		return channel.Name
	}

	return ""
}

// GetChannel necessary for SlackConnection
func (m *MockSlackConnection) GetChannel(id string) (*slack.Channel, error) {
	channel, ok := m.Channels[id]

	if !ok {
		return nil, core2.ErrChannelNotFound
	}

	return channel, nil
}

// GetChannels necessary for SlackConnection
func (m *MockSlackConnection) GetChannels() map[string]*slack.Channel {
	return m.Channels
}

// GetRtm necessary for SlackConnection (and does nothing)
func (m *MockSlackConnection) GetRtm() *slack.RTM {
	return nil
}

// GetClient necessary for SlackConnection (and does nothing)
func (m *MockSlackConnection) GetClient() *slack.Client {
	return nil
}

// IsConnected necessary for SlackConnection (and is always true)
func (m *MockSlackConnection) IsConnected() bool {
	return true
}
