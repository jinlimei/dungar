package core2

import (
	"log"
	"strings"
)

// EventResponseEnvelope is our response envelope with the event and its associated
// responses.
type EventResponseEnvelope struct {
	Event     *IncomingEvent
	Responses []*EventResponse
}

// EventResponse is our response to a given event within a specific location
type EventResponse struct {
	// PrefixUsername if true (and the event *does* have a valid user associated
	// with it) the message will be prefixed with a ping to said user
	PrefixUsername bool

	// ConsumedEvent if true, means that no other event handlers can respond
	// to this event
	ConsumedEvent bool

	// Contents is the response given from the bot to the channel
	Contents string
}

// IsEmpty returns true if the event contents is empty (after a strings.TrimSpace)
func (er *EventResponse) IsEmpty() bool {
	return strings.TrimSpace(er.Contents) == ""
}

// Format will format the event text to the appropriately expected protocol handling
func (er *EventResponse) Format(driver ProtocolDriver, event *IncomingEvent) string {
	if er.PrefixUsername {
		user, err := driver.GetUser(event.UserID, event.ServerID)

		if err != nil {
			log.Printf("[EventResponse.Format] Could not get user with ID '%s': %v", event.UserID, err)
			return er.Contents
		}

		return driver.PingUser(user) + er.Contents
	}

	return er.Contents
}

// ToResponse converts EventResponse to a Response
func (er *EventResponse) ToResponse() *Response {
	return &Response{
		PrefixUsername: er.PrefixUsername,
		Contents:       er.Contents,
		ResponseType:   ResponseTypeBasic,
		Origin:         OriginEvent,
	}
}
