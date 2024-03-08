package core2

// EmptyEventRsp is our shortcut of a non-consuming response to the event
// response pipeline
func EmptyEventRsp() []*EventResponse {
	return SingleEventRsp(&EventResponse{
		PrefixUsername: false,
		ConsumedEvent:  false,
		Contents:       "",
	})
}

// SingleEventRsp is a wrapper to make a single EventResponse into a slice
func SingleEventRsp(rsp *EventResponse) []*EventResponse {
	return []*EventResponse{rsp}
}

// MakeEventRsp is a shortcut for a simple event response with just a message
func MakeEventRsp(contents string) []*EventResponse {
	return SingleEventRsp(&EventResponse{
		PrefixUsername: false,
		ConsumedEvent:  true,
		Contents:       contents,
	})
}

// MakePrefixedEventRsp is a shortcut for a single event response with just a message
// that is directed to the user that initiated the event.
func MakePrefixedEventRsp(contents string) []*EventResponse {
	return SingleEventRsp(&EventResponse{
		PrefixUsername: true,
		ConsumedEvent:  true,
		Contents:       contents,
	})
}
