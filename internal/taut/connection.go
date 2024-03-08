package taut

import (
	"github.com/slack-go/slack"
)

// SlackConnection is our interface for all the expected functions
// using the slack api.
type SlackConnection interface {
	Connect(token string) error
	Disconnect() error

	GetUserName(id string) string
	GetUser(id string) (*slack.User, error)
	GetUsers() map[string]*slack.User

	GetChannelName(id string) string
	GetChannel(id string) (*slack.Channel, error)
	GetChannels() map[string]*slack.Channel

	GetRtm() *slack.RTM
	GetClient() *slack.Client

	IsConnected() bool
}

// Connect will attempt to make a connection to a slack server
// and establish an event handling loop
func (d *Driver) Connect(token string) {
	_ = d.Con.Connect(token)

	go d.Con.GetRtm().ManageConnection()
	go d.incomingEventLoop()
	go d.outgoingEventLoop()
}

// Disconnect will disconnect from our RTM slack connection
func (d *Driver) Disconnect() error {
	return d.Con.Disconnect()
}
