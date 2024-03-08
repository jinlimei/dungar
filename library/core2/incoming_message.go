package core2

import (
	"strings"
)

// MessageType is the different types of incoming backlog
//
//go:generate stringer -type=MessageType
type MessageType uint8

const (
	// MessageTypeUnknown is for unknown message types
	MessageTypeUnknown MessageType = iota
	// MessageTypeBasic is a standard, basic message
	MessageTypeBasic
	// MessageTypeMe is a /me message
	MessageTypeMe
	// MessageTypeChanged is a change of a previous message
	MessageTypeChanged
	// MessageTypeDeleted is a deleted message
	// This will also be an event from IncomingEvent
	MessageTypeDeleted
	// MessageTypeBroadcast is a broadcast-specific message (?)
	MessageTypeBroadcast
)

// IncomingMessage is our standardized incoming message type. Replaces the ConsumableMessage
// from Ye Olden Days.
type IncomingMessage struct {
	ID       string
	UserID   string
	ServerID string

	ChannelID    string
	SubChannelID string
	IsSubMessage bool

	IsReply   bool
	ReplyToID string

	Contents       string
	LowerContents  string
	ParsedContents *ParsedMessage

	Type        MessageType
	Attachments []IncomingAttachment
}

// String returns the message Contents
func (im *IncomingMessage) String() string {
	return im.Contents
}

// Lowered returns the message Contents but in all lower case.
// This will also set LowerContents if not already set
func (im *IncomingMessage) Lowered() string {
	if im.LowerContents == "" && im.Contents != "" {
		im.LowerContents = strings.ToLower(im.Contents)
	}

	return im.LowerContents
}

// HasValidUser returns whether or not the message has come with a valid UserID
func (im *IncomingMessage) HasValidUser() bool {
	return im.UserID != ""
}
