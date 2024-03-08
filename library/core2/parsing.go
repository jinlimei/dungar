package core2

import (
	"fmt"
	"strconv"
)

// TokenType is an enum for various token types for message tokens
//
//go:generate stringer -type=TokenType
type TokenType int

const (
	// TokenSpace is for spaces (" ")
	TokenSpace TokenType = iota
	// TokenWord is for classic safe words, which have no magic attached.
	TokenWord
	// TokenUserID is for <@UID> tokens
	TokenUserID
	// TokenChanID is for <#CID> tokens
	TokenChanID
	// TokenGroupID is for <#group> or <@group>
	TokenGroupID
	// TokenRoleID is primarily for discord's role identifiers (<@&RID>)
	TokenRoleID
	// TokenEmoticon is primarily for discords use of emoticons:
	// slack: :thing:
	// discord: <:thing:EmoticonID>
	TokenEmoticon
	// TokenURL is for <https://www> tokens
	TokenURL
)

// MessageToken is the various tokens that comprise a received message.
type MessageToken struct {
	Token    string
	Type     TokenType
	Value    *string
	Override *string
}

// ParsedMessage is the incoming/received message parsed into a set
// of tokens which could be reused to rebuild the raw message, or
// to convert into a new message.
type ParsedMessage struct {
	Tokens    []MessageToken
	Raw       string
	Converted string
}

// AddToken is a quickie shortcut to append tokens to the end of Tokens
func (pm *ParsedMessage) AddToken(token MessageToken) {
	pm.Tokens = append(pm.Tokens, token)
}

// IDTokens provides a list of all ID-style (TokenChanID and TokenUserID)
// tokens in the parsed message.
func (pm *ParsedMessage) IDTokens() []MessageToken {
	out := make([]MessageToken, 0)

	for _, tok := range pm.Tokens {
		if tok.Type == TokenChanID || tok.Type == TokenUserID || tok.Type == TokenGroupID || tok.Type == TokenRoleID {
			out = append(out, tok)
		}
	}

	return out
}

// URLTokens provides a list of all URL tokens (TokenURL)
// in the parsed message.
func (pm *ParsedMessage) URLTokens() []MessageToken {
	out := make([]MessageToken, 0)

	for _, tok := range pm.Tokens {
		if tok.Type == TokenURL {
			out = append(out, tok)
		}
	}

	return out
}

// TypeToString converts a TokenType into the string counterpart
// of its name
func (pm *ParsedMessage) TypeToString(t TokenType) string {
	switch t {
	case TokenWord:
		return "TokenWord"
	case TokenSpace:
		return "TokenSpace"
	case TokenURL:
		return "TokenURL"
	case TokenRoleID:
		return "TokenRoleID"
	case TokenUserID:
		return "TokenUserID"
	case TokenChanID:
		return "TokenChanID"
	case TokenGroupID:
		return "TokenGroupID"
	}

	return "UnknownToken" + strconv.Itoa(int(t))
}

// Dump outputs a list of tokens formatted in a nice way for console analysis
func (pm *ParsedMessage) Dump() {
	pad := func(val string, amt int) string {
		vLen := len(val)

		for vLen < amt {
			val += "."
			vLen++
		}

		return val
	}

	safe := func(p *string) string {
		if p == nil {
			return "<nil>"
		}

		return *p
	}

	maxTokenLen := 1

	for _, token := range pm.Tokens {
		if len(token.Token) > maxTokenLen {
			maxTokenLen = len(token.Token)
		}
	}

	fmt.Printf("RAW: '%s'\n", pm.Raw)

	for pos, token := range pm.Tokens {

		fmt.Printf(
			"%03d: TOKEN %s RAW %s (VALUE='%s', OVERRIDE='%s')\n",
			pos,
			pad(fmt.Sprintf("'%s' ", pm.TypeToString(token.Type)), 15),
			pad(fmt.Sprintf("'%s' ", token.Token), maxTokenLen+3),
			safe(token.Value),
			safe(token.Override),
		)
	}
}
