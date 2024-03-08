package core2

// ScheduledMessages is a typed slice/array of ScheduledMessage
type ScheduledMessages []*ScheduledMessage

// ToResponses will take all members and convert them to a Response via
// ScheduledMessage.ToResponse
func (sm ScheduledMessages) ToResponses() []*Response {
	out := make([]*Response, 0, len(sm))

	for _, msg := range sm {
		out = append(out, msg.ToResponse())
	}

	return out
}

// ToResponseEnvelopes will take all members and convert them to a slice of
// ResponseEnvelope's grouped by ScheduledMessage.ChannelID's
func (sm ScheduledMessages) ToResponseEnvelopes() []*ResponseEnvelope {
	msgToChan := make(map[string]ScheduledMessages)

	for _, msg := range sm {
		if _, ok := msgToChan[msg.ChannelID]; !ok {
			msgToChan[msg.ChannelID] = make(ScheduledMessages, 0)
		}

		msgToChan[msg.ChannelID] = append(msgToChan[msg.ChannelID], msg)
	}

	envelopes := make([]*ResponseEnvelope, 0, len(msgToChan))

	for chanID, msgs := range msgToChan {

		envelopes = append(envelopes, &ResponseEnvelope{
			Message: &IncomingMessage{
				ChannelID: chanID,
				Type:      MessageTypeBasic,
			},
			Responses: msgs.ToResponses(),
		})
	}

	return envelopes
}
