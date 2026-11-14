package core2

func isAnyCanceled(responses []*Response) bool {
	for _, rsp := range responses {
		// Well we did encounter a canceled message. Nothing goes out now.
		if rsp.IsCanceled() {
			return true
		}
	}

	return false
}

func isAnyHandled(responses []*Response) bool {
	// Can't handle message if there's no valid responses.
	if len(responses) == 0 {
		return false
	}

	for _, rsp := range responses {
		// Although these should both be true if CancelledMessage is
		// true, there might be a dumb case when they aren't
		if rsp.IsHandled() {
			return true
		}
	}

	return false
}

func didAnyConsume(responses []*Response) bool {
	for _, rsp := range responses {
		if rsp.IsConsumed() {
			return true
		}
	}

	return false
}
