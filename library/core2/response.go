package core2

import (
	"fmt"
	"log"
	"strings"
)

// Response is a struct of a message and whether or not to prefix
// that message with the respondents name.
//
// The table for what means what is as follows:
//
//   - CancelledMessage flags the message as HandledMessage & ConsumedMessage.
//     Regardless of how these other flags are set, CancelledMessage is an
//     overriding flag. This forces the response to get into the envelope slice
//     such that when we need to determine if the ConsumableMessage is bad, it's
//     always going to be present to say TRUE IT'S BAD DELETE
//
//   - ConsumedMessage lets the message handler loop know that this message has
//     been consumed, even if it hasn't been handled. This likely means that there
//     is a soft-stopping in the handler loop as opposed to a hard-stop like
//     CancelledMessage would bring.
//
//   - HandledMessage lets the message handler loop know that the Response
//     should be added to the outgoing slice to make an envelope.
type Response struct {
	// PrefixUsername when true will force the message in the response
	// to be built with <@user> Contents -- and the syntax for mention
	// will be chat-protocol dependent.
	PrefixUsername bool

	// MakeSubChannelIfPossible will indicate to the outgoing messaging to
	// make a sub-channel (slack thread or discord response) if it doesn't
	// already have one.
	MakeSubChannelIfPossible bool

	// ConsumedMessage means that, when handled, the message will not respond
	// to any further MessageEvHandler's and stop at this one.
	ConsumedMessage bool

	// HandledMessage means the message was not skipped. We actually
	// did something with the message which prompted this Response.
	HandledMessage bool

	// CancelledMessage means that the ConsumableMessage that was provided
	// caused the responses to be cancelled. Any cancellation in the response
	// listing should cancel all, and this will force-flag ConsumedMessage
	// and HandledMessage as true.
	CancelledMessage bool

	// Contents is Dungar's response back to the user.
	Contents string

	// ResponseType indicates how this response should be interpreted.
	ResponseType ResponseType

	// Origin indicates where this response is coming from
	Origin ResponseOrigin
}

// IsEmpty returns whether or not the contents is empty
func (rsp *Response) IsEmpty() bool {
	return strings.TrimSpace(rsp.Contents) == ""
}

// IsCancelled returns whether or not the response is cancelled.
func (rsp *Response) IsCancelled() bool {
	return rsp.CancelledMessage
}

// IsHandled returns whether or not the message is 'handled' i.e. the responses
// should be added to the envelope. This occurs when CancelledMessage and/or
// HandledMessage is true
func (rsp *Response) IsHandled() bool {
	return rsp.CancelledMessage || rsp.HandledMessage
}

// IsConsumed returns whether or not the message should be considered
// consumed. This occurs when ConsumedMessage and/or CancelledMessage is true
func (rsp *Response) IsConsumed() bool {
	return rsp.CancelledMessage || rsp.ConsumedMessage
}

// Prefixes changes PrefixUsername to true
func (rsp *Response) Prefixes() *Response {
	rsp.PrefixUsername = true
	return rsp
}

// Handles changes HandledMessage to true
func (rsp *Response) Handles() *Response {
	rsp.HandledMessage = true
	return rsp
}

// Consumes changes ConsumedMessage to true
func (rsp *Response) Consumes() *Response {
	rsp.ConsumedMessage = true
	return rsp
}

// Cancel will cancel the message in the post phase (or pre phase)
func (rsp *Response) Cancel() *Response {
	rsp.CancelledMessage = true
	return rsp
}

// IsValidOutgoing determines if the outgoing message is even valid
func (rsp *Response) IsValidOutgoing(driver ProtocolDriver, msg *IncomingMessage) bool {
	return strings.TrimSpace(rsp.Format(driver, msg)) != ""
}

// Format takes the PingUser of the frontend
func (rsp *Response) Format(driver ProtocolDriver, msg *IncomingMessage) string {
	if rsp.ResponseType == ResponseTypeReaction {
		return rsp.Contents
	}

	if rsp.PrefixUsername {
		user, err := driver.GetUser(msg.UserID, msg.ServerID)

		if err != nil {
			log.Printf("[Response.Format] Could not get user with ID '%s': %v", msg.UserID, err)

			return rsp.Contents
		}

		return driver.PingUser(user) + rsp.Contents
	}

	return rsp.Contents
}

// String converts to String
func (rsp *Response) String() string {
	return fmt.Sprintf(
		"handled='%v', prefix='%v', consume='%v', cancelled='%v', contents='%s', type='%s'",
		rsp.HandledMessage,
		rsp.PrefixUsername,
		rsp.ConsumedMessage,
		rsp.CancelledMessage,
		rsp.Contents,
		rsp.ResponseType,
	)
}
