package cleaner

func peekToken(tl []Token, pos int) (Token, bool) {
	if pos >= len(tl) {
		return Token{}, false
	}

	if pos < 0 {
		return Token{}, false
	}

	return tl[pos], true
}

func peekUntilTokens(tl []Token, pos int, tts... TokenType) (int, TokenType) {
	if pos >= len(tl) {
		return -1, TokenEOL
	}

	if pos < 0 {
		return -1, TokenEOL
	}

	for ; pos < len(tl); pos++ {
		for _, tt := range tts {
			if tl[pos].Type == tt {
				return pos, tt
			}
		}
	}

	return -1, TokenEOL
}

func peekUntilToken(tl []Token, pos int, t TokenType) int {
	if pos >= len(tl) {
		return -1
	}

	if pos < 0 {
		return -1
	}

	for ; pos < len(tl); pos++ {
		if tl[pos].Type == t {
			return pos
		}
	}

	return -1
}

func getTokenGroup(tl []Token, start, end int) ([]Token, bool) {
	if start < 0 || end < 0 {
		return nil, false
	}

	tLen := len(tl)

	if start >= tLen || end >= tLen {
		return nil, false
	}

	return tl[start:end], true
}
