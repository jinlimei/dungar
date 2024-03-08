package cleaner

import (
	"strings"
)

const nilRune = rune(0x00)

var (
	unwantedReplacer = strings.NewReplacer(
		string([]rune{nilRune}), "",
		string([]rune{'\r'}), "",
		string([]rune{8203}), "",
		"‘", "'",
		"’", "'",
		"“", "\"",
		"”", "\"",
		"\t", " ",
	)
)

// TokenizePrep prepares a string for tokenizing with some common actions
// such as: trimming space, removing \r, etc.
func TokenizePrep(str string) string {
	str = strings.TrimSpace(str)
	str = unwantedReplacer.Replace(str)

	for strings.Contains(str, "  ") {
		str = strings.ReplaceAll(str, "  ", " ")
	}

	return str
}

// Tokenize takes the incoming string and variant and builds a TokenList
// based off of that handling.
func Tokenize(str string, variant TokenVariant) TokenList {
	str = TokenizePrep(str)

	var (
		tokens = make([]Token, 0)
		runes  = []rune(str)
		rLen   = len(runes)
	)

	// Nothing to do here I'm afraid
	if rLen == 0 {
		return TokenList{
			tokens: []Token{},
		}
	}

	var (
		pos       = 0
		tokStart  = 0
		lastSpace = 0

		chr     rune
		next    rune
		inBrace bool
	)

	for ; pos <= rLen; pos++ {
		chr = peek(runes, pos)
		next = peek(runes, pos+1)

		switch chr {
		case nilRune:
			if inBrace {
				tokens = append(tokens, Token{
					Type:     TokenMalformed,
					StartPos: tokStart,
					EndPos:   pos,
					Value:    getRuneRange(runes, tokStart, pos),
				})

				break
			}

			if intAbs(pos-lastSpace) > 0 {
				tokens = append(tokens, Token{
					Type:     TokenWord,
					StartPos: lastSpace,
					EndPos:   pos,
					Value:    getRuneRange(runes, lastSpace, pos),
				})
			}

			//spew.Dump(getRuneRange(runes, lastSpace, pos))

			tokens = append(tokens, Token{
				Type:     TokenEOL,
				StartPos: pos,
				EndPos:   pos + 1,
				Value:    []rune{nilRune},
			})

		case ' ', '\n':
			if inBrace {
				break
			}

			tok := TokenSpace
			if chr == '\n' {
				tok = TokenNewLine
			}

			if intAbs(pos-lastSpace) > 0 {
				tokens = append(tokens, Token{
					Type:     TokenWord,
					StartPos: lastSpace,
					EndPos:   pos,
					Value:    getRuneRange(runes, lastSpace, pos),
				})
			}

			tokens = append(tokens, Token{
				Type:     tok,
				StartPos: pos,
				EndPos:   pos + 1,
				Value:    []rune{chr},
			})

			lastSpace = pos + 1

		case '(':
			// We only do this for the first segment
			if variant != VariantXMPP || pos != 0 {
				break
			}

			prEnd := inlinePeekUntil(runes, pos+1, ')')

			if intAbs(pos - lastSpace) > 0 {
				tokens = append(tokens, Token{
					Type:     TokenWord,
					StartPos: lastSpace,
					EndPos:   pos,
					Value:    getRuneRange(runes, lastSpace, pos),
				})

				lastSpace = pos + 1
			}

			if prEnd > 0 {
				rng := getRuneRange(runes, pos, prEnd+1)

				if isMaybeTimestamp(rng) {
					tokens = append(tokens, Token{
						Type:     TokenSpecial,
						StartPos: pos,
						EndPos:   prEnd + 1,
						Value:    rng,
					})

					pos = prEnd
					lastSpace = prEnd + 1
				}
			}

		case '<':
			if variant.isSimpleVariant() {
				break
			}

			braceEnd := inlinePeekUntil(runes, pos+1, '>')
			if braceEnd < 0 {
				break
			}

			if intAbs(pos-lastSpace) > 0 {
				tokens = append(tokens, Token{
					Type:     TokenWord,
					StartPos: lastSpace,
					EndPos:   pos,
					Value:    getRuneRange(runes, lastSpace, pos),
				})

				lastSpace = pos + 1
			}

			tokens = append(tokens, Token{
				Type:     TokenSpecial,
				StartPos: pos,
				EndPos:   braceEnd + 1,
				Value:    getRuneRange(runes, pos, braceEnd+1),
			})

			pos = braceEnd
			lastSpace = braceEnd + 1

		case '`':
			if !variant.isMarkdownVariant() {
				break
			}

			if intAbs(pos-lastSpace) > 0 {
				tokens = append(tokens, Token{
					Type:     TokenWord,
					StartPos: lastSpace,
					EndPos:   pos,
					Value:    getRuneRange(runes, lastSpace, pos),
				})

				lastSpace = pos + 1
			}

			isThree := peek(runes, pos+1) == '`' && peek(runes, pos+2) == '`'
			// We want to do something a bit more fancy when we're looking for
			// the triple backticks
			if isThree && rLen > pos+3 {
				var (
					bPos  = pos + 3
					found = false
					c1    rune
					c2    rune
					c3    rune
				)

				for ; bPos < rLen; bPos++ {
					c1 = peek(runes, bPos)
					c2 = peek(runes, bPos+1)
					c3 = peek(runes, bPos+2)

					if c1 == '`' && c2 == '`' && c3 == '`' {
						found = true
						break
					}
				}

				if !found {
					tokens = append(tokens, Token{
						Type:     TokenMalformed,
						StartPos: pos,
						EndPos:   rLen,
						Value:    getRuneRange(runes, pos, rLen),
					})

					pos = rLen - 1
					lastSpace = rLen
					break
				}

				tokens = append(tokens, Token{
					Type:     TokenSpecial,
					StartPos: pos,
					EndPos:   bPos + 3,
					Value:    getRuneRange(runes, pos, bPos+3),
				})

				// but why?
				// next round pos is incremented
				// so this becomes defacto bPos+3 so we should get the character
				// after we want. weird?
				pos = bPos + 2
				lastSpace = pos + 1
				//log.Printf("pos: %d, peek: '%v'\n",
				//	pos, peek(runes, pos))
			} else {
				next := peekUntil(runes, pos+1, '`')
				btCnt := 0
				if next > 0 {
					for p1 := next; p1 < rLen; p1++ {
						if peek(runes, p1) == '`' {
							btCnt++
						} else {
							break
						}
					}
					//nextThree := peek(runes, next+1) == '`' && peek(runes, next+2) == '`' && peek(runes, next+3) == '`'
					//log.Printf("Count of next backticks: %d\n", btCnt)
					//log.Printf("Looking at (potential): %s\n", string(getRuneRange(runes, pos, next+1)))
					//log.Printf("Looking at (next,next+btCnt): %s\n", string(getRuneRange(runes, next, next+btCnt)))
				}

				if next > pos && (btCnt == 1 || btCnt%2 == 0 || peekUntil(runes, next+btCnt+1, '`') < 0) {
					tokens = append(tokens, Token{
						Type:     TokenSpecial,
						StartPos: pos,
						EndPos:   next + 1,
						Value:    getRuneRange(runes, pos, next+1),
					})

					pos = next
					lastSpace = pos + 1
				}
			}

		case '>':
			//if inBrace {
			//	tokens = append(tokens, Token{
			//		Type:  TokenSpecial,
			//		Value: getRuneRange(runes, tokStart+1, pos),
			//	})
			//
			//	tokStart = pos + 1
			//	lastSpace = pos + 1
			//}

		case ':':
			if !variant.isEmojiVariant() {
				break
			}

			if inBrace {
				break
			}

			if next == ' ' {
				break
			}

			if next == ':' {
				break
			}

			//log.Printf("Looking at colon at pos %d\n", pos)
			// we're trying to match emoticons here, so any 'weird' characters should
			// exclude us from matching said emoticons
			peekPos, peekChr := peekUntilSet(runes, pos+2, ':', ' ', '?', '/', '\n', '`')

			if peekPos < 0 {
				break
			}

			if peekChr != ':' {
				break
			}

			if intAbs(pos-lastSpace) > 0 {
				tokens = append(tokens, Token{
					Type:     TokenWord,
					StartPos: lastSpace,
					EndPos:   pos,
					Value:    getRuneRange(runes, lastSpace, pos),
				})
			}

			tokens = append(tokens, Token{
				Type:     TokenEmoticon,
				StartPos: pos,
				EndPos:   peekPos + 1,
				Value:    getRuneRange(runes, pos, peekPos+1),
			})

			lastSpace = peekPos + 1
			tokStart = peekPos
			pos = peekPos
		} // end of switch
	} // end of for

	output := make([]Token, 0, len(tokens))

	//DebugPrintTokenList(tokens)

	for _, token := range tokens {
		//log.Println(token.DebugString())
		updated := token.Parse(variant)
		output = append(output, updated...)
	}

	tokens = output
	output = make([]Token, 0, len(tokens))
	for pos, token := range tokens {
		if token.Type == TokenSpace && pos == 0 {
			continue
		}

		output = append(output, token)
	}

	return TokenList{
		tokens: output,
	}
}
