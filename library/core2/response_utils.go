package core2

// EmptyRsp provides an empty response
func EmptyRsp() []*Response {
	return SingleRsp(&Response{
		PrefixUsername:   false,
		ConsumedMessage:  false,
		HandledMessage:   false,
		CancelledMessage: false,
		Contents:         "",
		ResponseType:     ResponseTypeBasic,
		Origin:           OriginEvent,
	})
}

// MakeRsp will make a valid response with the incoming msg
func MakeRsp(msg string) *Response {
	return &Response{
		PrefixUsername:   false,
		ConsumedMessage:  true,
		HandledMessage:   true,
		CancelledMessage: false,
		Contents:         msg,
		ResponseType:     ResponseTypeBasic,
		Origin:           OriginEvent,
	}
}

// ReactionRsp is for emoji-reaction responses to a given message
func ReactionRsp(emoji string) *Response {
	r := MakeRsp(emoji)
	r.ResponseType = ResponseTypeReaction

	return r
}

// SingleReactionRsp is a shortcut for SingleRsp and ReactionRsp
func SingleReactionRsp(emoji string) []*Response {
	return SingleRsp(ReactionRsp(emoji))
}

// DirectMessageRsp is for when the bot will respond to the user via a DM even if the
// original comment was not in a DM
func DirectMessageRsp(msg string) *Response {
	r := MakeRsp(msg)
	r.ResponseType = ResponseTypeDirectMessage

	return r
}

// SingleDirectMessageRsp is a shortcut for SingleRsp and DirectMessageRsp
func SingleDirectMessageRsp(msg string) []*Response {
	return SingleRsp(DirectMessageRsp(msg))
}

// PrefixedSingleRsp is SingleRsp and MakeRsp but with Response.PrefixUsername set to true
func PrefixedSingleRsp(msg string) []*Response {
	r := MakeRsp(msg)
	r.PrefixUsername = true

	return SingleRsp(r)
}

// CancelledRsp is for a cancelled/rejected response
func CancelledRsp() []*Response {
	return SingleRsp(&Response{
		PrefixUsername:   false,
		ConsumedMessage:  true,
		HandledMessage:   true,
		CancelledMessage: true,
		Contents:         "",
		ResponseType:     ResponseTypeBasic,
		Origin:           OriginEvent,
	})
}

// MakeSingleRsp is a shortcut for SingleRsp and MakeRsp
func MakeSingleRsp(msg string) []*Response {
	return SingleRsp(MakeRsp(msg))
}

// SingleRsp will wrap a single Response and convert it into a slice
func SingleRsp(rsp *Response) []*Response {
	return []*Response{rsp}
}
