package core2

// ProtocolDriver is a standard way of handling chat connections.
// The hopeful goal is to have Dungar support multiple platforms
type ProtocolDriver interface {
	DriverName() string
	GetUserName(userID, serverID string) string
	GetUser(userID, serverID string) (User, error)
	GetUsers(serverID string) map[string]User

	GetBotUser() BotUser

	GetChannelName(channelID, serverID string) string
	GetChannel(channelID, serverID string) (Channel, error)
	GetChannels(serverID string) map[string]Channel

	PingUser(user User) string

	SetCore(svc *Service)
	Disconnect() error
}
