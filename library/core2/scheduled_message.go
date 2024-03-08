package core2

import (
	"fmt"
	"strings"
	"time"
)

// ScheduledMessage is our structure on handling backlog
// which aren't in response to a IncomingMessage but instead
// generated arbitrarily (usually via time or something)
type ScheduledMessage struct {
	ChannelID        string
	ServerID         string
	ChannelArchetype string
	Cancelled        bool
	Contents         string
	SentAt           time.Time
}

// MakeScheduledMessage will create a basic ScheduledMessage
func MakeScheduledMessage(archetype, contents string, sentAt time.Time) *ScheduledMessage {
	return &ScheduledMessage{
		ChannelArchetype: archetype,
		Contents:         contents,
		SentAt:           sentAt,
	}
}

// IsEmpty returns whether, when trimming for space, the message is empty
func (sr *ScheduledMessage) IsEmpty() bool {
	return len(strings.TrimSpace(sr.Contents)) == 0
}

// IsCancelled returns whether the message is cancelled
func (sr *ScheduledMessage) IsCancelled() bool {
	return sr.Cancelled
}

// Cancel sets the message to cancelled
func (sr *ScheduledMessage) Cancel() *ScheduledMessage {
	sr.Cancelled = true
	return sr
}

// String provides some debug output for this message
func (sr *ScheduledMessage) String() string {
	return fmt.Sprintf(
		"channel-id='%s',channel-archetype='%s', cancelled='%v', contents='%s'",
		sr.ChannelID,
		sr.ChannelArchetype,
		sr.Cancelled,
		sr.Contents,
	)
}

// IsReady determines if all the special components for a scheduled message have been filled out
func (sr *ScheduledMessage) IsReady() bool {
	return !sr.Cancelled && sr.ChannelID != "" && sr.ServerID != "" && sr.Contents != ""
}

// ToResponse converts our simplified ScheduledMessage into a Response
func (sr *ScheduledMessage) ToResponse() *Response {
	return &Response{
		PrefixUsername:   false,
		ConsumedMessage:  true,
		HandledMessage:   true,
		CancelledMessage: sr.Cancelled,
		Contents:         sr.Contents,
		ResponseType:     ResponseTypeBasic,
		Origin:           OriginScheduled,
	}
}
