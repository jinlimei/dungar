package core2

import "fmt"

// ChannelType is the means of establishing special channel types for working out how
// to react to them in various circumstances (through triggers, etc.)
//
//go:generate stringer -type=ChannelType
type ChannelType uint8

const (
	// ChannelUnknown is our base channel type and is used as the "default" if no code
	// actually specified a valid type
	ChannelUnknown ChannelType = iota
	// ChannelPublic is a standard slack/discord/IRC channel
	ChannelPublic
	// ChannelPrivileged is a private Slack channel, or other instances of a channel
	// which could conceivably be private
	ChannelPrivileged
	// ChannelDirectMessage is a 1-on-1 Direct Message between the bot and a user
	ChannelDirectMessage
	// ChannelGroupMessage is a multi-person direct message between the bot and users
	ChannelGroupMessage

	// ChannelReadOnly is for channels whose entire existence will _only_ be read-only
	// for the bot. This generally includes channels which may not actually be read-only
	// for admins or other users.
	ChannelReadOnly

	// ChannelVoice is a type of voice channel! Who knew ~
	ChannelVoice

	// ChannelForum handles the very special situation on Discord where
	// there are pseudo-channels which also behave like forums.
	ChannelForum
)

// Channel is our standardized struct for working with group chats across
// various protocols (primarily Slack)
type Channel struct {
	ID             string      `json:"id"`
	ServerID       string      `json:"server_id"`
	Name           string      `json:"name"`
	NameNormalized string      `json:"name_normalized"`
	PreviousNames  []string    `json:"previous_names"`
	Topic          string      `json:"topic"`
	Type           ChannelType `json:"type"`
	CanPost        bool        `json:"can_post"`
	Archetype      string      `json:"archetype"`
}

// String returns the Channel in a string-debugish format
func (c *Channel) String() string {
	if c == nil {
		return "core.Channel(<nil>)"
	}

	return fmt.Sprintf(
		"Channel(ID='%s',Name='%s',NameNormalized='%s',Type='%s'",
		c.ID,
		c.Name,
		c.NameNormalized,
		c.Type.String(),
	)
}
