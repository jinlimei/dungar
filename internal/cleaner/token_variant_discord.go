package cleaner

func (t *Token) parseSpecialForDiscord() []Token {
	// emoji: \u003c:blobyes:900369292292403211\u003e
	// @person: \u003c@569251658773692538\u003e test
	// @person: \u003c@569251658773692538\u003e do you have content now

	var (
		offs = t.StartPos
		out  = make([]Token, 0)
		wLen = len(t.Value)
		pos  = 0
		p1   = 0
		chr  rune
		nxt  rune

		tok1    TokenType
		peekMid int
		peekEnd int
	)

	if wLen == 2 {
		t.Type = TokenWord
		return []Token{*t}
	}

	for ; pos <= wLen; pos++ {
		chr = peek(t.Value, pos)
		nxt = peek(t.Value, pos+1)

		switch chr {
		case '<':
			if intAbs(pos-p1) > 0 {
				out = append(out, Token{
					Type:     TokenWord,
					StartPos: offs + p1,
					EndPos:   offs + pos + 1,
					Value:    getRuneRange(t.Value, p1, pos),
				})

				p1 = pos + 1
			}

			switch nxt {
			case ':':
				tok1 = TokenEmoticon
			case '@':
				tok1 = TokenMentionUser

				if peek(t.Value, pos+2) == '&' {
					tok1 = TokenMentionRole
				}

			case '#':
				tok1 = TokenMentionChannel
			default:
				tok1 = TokenURL
			}

			switch tok1 {
			case TokenMentionRole:

				// TODO test
				out = append(out, Token{
					Type:     tok1,
					StartPos: offs + pos,
					EndPos:   offs + wLen,
					Value:    getRuneRange(t.Value, pos, wLen),
				})

			case TokenEmoticon:
				peekMid = peekUntil(t.Value, pos+2, ':')
				peekEnd, _ = peekUntilSet(t.Value, pos+2, ' ', '>', '\u0000')

				if peekEnd < 0 {
					peekEnd = wLen
				}

				if peekMid < 0 {
					out = append(out, Token{
						Type:     TokenMalformed,
						StartPos: offs + pos,
						EndPos:   offs + wLen,
						Value:    getRuneRange(t.Value, pos, wLen),
					})

					pos = wLen - 1
					p1 = wLen

					break
				}

				// We have an end : and it's before the end
				// so it's probably ':blobyes:1234567890
				if peekMid > 0 && peekMid < peekEnd {
					out = append(out, Token{
						Type:     tok1,
						StartPos: offs + pos,
						EndPos:   offs + peekEnd,
						Value:    getRuneRange(t.Value, pos+1, peekMid+1),
						Override: getRuneRange(t.Value, pos+1, peekEnd),
					})

					pos = peekEnd
					p1 = peekEnd + 1
				} else {
					out = append(out, Token{
						Type:     tok1,
						StartPos: offs + pos,
						EndPos:   offs + peekEnd,
						Value:    getRuneRange(t.Value, pos+1, peekEnd),
					})

					pos = peekEnd
					p1 = peekEnd + 1
				}

			case TokenMentionUser, TokenMentionChannel, TokenURL:
				peekMid = peekUntil(t.Value, pos+1, '|')
				peekEnd = peekUntil(t.Value, pos+1, '>')

				if peekEnd < 0 {
					out = append(out, Token{
						Type:     TokenMalformed,
						StartPos: offs + pos,
						EndPos:   offs + wLen,
						Value:    getRuneRange(t.Value, pos, wLen),
					})

					pos = wLen - 1
					p1 = wLen

					break
				}

				// We have a midpoint and a valid end!
				if peekEnd > 0 && peekMid > 0 && peekMid < peekEnd {
					out = append(out, Token{
						Type:     tok1,
						StartPos: offs + pos,
						EndPos:   offs + peekEnd,
						Value:    getRuneRange(t.Value, pos+1, peekMid),
						Override: getRuneRange(t.Value, peekMid+1, peekEnd),
					})

					pos = peekEnd
					p1 = peekEnd + 1
				} else if peekEnd > 0 && (peekMid < 0 || peekMid > peekEnd) {
					out = append(out, Token{
						Type:     tok1,
						StartPos: offs + pos,
						EndPos:   offs + peekEnd,
						Value:    getRuneRange(t.Value, pos+1, peekEnd),
					})

					pos = peekEnd
					p1 = peekEnd + 1
				}
			}
		}
	}

	return out
}
