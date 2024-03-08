package core2

// ResponseType indicates how a Response should be handled/interpreted by a given
// ProtocolDriver
//
//go:generate stringer -type=ResponseType
type ResponseType uint8

const (
	// ResponseTypeUnknown is when a ResponseType has not been set (defaults to 0)
	ResponseTypeUnknown ResponseType = iota
	// ResponseTypeBasic is a standard channel (or sub-channel) message
	ResponseTypeBasic
	// ResponseTypeReaction is a reaction emoji to a specific message
	ResponseTypeReaction
	// ResponseTypeDirectMessage will DM the individual it is responding to instead of
	// messaging the channel
	ResponseTypeDirectMessage
)
