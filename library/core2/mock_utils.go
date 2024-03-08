package core2

import "fmt"

// NewMockService will create a new Service with a MockProtocolDriver associated
// with it for testing purposes
func NewMockService() (*Service, *MockProtocolDriver) {
	mock := &MockProtocolDriver{
		Core:     nil,
		Users:    make(map[string]User),
		Channels: make(map[string]Channel),
		BotUser:  BotUser{},
	}

	mock.Core = New(mock)

	return mock.Core, mock
}

// MockProtocolDriver is our mock implementation of a protocol driver
type MockProtocolDriver struct {
	Core     *Service
	Users    map[string]User
	Channels map[string]Channel
	BotUser  BotUser
}

// DriverName returns "mock"
func (mpd *MockProtocolDriver) DriverName() string {
	return "mock"
}

// SetBotUser will set the BotUser
func (mpd *MockProtocolDriver) SetBotUser(user BotUser) {
	mpd.BotUser = user
}

// SetUser will create a new User
func (mpd *MockProtocolDriver) SetUser(id, name string) {
	mpd.Users[id] = User{
		ID:      id,
		Name:    name,
		IsBot:   false,
		IsAdmin: false,
	}
}

// SetChannel will create a new Channel
func (mpd *MockProtocolDriver) SetChannel(id, name, archetype string, chanType ChannelType) {
	mpd.Channels[id] = Channel{
		ID:             id,
		Name:           name,
		NameNormalized: name,
		Archetype:      archetype,
		Type:           chanType,
	}
}

// GetUserName is our ProtocolDriver.GetUserName implementation
func (mpd *MockProtocolDriver) GetUserName(userID, serverID string) string {
	user, ok := mpd.Users[userID]

	if !ok {
		return ""
	}

	return user.Name
}

// GetUser is our ProtocolDriver.GetUser implementation
func (mpd *MockProtocolDriver) GetUser(userID, serverID string) (User, error) {
	user, ok := mpd.Users[userID]

	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}

// GetUsers is our ProtocolDriver.GetUsers implementation
func (mpd *MockProtocolDriver) GetUsers(serverID string) map[string]User {
	return mpd.Users
}

// GetBotUser is our ProtocolDriver.GetBotUser implementation
func (mpd *MockProtocolDriver) GetBotUser() BotUser {
	return mpd.BotUser
}

// GetChannelName is our ProtocolDriver.GetChannelName implementation
func (mpd *MockProtocolDriver) GetChannelName(channelID, serverID string) string {
	channel, ok := mpd.Channels[channelID]

	if !ok {
		return ""
	}

	return channel.Name
}

// GetChannel is our ProtocolDriver.GetChannel implementation
func (mpd *MockProtocolDriver) GetChannel(channelID, serverID string) (Channel, error) {
	channel, ok := mpd.Channels[channelID]

	if !ok {
		return Channel{}, ErrChannelNotFound
	}

	return channel, nil
}

// GetChannels is our ProtocolDriver.GetChannels implementation
func (mpd *MockProtocolDriver) GetChannels(serverID string) map[string]Channel {
	return mpd.Channels
}

// PingUser is our ProtocolDriver.PingUser implementation
func (mpd *MockProtocolDriver) PingUser(user User) string {
	return fmt.Sprintf("<@%s>", user.ID)
}

// SetCore is our ProtocolDriver.SetCore implementation
func (mpd *MockProtocolDriver) SetCore(svc *Service) {
	mpd.Core = svc
}

// Disconnect will do nothing for this mock driver
func (mpd *MockProtocolDriver) Disconnect() error {
	return nil
}
