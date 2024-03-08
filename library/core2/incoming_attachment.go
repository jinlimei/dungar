package core2

import (
	"strings"
)

// AttachmentType is the different kind of types of attachments to a given message or event
//
//go:generate stringer -type=AttachmentType
type AttachmentType uint8

const (
	// AttachmentUnknown is our default unknown attachment type
	AttachmentUnknown AttachmentType = iota

	// AttachmentImage the attachment is a simple image
	AttachmentImage

	// AttachmentTweet the attachment is explicitly a tweet
	AttachmentTweet

	// AttachmentLinkUnfurl is when the attachment is an unfurled link. Could be reddit or
	// many other websites where this sort of thing is standardized
	AttachmentLinkUnfurl
)

// IncomingAttachment is our standardized/generic struct for message or event attachments
type IncomingAttachment struct {
	AuthorName  string
	AuthorURL   string
	Title       string
	TitleLink   string
	ImageURL    string
	ServiceName string
	Fallback    string
	Type        AttachmentType
}

// String attempts to provide a debug-string output for the given IncomingAttachment
func (ia *IncomingAttachment) String() string {
	out := make([]string, 0)

	out = append(out, "```")

	if ia.AuthorName != "" {
		out = append(out, "Author: "+ia.AuthorName)
	}

	if ia.Title != "" {
		out = append(out, "Title: "+ia.Title)
	}

	if ia.TitleLink != "" {
		out = append(out, "(Title) Link: Yes")
	}

	if ia.ServiceName != "" {
		out = append(out, "Service Name: "+ia.ServiceName)
	}

	if ia.Fallback != "" {
		out = append(out, "Fallback: "+ia.Fallback)
	}

	out = append(out, "Type: "+ia.Type.String())
	out = append(out, "```")

	return strings.Join(out, "\n")
}
