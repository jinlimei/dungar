package core2

import "errors"

var (
	// ErrChannelNotFound is when a channel was failed to be found (either because it
	// doesn't exist or it was an error)
	ErrChannelNotFound = errors.New("could not find channel with provided id")
)

// GetChannel will attempt to return the channel based on the given provided unique ID
// from the ProtocolDriver
func (s *Service) GetChannel(channelID, serverID string) (Channel, error) {
	return s.driver.GetChannel(channelID, serverID)
}

// GetChannelName will return only the channel name based on the provided unique ID
// from the ProtocolDriver
func (s *Service) GetChannelName(channelID, serverID string) string {
	return s.driver.GetChannelName(channelID, serverID)
}

// GetChannels will return all (currently) stored channels from the ProtocolDriver
func (s *Service) GetChannels(serverID string) map[string]Channel {
	return s.driver.GetChannels(serverID)
}
