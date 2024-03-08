package cleaner

import "strings"

func (t *Token) parseWordForChat(variant TokenVariant) []Token {
	var (
		str   = string(t.Value)
		offs  = t.StartPos
		wLen  = len(t.Value)
		pos   = 0
		p1    = 0
		isNum = false
		out   = make([]Token, 0)
		chr   rune
		tok1  TokenType
		tok2  TokenType
	)

	if strings.HasPrefix(str, "https://") || strings.HasPrefix(str, "http://") {
		t.Type = TokenURL
		return []Token{*t}
	}

	for ; pos <= wLen; pos++ {
		chr = peek(t.Value, pos)
		//log.Printf("chr '%v' at pos '%d'\n", string(chr), pos)

		switch chr {
		case '@', '#':
			peekPos := wLen

			switch variant {
			case VariantDiscord, VariantSlack:
				switch chr {
				case '@':
					tok1 = TokenMentionUser
				case '#':
					tok1 = TokenMentionChannel
				default:
					tok1 = TokenURL
				}

				// Discord special-case: @&role
				if peek(t.Value, pos+1) == '&' {
					tok1 = TokenMentionRole
				}

				out = append(out, Token{
					Type:     tok1,
					StartPos: offs + pos,
					EndPos:   offs + peekPos + 1,
					Value:    getRuneRange(t.Value, pos, peekPos),
				})

				pos = peekPos
				p1 = peekPos + 1

			case VariantTwitter:
				switch chr {
				case '@':
					tok1 = TokenMentionUser
				case '#':
					tok1 = TokenHashTag
				case '$':
					// Although not technically a hashtag (it's a stock ticker or ... crypto thing?)
					// we will treat it as a hashtag because it's as cursed as one.
					tok1 = TokenHashTag
				}

				out = append(out, Token{
					Type:     tok1,
					StartPos: offs + pos,
					EndPos:   offs + peekPos + 1,
					Value:    getRuneRange(t.Value, pos, peekPos),
				})

				pos = peekPos
				p1 = peekPos + 1
			}

		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ',':
			if !isNum {
				if chr == ',' {
					break
				} else if intAbs(pos-p1) > 0 {
					//log.Printf("Making word with p1=%d and pos=%d (%v)\n", p1, pos, getRuneRange(t.Value, p1, pos))
					out = append(out, Token{
						Type:     TokenWord,
						StartPos: offs + p1,
						EndPos:   offs + pos,
						Value:    getRuneRange(t.Value, p1, pos),
					})
				}

				isNum = true
				p1 = pos
			}
		case '.':
			if isNum {
				if !isDigit(peek(t.Value, pos+1)) {
					isNum = false
					out = append(out, Token{
						Type:     TokenNumber,
						StartPos: offs + p1,
						EndPos:   offs + pos,
						Value:    getRuneRange(t.Value, p1, pos),
					})

					p1 = pos + 1
				}

				break
			} else if intAbs(pos-p1) > 0 {
				out = append(out, Token{
					Type:     TokenWord,
					StartPos: offs + p1,
					EndPos:   offs + pos,
					Value:    getRuneRange(t.Value, p1, pos),
				})
			}

			next := peek(t.Value, pos+1)
			if next == '.' {
				// We have a series of periods, let's see
				// if we can just consume all of them
				n := pos + 1
				for ; n < wLen; n++ {
					if peek(t.Value, n) != '.' {
						break
					}
				}

				out = append(out, Token{
					Type:     TokenEllipsis,
					StartPos: offs + pos,
					EndPos:   offs + n,
					Value:    getRuneRange(t.Value, pos, n),
				})

				pos = n
				p1 = n
			} else if peeped(next, '!', '?') {
				n := pos + 1
			nPosLoop2:
				for ; n < wLen; n++ {
					switch peek(t.Value, n) {
					case '!', '?', '.':
					default:
						break nPosLoop2
					}
				}

				out = append(out, Token{
					Type:     TokenSentEnd,
					StartPos: offs + pos,
					EndPos:   offs + n + 1,
					Value:    getRuneRange(t.Value, pos, n),
				})

				pos = n
				p1 = n + 1

			} else {
				out = append(out, Token{
					Type:     TokenPeriod,
					StartPos: offs + pos,
					EndPos:   offs + pos + 1,
					Value:    []rune{'.'},
				})

				pos = pos
				p1 = pos + 1
			}

		case nilRune:
			if intAbs(pos-p1) == 0 {
				break
			}

			tok1 = TokenWord

			if isNum {
				isNum = false
				tok1 = TokenNumber
			}

			//log.Printf("Making word with p1=%d and pos=%d (%v from %v)\n", p1, pos, getRuneRange(t.Value, p1, pos), t.Value)
			out = append(out, Token{
				Type:     tok1,
				StartPos: offs + p1,
				EndPos:   offs + pos,
				Value:    getRuneRange(t.Value, p1, pos),
			})

		case '!', '?':
			tok1 = TokenWord

			if isNum {
				isNum = false
				tok1 = TokenNumber
			}

			if intAbs(pos-p1) > 0 {
				out = append(out, Token{
					Type:     tok1,
					StartPos: offs + p1,
					EndPos:   offs + pos,
					Value:    getRuneRange(t.Value, p1, pos),
				})
			}

			nPos := pos + 1
		nPosLoop:
			for ; nPos < wLen; nPos++ {
				switch peek(t.Value, nPos) {
				case '!', '?', '.':
				default:
					break nPosLoop
				}
			}

			out = append(out, Token{
				Type:     TokenSentEnd,
				StartPos: offs + pos,
				EndPos:   offs + nPos + 1,
				Value:    getRuneRange(t.Value, pos, nPos),
			})

			pos = nPos
			p1 = pos + 1

		case '(', ')', '{', '}', '[', ']':
			tok1 = TokenWord

			if isNum {
				isNum = false
				tok1 = TokenNumber
			}

			if intAbs(pos-p1) > 0 {
				out = append(out, Token{
					Type:     tok1,
					StartPos: offs + p1,
					EndPos:   offs + pos,
					Value:    getRuneRange(t.Value, p1, pos),
				})
			}

			switch chr {
			case '(':
				tok2 = TokenOpenParens
			case ')':
				tok2 = TokenCloseParens
			case '{':
				tok2 = TokenOpenBracket
			case '}':
				tok2 = TokenCloseBracket
			case '[':
				tok2 = TokenOpenBrace
			case ']':
				tok2 = TokenCloseBrace
			}

			out = append(out, Token{
				Type:     tok2,
				StartPos: offs + pos,
				EndPos:   offs + pos + 1,
				Value:    []rune{chr},
			})

		//case ':':
		//	tok1 = TokenWord
		//
		//	if isNum {
		//		isNum = false
		//		tok1 = TokenNumber
		//
		//		if intAbs(pos - p1) > 0 {
		//			out = append(out, Token{
		//				Type:     tok1,
		//				StartPos: p1,
		//				EndPos:   pos,
		//				Value:    getRuneRange(t.Value, p1, pos),
		//			})
		//		}
		//	}
		//
		//	nextColon := peekUntil(t.Value, pos + 1, ':')
		//	if nextColon < 0 {
		//		break
		//	}
		//
		//	if intAbs(pos - p1) > 0 {
		//		out = append(out, Token{
		//			Type:     tok1,
		//			StartPos: p1,
		//			EndPos:   pos,
		//			Value:    getRuneRange(t.Value, p1, pos),
		//		})
		//	}
		//
		//	out = append(out, Token{
		//		Type:     TokenEmoticon,
		//		StartPos: pos,
		//		EndPos:   nextColon,
		//		Value:    getRuneRange(t.Value, pos, nextColon),
		//	})
		//
		//	pos = nextColon + 1

		case '"', '\'':
			tok2 = TokenSingleQuote

			if chr == '"' {
				tok2 = TokenDoubleQuote
			}

			tok1 = TokenWord

			if isNum {
				isNum = false
				tok1 = TokenNumber

				if intAbs(pos-p1) > 0 {
					out = append(out, Token{
						Type:     tok1,
						StartPos: offs + p1,
						EndPos:   offs + pos,
						Value:    getRuneRange(t.Value, p1, pos),
					})
				}
			}

			if peeped(peek(t.Value, pos-1), ' ', nilRune) {
				out = append(out, Token{
					Type:     tok2,
					StartPos: offs + pos,
					EndPos:   offs + pos + 1,
					Value:    []rune{chr},
				})

				p1 = pos + 1
			} else if peeped(peek(t.Value, pos+1), ' ', nilRune) {
				if intAbs(pos-p1) > 0 {
					out = append(out, Token{
						Type:     tok1,
						StartPos: offs + p1,
						EndPos:   offs + pos,
						Value:    getRuneRange(t.Value, p1, pos),
					})
				}

				out = append(out, Token{
					Type:     tok2,
					StartPos: offs + pos,
					EndPos:   offs + pos + 1,
					Value:    []rune{chr},
				})

				p1 = pos + 1
			}

			//tok1 = TokenWord
			//tok2 = TokenSingleQuote
			//if chr == '"' {
			//	tok2 = TokenDoubleQuote
			//}
			//
			//if isNum {
			//	isNum = false
			//	tok1 = TokenNumber
			//}
			//
			//if intAbs(pos-p1) > 0 {
			//	out = append(out, Token{
			//		Type:     tok1,
			//		StartPos: offs + p1,
			//		EndPos:   offs + pos,
			//		Value:    getRuneRange(t.Value, p1, pos),
			//	})
			//}
			//
			//out = append(out, Token{
			//	Type:     tok2,
			//	StartPos: offs + pos,
			//	EndPos:   offs + pos + 1,
			//	Value:    []rune{chr},
			//})
			//
			//p1 = pos + 1

		default:
			if isNum {
				isNum = false
				out = append(out, Token{
					Type:     TokenNumber,
					StartPos: offs + p1,
					EndPos:   offs + pos,
					Value:    getRuneRange(t.Value, p1, pos),
				})

				p1 = pos
			}
		}
	}

	//spew.Dump(out[0].Type, out[0].Value, out[0].StartPos, out[0].EndPos)
	return out
}
