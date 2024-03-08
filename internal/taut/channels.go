package taut

import (
	"log"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// GetChannel returns a channel struct and if it was a successful grab
func (d *Driver) GetChannel(channelID, _ string) (core2.Channel, error) {
	if channel, ok := d.chanCache[channelID]; ok {
		return d.translateSlackChannel(channel), nil
	}

	return d.getChannelReal(channelID)
}

// GetChannels returns a map of all channels (with their channel IDs as the key)
func (d *Driver) GetChannels(_ string) map[string]core2.Channel {
	return d.getChannelsReal()
}

// GetChannelName converts the channel id to a channels name
func (d *Driver) GetChannelName(channelID, serverID string) string {
	channel, err := d.GetChannel(channelID, serverID)

	if err != nil {
		log.Printf("GetChannelName('%s'): %v", channelID, err)

		return ""
	}

	return channel.Name
}
