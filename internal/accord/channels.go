package accord

import (
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
)

// GetChannel retrieves the core2.Channel from its unique ID
func (d *Driver) GetChannel(channelID, serverID string) (core2.Channel, error) {
	channels := d.getCoreChannels()

	for _, channel := range channels {
		if channel.ID == channelID {
			return channel, nil
		}
	}

	return core2.Channel{}, core2.ErrChannelNotFound
}

// GetChannels returns channels with the ID as a key in a map
func (d *Driver) GetChannels(serverID string) map[string]core2.Channel {
	guild := d.getOrMakeGuild(serverID)
	converted := make(map[string]core2.Channel, len(guild.channelCache))

	for _, channel := range guild.channelCache {
		converted[channel.ID] = d.translateDiscordChannel(channel)
	}

	return converted
}

// GetChannelName returns the channel name from its ID
func (d *Driver) GetChannelName(channelID, serverID string) string {
	channel, err := d.GetChannel(channelID, serverID)

	if err != nil {
		log.Printf("ERROR GetChannelName('%s') = %v", channelID, err)

		return "@" + channel.ID
	}

	return "@" + channel.Name
}
