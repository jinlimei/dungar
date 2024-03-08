package core2

// EventType is our different kinds of event types for a given IncomingEvent
//
//go:generate stringer -type=EventType
type EventType uint8

const (
	// EventUnknown is the standard default event
	EventUnknown EventType = iota
	// EventChannelJoin is when a user joins a Channel
	EventChannelJoin
	// EventChannelLeave is when a user leaves a channel
	EventChannelLeave
	// EventChannelTopic is when a user changes the channel topic
	EventChannelTopic
	// EventChannelPurpose is when a user changes the channel purpose
	EventChannelPurpose
	// EventChannelName is when a user changes the channel name
	EventChannelName
	// EventChannelArchive is when a user archives the channel
	EventChannelArchive
	// EventFileShare is when a user shares a file
	EventFileShare
	// EventPinnedItem is when a user pins an item
	EventPinnedItem
	// EventUnpinnedItem is when a user unpins an item
	EventUnpinnedItem
	// EventReactionAdd is when a reaction is added to a message
	EventReactionAdd
	// EventReactionRemove is when a reaction is removed from a message
	EventReactionRemove
	// EventMessageDelete is when a message is/was deleted
	EventMessageDelete
)

// IsChannelEvent refers to an event for the specific state of a channel
func (et EventType) IsChannelEvent() bool {
	return et == EventChannelJoin || et == EventChannelLeave || et == EventChannelTopic ||
		et == EventChannelPurpose || et == EventChannelName || et == EventChannelArchive
}

// IncomingEvent is our standardized struct for incoming non-message events.
type IncomingEvent struct {
	ID          string
	UserID      string
	ServerID    string
	ChannelID   string
	Text        string
	Type        EventType
	Attachments []IncomingAttachment
}
