package taut

import (
	"log"

	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// RealSlackConnection is our proper, real implementation of SlackConnection
type RealSlackConnection struct {
	rtm    *slack.RTM
	client *slack.Client

	users    map[string]*slack.User
	channels map[string]*slack.Channel
}

// NewRealSlackConnection provides a new instance of RealSlackConnection with the necessary
// members initialized. Necessary for having the RealSlackConnection work
func NewRealSlackConnection() *RealSlackConnection {
	return &RealSlackConnection{
		rtm:      nil,
		client:   nil,
		users:    make(map[string]*slack.User),
		channels: make(map[string]*slack.Channel),
	}
}

// Connect takes the token and establishes the RTM and Client APIs
func (r *RealSlackConnection) Connect(token string) error {
	api := slack.New(token)

	r.rtm = api.NewRTM()
	r.client = api

	return nil
}

// Disconnect attempts to disconnect from the RTM API
func (r *RealSlackConnection) Disconnect() error {
	return r.rtm.Disconnect()
}

// GetUserName attempts to retrieve a user name via the unique user ID
func (r *RealSlackConnection) GetUserName(id string) string {
	user, err := r.GetUser(id)

	if err != nil {
		return ""
	}

	return user.Name
}

// prepUserCache makes an initial effort to setup the user cache for user retrieval
func (r *RealSlackConnection) prepUserCache() {
	if len(r.users) > 0 {
		return
	}

	retrieved, err := r.client.GetUsers()
	utils.HaltingError("slack retrieveUsers", err)

	users := make(map[string]*slack.User)

	for k := 0; k < len(retrieved); k++ {
		user := &retrieved[k]

		users[user.ID] = user
	}

	r.users = users
}

// GetUser attempts to retrieve a slack.User via their unique ID
func (r *RealSlackConnection) GetUser(id string) (*slack.User, error) {
	r.prepUserCache()

	user, ok := r.users[id]

	if !ok {
		return nil, core2.ErrUserNotFound
	}

	return user, nil
}

// GetUsers returns the data in the user cache, which is a set of slack.User's
func (r *RealSlackConnection) GetUsers() map[string]*slack.User {
	r.prepUserCache()

	return r.users
}

// GetChannelName attempts to provide the current "channel" name for a given unique ID
// note that if it fails, the channel name will be '$CC$' + id
func (r *RealSlackConnection) GetChannelName(id string) string {
	channel, err := r.GetChannel(id)

	if err != nil {
		log.Printf("[ERROR] Failed to GetChannel(id=%s): %v", id, err)
		return "$CC$" + id
	}

	return channel.Name
}

func (r *RealSlackConnection) prepChannelCache() {
	if len(r.channels) > 0 {
		return
	}

	convos, cur, err := r.client.GetConversations(&slack.GetConversationsParameters{
		Cursor:          "",
		ExcludeArchived: true,
		Types:           []string{"public_channel", "private_channel", "mpim", "im"},
	})

	if err != nil {
		utils.HaltingError("retrieveChats 1", err)
	}

	channels := make([]slack.Channel, 0)
	channels = append(channels, convos...)

	for cur != "" {
		convos, cur, err = r.client.GetConversations(&slack.GetConversationsParameters{
			Cursor:          cur,
			ExcludeArchived: true,
			Types:           []string{"public_channel", "private_channel", "mpim", "im"},
		})

		if err != nil {
			utils.HaltingError("retrieveChats 2", err)
		}

		channels = append(channels, convos...)
	}

	output := make(map[string]*slack.Channel, len(channels))

	// Because we need to grab via pointers, we can't use a 'range' loop
	for k := 0; k < len(channels); k++ {
		channel := &channels[k]

		output[channel.ID] = channel
	}

	log.Printf("Added %d channels (group conversation things) to channel list",
		len(channels))

	r.channels = output
}

// GetChannel attempts to retrieve a given channel via its unique ID and returns a *slack.Channel
func (r *RealSlackConnection) GetChannel(id string) (*slack.Channel, error) {
	r.prepChannelCache()

	channel, ok := r.channels[id]

	if !ok {
		return nil, core2.ErrChannelNotFound
	}

	return channel, nil
}

// GetChannels returns the channel cache result
func (r *RealSlackConnection) GetChannels() map[string]*slack.Channel {
	r.prepChannelCache()

	return r.channels
}

// GetRtm returns the slack.RTM client
func (r *RealSlackConnection) GetRtm() *slack.RTM {
	return r.rtm
}

// GetClient returns the slack.Client
func (r *RealSlackConnection) GetClient() *slack.Client {
	return r.client
}

// IsConnected makes a soft attempt at determining whether or not we're connected.
func (r *RealSlackConnection) IsConnected() bool {
	return r.rtm != nil && r.client != nil
}
