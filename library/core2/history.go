package core2

// GetPreviousMessages will attempt to retrieve up to count previous backlog.
// If there are less previous backlog than count, the total amount of tracked
// previous backlog will be provided.
func (s *Service) GetPreviousMessages(count int) []*IncomingMessage {
	if count < 0 {
		return nil
	}

	if count > len(s.backlog) {
		count = len(s.backlog)
	}

	outbound := make([]*IncomingMessage, 0, count)

	for k := 0; k < count; k++ {
		outbound = append(outbound, s.backlog[k])
	}

	return outbound
}

// GetPreviousMessagesByChannel is similar to GetPreviousMessages but will explicitly only
// look for previous backlog in a given chaneID
func (s *Service) GetPreviousMessagesByChannel(chanID string, count int) []*IncomingMessage {
	if chanID == "" || count < 0 {
		return nil
	}

	if count > len(s.backlog) {
		count = len(s.backlog)
	}

	outbound := make([]*IncomingMessage, 0, count)

	for k, v := 0, 0; v < len(s.backlog) && k < count; k++ {
		for ; v < len(s.backlog); v++ {
			if s.backlog[v].ChannelID == chanID {
				outbound = append(outbound, s.backlog[v])
				break
			}
		}
	}

	return outbound
}
