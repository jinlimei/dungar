package cleaner

// SimpleToken is a simplified struct of Token which is useful for ingesting
// into Markov3 and maybe other things.
type SimpleToken struct {
	Type     TokenType
	Value    string
	Override string
}

// SimpleTokenList is a listing for SimpleToken
type SimpleTokenList struct {
	// Groupings provides a group of tokens by TokenType
	Groupings map[TokenType][]SimpleToken
	// Tokens is our slice of SimpleToken
	Tokens []SimpleToken
}

// GetSimpleTokenList will return a SimpleTokenList & parse out the stuff
// based on a predefined set of rules.
func (tl *TokenList) GetSimpleTokenList() SimpleTokenList {
	if tl.simple == nil {
		stl := &SimpleTokenList{
			Groupings: make(map[TokenType][]SimpleToken, 25),
			Tokens:    make([]SimpleToken, 0, len(tl.tokens)),
		}

		var (
			tLen = len(tl.tokens)
			pos  = 0
			//cOk  bool
			//nOk  bool
			cur Token
			//nxt  Token
			tmp Token
			wrk SimpleToken

			p1 int
			//p2 int
		)

	tokenLoop:
		for ; pos < tLen; pos++ {
			cur, _ = peekToken(tl.tokens, pos)
			//nxt, nOk = peekToken(tl.tokens, pos+1)
			wrk = SimpleToken{
				Type:     TokenEOL,
				Value:    "",
				Override: "",
			}

			switch cur.Type {
			case TokenDoubleQuote, TokenSingleQuote, TokenOpenParens, TokenCloseParens,
				TokenOpenBrace, TokenCloseBrace, TokenOpenBracket, TokenCloseBracket, TokenMalformed, TokenCode:
				// do nothing with these
			case TokenEOL:
				break tokenLoop
			case TokenSpace:
				// do nothing with this too!
			case TokenSentEnd, TokenPeriod, TokenEllipsis:
				wrk.Type = TokenSentEnd
				wrk.Value = string(cur.Value)
				wrk.Override = ""

				stl.Tokens = append(stl.Tokens, wrk)

			case TokenURL, TokenEmoticon, TokenMentionUser, TokenMentionRole,
				TokenMentionChannel, TokenHashTag:
				wrk.Type = cur.Type
				wrk.Value = string(cur.Value)
				wrk.Override = string(cur.Override)

				stl.Tokens = append(stl.Tokens, wrk)

			case TokenWord, TokenNumber, TokenCommand:
				wrk.Type = TokenWord
				wrk.Value = string(cur.Value)
				wrk.Override = string(cur.Override)

				for p1 = pos + 1; p1 < tLen; p1++ {
					tmp, _ = peekToken(tl.tokens, p1)
					if !isWordConsumableToken(tmp.Type) {
						break
					}

					wrk.Value += string(tmp.Value)
					wrk.Override += string(tmp.Override)
				} // end of for p1=pos+1

				stl.Tokens = append(stl.Tokens, wrk)
				pos = p1
			} // end of switch cur.Type
		} // end of tokenLoop

		newList := make([]SimpleToken, 0)
		for _, tok := range stl.Tokens {
			w := cleanWord(tok.Value)
			if w == "" {
				continue
			}

			tok.Value = w
			newList = append(newList, tok)
		}

		stl.Tokens = newList

		for _, tok := range stl.Tokens {
			_, ok := stl.Groupings[tok.Type]
			if !ok {
				stl.Groupings[tok.Type] = make([]SimpleToken, 0)
			}

			stl.Groupings[tok.Type] = append(stl.Groupings[tok.Type], tok)
		} // end of stl.Tokens loop for grouping

		tl.simple = stl
	} // end of tl.simple == nil

	return *tl.simple
}

func isWordConsumableToken(tt TokenType) bool {
	var valid = []TokenType{
		TokenWord,
		TokenNumber,
		TokenPeriod,
		TokenSentEnd,
		TokenEllipsis,
		TokenCommand,
	}

	for _, vt := range valid {
		if tt == vt {
			return true
		}
	}

	return false
}
