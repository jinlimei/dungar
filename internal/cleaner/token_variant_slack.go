package cleaner

func (t *Token) parseSpecialForSlack() []Token {
	var (
		offs = t.StartPos
		out  = make([]Token, 0)
		wLen = len(t.Value)
		pos  = 0
		p1   = 0
		chr  rune
		nxt  rune

		tok1 TokenType

		peekMid int
		peekEnd int
	)

	if t.Value[0] == '`' && t.Value[wLen-1] == '`' {
		t.Type = TokenCode
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
			case '@':
				tok1 = TokenMentionUser
			case '#':
				tok1 = TokenMentionChannel
			default:
				tok1 = TokenURL
			}

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
		} // end of switch chr
	} // end of for ; pos <= wLen; pos++

	return out
}
