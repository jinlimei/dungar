package core2

import "strings"

func makeIncomingMessage(contents string, msgType MessageType) *IncomingMessage {
	return &IncomingMessage{
		IsSubMessage:  false,
		Contents:      contents,
		LowerContents: strings.ToLower(contents),
		Type:          msgType,
	}
}
