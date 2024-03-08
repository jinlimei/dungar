package markov3

import (
	"fmt"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// ChatTokenType is our underlying type for tokens
type ChatTokenType uint8

const (
	// ChatTokenWord represents an arbitrary string
	ChatTokenWord ChatTokenType = iota
	// ChatTokenSpace is well.. space
	ChatTokenSpace
	// ChatTokenNewLine identifies newlines in a msg
	ChatTokenNewLine
	// ChatTokenSpecial is an intermediate potentially converting
	// to other tokens (like ChatTokenURL or ChatTokenEmoticon, etc.)
	ChatTokenSpecial
	// ChatTokenMalformed is for tokens that are borked
	ChatTokenMalformed
	// ChatTokenURL represents URLs
	ChatTokenURL
	// ChatTokenEmoticon is for :condi:
	ChatTokenEmoticon
	// ChatTokenCode handles `code` and ```code```
	ChatTokenCode
	// ChatTokenQuote is not working right now
	// TODO fix this
	ChatTokenQuote
)

// ChatSubTokenType is a sub-type primarily for ChatTokenSpecial but also
// provides values for ChatSubTokenURL and ChatSubTokenEmoticon
type ChatSubTokenType string

const (
	// ChatSubTokenURL is a lower-level handler for ChatTokenURL
	ChatSubTokenURL ChatSubTokenType = "\u0000URL\u0000"

	// ChatSubTokenEmoticon is our lower-level handler for ChatTokenEmoticon
	ChatSubTokenEmoticon ChatSubTokenType = "\u0000EMOTICON\u0000"

	// ChatSubTokenUser is for user-mention
	ChatSubTokenUser ChatSubTokenType = "\u0000USER\u0000"

	// ChatSubTokenChannel is for channel-mention
	ChatSubTokenChannel ChatSubTokenType = "\u0000CHANNEL\u0000"

	// ChatSubTokenRole is for role-mention (although not explicitly supported
	// on all chat programs)
	ChatSubTokenRole ChatSubTokenType = "\u0000ROLE\u0000"

	// ChatSubTokenUnknown lets us map to an unknown token when something
	// went wrong with the matching.
	ChatSubTokenUnknown ChatSubTokenType = "\u0000UNKNOWN\u0000"
)

// TokenTypeToString takes the incoming t
// and returns a string representation of it
func TokenTypeToString(t ChatTokenType) string {
	switch t {
	case ChatTokenWord:
		return "ChatTokenWord"
	case ChatTokenNewLine:
		return "ChatTokenNewLine"
	case ChatTokenSpace:
		return "ChatTokenSpace"
	case ChatTokenSpecial:
		return "ChatTokenSpecial"
	case ChatTokenMalformed:
		return "ChatTokenMalformed"
	case ChatTokenURL:
		return "ChatTokenURL"
	case ChatTokenEmoticon:
		return "ChatTokenEmoticon"
	case ChatTokenCode:
		return "ChatTokenCode"
	case ChatTokenQuote:
		return "ChatTokenQuote"
	}
	return ""
}

// ChatToken is our struct for holding the type and value
type ChatToken struct {
	// TokenType
	TokenType ChatTokenType
	// Value is our value associated with this type
	Value []rune
}

// String converts the token into a debug string
func (ct *ChatToken) String() string {
	return fmt.Sprintf("%v: '%s'", TokenTypeToString(ct.TokenType), string(ct.Value))
}

// Analyze will look at the tokens Value and determine
// if its type should change as a result
func (ct *ChatToken) Analyze() []ChatToken {
	out := make([]ChatToken, 1)
	out[0] = *ct

	if ct.TokenType == ChatTokenEmoticon {
		colonCount := 0
		for _, c := range ct.Value {
			if c == ':' {
				colonCount++
			}
		}

		// oh boy! this means our thing is actually
		// :emoticon::emoticon:
		if colonCount > 2 {
			out = make([]ChatToken, 0)
			vLen := len(ct.Value)
			peek := func(idx int) rune {
				if idx < 0 {
					return rune(0x00)
				}

				if idx >= vLen {
					return rune(0x00)
				}

				return ct.Value[idx]
			}

			peekUntil := func(start int, until rune) int {
				for nxt := start; nxt <= vLen; nxt++ {
					if nxt < vLen && ct.Value[nxt] == until {
						return nxt
					}
				}

				return -1
			}

			getRange := func(start, end int) []rune {
				if end < 0 || start < 0 || end <= start || start > vLen || end > vLen {
					panic(fmt.Sprintf("OHGOD: '%v' (start=%d,end=%d,maxLen=%d)", string(ct.Value), start, end, vLen))
				}

				return ct.Value[start:end]
			}

			var (
				spt1 int
				spt2 int
			)

			var chr rune
			for pos := 0; pos <= vLen; pos++ {
				chr = peek(pos)

				switch chr {
				case rune(0x00):
					if spt2 == pos {
						break
					}

					out = append(out, ChatToken{
						TokenType: ChatTokenWord,
						Value:     getRange(spt2, pos),
					})

				case ':':
					spt1 = peekUntil(pos+1, ':')
					if spt1 < 0 {
						break
					}

					if spt1 == pos+1 {
						break
					}

					out = append(out, ChatToken{
						TokenType: ChatTokenEmoticon,
						Value:     getRange(pos, spt1+1),
					})

					pos = spt1 + 1
					spt2 = spt1 + 1
				}
			}
		}
	} else if ct.TokenType == ChatTokenSpecial {
		// inside of the <>'s

		//log.Printf("ChatTokenSpecial: %v\n", string(ct.Value))
		if len(ct.Value) == 0 {
			log.Printf("Encountered zero-length token, returning nothin' (%+v)\n", ct)
			return nil
		} else if len(ct.Value) == 1 {
			ct.TokenType = ChatTokenWord
		} else {
			inner := string(ct.Value)

			if utils.IsURL(inner) {
				ct.TokenType = ChatTokenURL
				ct.Value = []rune(inner)
			}
		}

		out[0] = *ct
	} else if ct.TokenType == ChatTokenWord && strings.HasPrefix(string(ct.Value), "http") {
		if utils.IsURL(string(ct.Value)) {
			ct.TokenType = ChatTokenURL
			out[0] = *ct
		}
	}

	return out
}

// TokenList is our valid list of ChatToken's
// parsed from a string of some kind
type TokenList struct {
	rawMsg string
	// Tokens is our ChatToken output
	Tokens []ChatToken

	Emoticons []string
	URLs      []string
	Users     []string
	Channels  []string
}

// RawString provides the raw untouched
// message that we are operating over.
func (tl *TokenList) RawString() string {
	return tl.rawMsg
}

// String outputs all the chat tokens values
func (tl *TokenList) String() string {
	output := ""
	for _, token := range tl.Tokens {
		output += string(token.Value)
	}

	return output
}

// MarkovConsumable provides a list of tokens that make
// sense to be consumed by a markov generator
func (tl *TokenList) MarkovConsumable() []string {
	output := make([]string, 0)

	for _, token := range tl.Tokens {
		switch token.TokenType {
		case ChatTokenSpace:
			// do nothing
		case ChatTokenWord:
			output = append(output, cleanWord(string(token.Value)))
		case ChatTokenURL:
			output = append(output, string(ChatSubTokenURL))
		case ChatTokenEmoticon:
			output = append(output, string(ChatSubTokenEmoticon))
		case ChatTokenCode:
			// do nothing
		case ChatTokenQuote:
			// dunno bud
		case ChatTokenMalformed:
			output = append(output, string(token.Value))
		case ChatTokenSpecial:
			if len(token.Value) > 0 {
				id := token.Value[0]

				if id == '@' && token.Value[1] != '&' {
					if !utils.StringInSlice(string(token.Value), tl.Users) {
						tl.Users = append(tl.Users, string(token.Value))
					}

					output = append(output, string(ChatSubTokenUser))
				} else if id == '#' {
					if !utils.StringInSlice(string(token.Value), tl.Channels) {
						tl.Channels = append(tl.Channels, string(token.Value))
					}

					output = append(output, string(ChatSubTokenChannel))
				} else if id == '&' || (id == '@' && token.Value[1] == '&') {
					output = append(output, string(ChatSubTokenRole))
				} else {
					db.LogIssue("unknown special token",
						fmt.Sprintf("unknown special token '%s'", string(token.Value)),
						fmt.Sprintf("Msg: '%s'", tl.RawString()))

					log.Printf("Unknown Token '%v' declared!\n", string(token.Value))

					output = append(output, string(ChatSubTokenUnknown))
				}
			} else {
				db.LogIssue("unknown special token",
					"unknown special token (zero-length!)",
					fmt.Sprintf("Msg: %v\n", tl.RawString()))

				log.Printf("Unknown zero-length Token on raw message '%v'\n",
					tl.RawString())
			}
		default:
			db.LogIssue("unknown token used",
				fmt.Sprintf("token '%d' is unknown or unhandled.", token.TokenType),
				fmt.Sprintf("Msg: '%s'\nToken: '%+v'", tl.RawString(), token))

			log.Printf("Unknown token used '%d' on raw message '%v'\n",
				token.TokenType, tl.RawString())
		}

	}

	return output
}

var tokenReplacer = strings.NewReplacer(
	string([]rune{rune(0x00)}), "",
	"&lt;", "",
	"&gt;", "",
	"&amp;", "&",
)

// Tokenize will take the incoming string and provide
// a TokenList to work with for things like markovs
func Tokenize(str string) TokenList {
	tokens := make([]ChatToken, 0)

	str = tokenReplacer.Replace(str)

	runes := []rune(str)

	strLen := len(runes)

	peek := func(idx int) rune {
		if idx < 0 {
			return rune(0x00)
		}

		if idx >= strLen {
			return rune(0x00)
		}

		return runes[idx]
	}

	peekUntil := func(pos int, until rune) int {
		for p := pos; p <= strLen; p++ {
			//log.Printf("peekUntil: p=%d,max=%d,chr='%c',until='%c' (%d)\n",
			//	p, strLen, peek(p), until, until)
			if peek(p) == until {
				return p
			}
		}

		return -1
	}

	getRange := func(start, end int) []rune {
		if end < 0 || start < 0 || end <= start || start > len(runes) || end > len(runes) {
			panic(fmt.Sprintf("OHGOD: '%v' (start=%d,end=%d,maxLen=%d)", str, start, end, strLen))
		}

		return runes[start:end]
	}

	inSpecialToken := false
	isMention := false
	isChannel := false

	tokenStart := 0
	lastSpace := 0

	emoticons := make([]string, 0)

	var chr rune
	for pos := 0; pos <= strLen; pos++ {
		chr = peek(pos)
		//x := string([]rune{chr})
		//log.Printf("x: %v\n", x)
		//log.Printf("chr: '%c' (%d), pos='%d', lastSpace: %d, tokenStart: %d, progress: '%s'\n",
		//	chr, chr, pos, lastSpace, tokenStart, string(runes[0:pos]))

		switch chr {
		case rune(0x00):
			if inSpecialToken {
				tokens = append(tokens, ChatToken{
					TokenType: ChatTokenMalformed,
					Value:     getRange(tokenStart, pos),
				})

				break
			}

			if lastSpace == pos {
				break
			}

			tokens = append(tokens, ChatToken{
				TokenType: ChatTokenWord,
				Value:     getRange(lastSpace, pos),
			})

		case '@':

			if (peek(pos-1) == ' ' || pos-1 < 0) && peek(pos+1) != '@' {
				isMention = true
			}

		case '#':

			if (peek(pos-1) == ' ' || pos-1 < 0) && peek(pos+1) != '#' {
				isChannel = true
			}

		case '`':
			bigCode := false

			if peek(pos+1) == '`' && peek(pos+2) == '`' {
				bigCode = true
			}

			//log.Printf("thingy encountered, bigCode=%v\n", bigCode)

			if !bigCode {
				end := peekUntil(pos+1, '`')

				if end < 0 {
					break
				}

				if (pos + 1) != end {
					tokens = append(tokens, ChatToken{
						TokenType: ChatTokenCode,
						Value:     getRange(pos+1, end),
					})
				} else {
					tokens = append(tokens, ChatToken{
						TokenType: ChatTokenCode,
						Value:     []rune{},
					})
				}

				lastSpace = end + 1
				tokenStart = end + 1
				pos = end
			} else {
				//log.Printf("pos-1 = '%c', pos = '%c', pos+1 = '%c', pos+2 = '%c'\n",
				//	peek(pos - 1), chr, peek(pos + 1), peek(pos + 2))

				pos += 3
				end := peekUntil(pos, '`')

				if end < 0 {
					tokens = append(tokens, ChatToken{
						TokenType: ChatTokenCode,
						Value:     []rune{},
					})
					break
				}

				//log.Printf("end = %d, end-1 = '%c', end = '%c', end+1 = '%c', end+2 = '%c'\n",
				//	end, peek(end - 1), peek(end), peek(end + 1), peek(end + 2))

				for true {
					if peek(end+1) == '`' && peek(end+2) == '`' {
						tokens = append(tokens, ChatToken{
							TokenType: ChatTokenCode,
							Value:     getRange(pos-3, end+3),
						})

						lastSpace = end + 3
						tokenStart = end + 3
						pos = end + 3
						break
					}

					end = peekUntil(end+1, '`')
				}
			}

		case ' ', '\n':
			// We can handle both space and newline in the same
			// way for the most part. Just need to define different tokens.

			if inSpecialToken {
				break
			}

			spcToken := ChatToken{
				TokenType: ChatTokenSpace,
				Value:     []rune{' '},
			}

			if chr == '\n' {
				spcToken.TokenType = ChatTokenNewLine
				spcToken.Value = []rune{'\n'}
			}

			if lastSpace == pos {
				tokens = append(tokens, spcToken)

				lastSpace++
				break
			}

			// Add Space!

			tt := ChatTokenWord
			if isMention || isChannel {
				tt = ChatTokenSpecial
			}

			isMention = false
			isChannel = false

			tokens = append(tokens, ChatToken{
				TokenType: tt,
				Value:     getRange(lastSpace, pos),
			})

			tokens = append(tokens, spcToken)

			lastSpace = pos + 1
		case ':':
			if inSpecialToken {
				break
			}

			// Yeah so, words can end with a colon and _not_ be an emoji.
			if peek(pos+1) == ' ' {
				break
			}

			// Similarly: There's no way emoticons are going to be zero-length
			if peek(pos+1) == ':' {
				//pos++
				break
			}

			colEnd := peekUntil(pos+1, ':')
			spcEnd := peekUntil(pos+1, ' ')

			// This is _probably_ due to both being -1
			if spcEnd == colEnd {
				//				db.LogIssue("cleaner_chat", "spcEnd == colEnd", fmt.Sprintf(`
				//colEnd: %d
				//spcEnd: %d
				//message: '%s'
				//`, colEnd, spcEnd, str))

				if spcEnd != -1 {
					db.LogIssue("cleaner_chat", "spcEnd == colEnd but not -1", fmt.Sprintf(`
colEnd: %d
spcEnd: %d
message: '%s'
`, colEnd, spcEnd, str))
				}

				break
			}

			// Easy bail: there's a space in which means this is
			// not gonna be a thing we care about anymore.
			if spcEnd < colEnd && peek(colEnd+1) != rune(0x00) {
				break
			}

			// At this point: either spcEnd is > colEnd or
			// colEnd is last and 0x00 is actually the end.

			if colEnd > spcEnd {
				spcEnd = colEnd + 1
			}

			if peek(spcEnd-1) == ':' {
				emote := getRange(pos, spcEnd)
				tokens = append(tokens, ChatToken{
					TokenType: ChatTokenEmoticon,
					Value:     emote,
				})

				if !utils.StringInSlice(string(emote), emoticons) {
					emoticons = append(emoticons, string(emote))
				}

				pos = spcEnd - 1
				tokenStart = spcEnd
				lastSpace = spcEnd
			}

		case '<':
			if inSpecialToken {
				tokenStart = pos
			}

			// We're actually in an emoticon, I think.
			if peek(pos+1) == ':' && pos+2 < strLen {
				lastBrace := peekUntil(pos, '>')
				lastColon := peekUntil(pos+2, ':')
				lastSpace := peekUntil(pos, ' ')

				// Yup we're in a discord emoticon
				if lastBrace > lastColon && lastBrace > lastSpace {
					emote := getRange(pos+1, lastColon+1)
					tokens = append(tokens, ChatToken{
						TokenType: ChatTokenEmoticon,
						// Get just the colon emote stuff
						Value: emote,
					})

					if !utils.StringInSlice(string(emote), emoticons) {
						emoticons = append(emoticons, string(emote))
					}

					// Go over the last '>'
					pos = lastBrace + 1
					tokenStart = pos
				}
			} else {
				inSpecialToken = true
				tokenStart = pos
			}
		case '>':

			if inSpecialToken {
				tokens = append(tokens, ChatToken{
					TokenType: ChatTokenSpecial,
					// make sure we exclude the brackets because they are not necessary
					// our raw message input may not actually contain these.
					Value: getRange(tokenStart+1, pos),
				})

				tokenStart = pos + 1
				lastSpace = pos + 1
			}

			inSpecialToken = false
			// end of case
		} // end of switch
	} // end of for

	outgoing := make([]ChatToken, 0)
	for _, token := range tokens {
		res := token.Analyze()

		if len(res) > 1 {
			for _, tt := range res {
				if tt.TokenType == ChatTokenEmoticon && !utils.StringInSlice(string(tt.Value), emoticons) {
					emoticons = append(emoticons, string(tt.Value))
				}

				outgoing = append(outgoing, tt)
			}
		} else {
			if res[0].TokenType == ChatTokenEmoticon && !utils.StringInSlice(string(res[0].Value), emoticons) {
				emoticons = append(emoticons, string(res[0].Value))
			}

			outgoing = append(outgoing, res[0])
		}
	}

	return TokenList{
		rawMsg:    str,
		Tokens:    outgoing,
		Emoticons: emoticons,
		URLs:      make([]string, 0),
		Users:     make([]string, 0),
		Channels:  make([]string, 0),
	}
}
