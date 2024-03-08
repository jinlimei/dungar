package taut

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) getChannelReal(id string) (core2.Channel, error) {
	slackChannel, err := d.Con.GetChannel(id)

	if err != nil {
		return core2.Channel{}, err
	}

	log.Printf("getChannelReal(%s) = ID=%s, Name=%s, Type=%s", id, slackChannel.ID, slackChannel.Name, d.inferChannelType(slackChannel))

	if slackChannel.ID != id {
		return core2.Channel{}, fmt.Errorf("failed to retrieve slack channel '%s': received wrong channel '%s'",
			id, slackChannel.ID)
	}

	d.chanCache[id] = slackChannel

	return d.translateSlackChannel(slackChannel), nil
}

func (d *Driver) getChannelsReal() map[string]core2.Channel {
	channels := d.Con.GetChannels()

	out := make(map[string]core2.Channel, len(channels))

	for _, channel := range channels {
		d.chanCache[channel.ID] = channel
		out[channel.ID] = d.translateSlackChannel(channel)
	}

	return out
}

// inferChannelType takes the insane slack.Channel standard and attempts to infer a
// generalized channel type across protocols
func (d *Driver) inferChannelType(channel *slack.Channel) core2.ChannelType {
	if channel.IsMpIM {
		return core2.ChannelGroupMessage
	}

	if channel.IsIM {
		return core2.ChannelDirectMessage
	}

	// IsGroup is only true if it'd a private channel created before March 2021.
	// If it'd a private channel afterwards, it'd IsChannel=true, IsPrivate=true
	if channel.IsGroup || (channel.IsChannel && channel.IsPrivate) {
		return core2.ChannelPrivileged
	}

	return core2.ChannelPublic
}

func (d *Driver) getChannelName(channel *slack.Channel) string {
	if channel.Name != "" {
		return channel.Name
	}

	switch d.inferChannelType(channel) {
	case core2.ChannelGroupMessage:
		return "&" + channel.ID
	case core2.ChannelDirectMessage:
		// A direct message, the channel'd ID should be the user of the message
		user, err := d.GetUser(channel.User, "")

		if err != nil {
			return "@" + user.Name
		}

		log.Printf(
			"[CORE-C000] Unable to find user '%s' for DM channel, using channel ID '%s' instead",
			channel.User,
			channel.ID,
		)

		// If we failed, we'll just use the ID
		return "@" + channel.ID
	}

	// We've hella failed somehow
	spew.Dump(channel)
	panic(fmt.Sprintf("Failed to get channel name for channel '%s'", channel.ID))
}

func (d *Driver) translateSlackChannel(channel *slack.Channel) core2.Channel {
	return core2.Channel{
		ID:             channel.ID,
		ServerID:       d.teamID,
		Name:           channel.Name,
		NameNormalized: channel.NameNormalized,
		PreviousNames:  make([]string, 0),
		Topic:          channel.Topic.Value,
		Type:           d.inferChannelType(channel),
	}
}
