package markov4

// getTokenID lets us use a -1 and also +(over-max) and not
// have things bork due to out of range errors.
func getTokenID(ids []TokenID, maxLen, pos int) TokenID {
	if pos < 0 {
		return LeftMost
	}

	if (pos + 1) > maxLen {
		return RightMost
	}

	return ids[pos]
}
