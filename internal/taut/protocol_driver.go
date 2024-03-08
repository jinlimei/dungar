package taut

import (
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// Driver is the slack-specific implementation of core2.ProtocolDriver
type Driver struct {
	// Con the implementation of the SlackConnection!
	Con SlackConnection

	botUser  *core2.BotUser
	outgoing *core2.OutgoingResponses
	core     *core2.Service
	teamID   string

	userCache map[string]*slack.User
	chanCache map[string]*slack.Channel
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "slack"
}

// SetCore allows driver to have access to the core service
func (d *Driver) SetCore(svc *core2.Service) {
	d.core = svc
}

// New provides a (safe) means for a new instance of Driver
// It isn't explicitly necessary as all important members interacted
// with here are also public (or usable via Driver.InitCache)
func New(con SlackConnection) *Driver {
	d := &Driver{
		Con: con,
	}

	d.InitCache()

	return d
}

// InitCache will initialize the userCache and chanCache members
func (d *Driver) InitCache() {
	d.userCache = make(map[string]*slack.User, 70)
	d.chanCache = make(map[string]*slack.Channel, 70)
}
