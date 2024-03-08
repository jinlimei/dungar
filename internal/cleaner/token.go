package cleaner

import (
	"fmt"
	"log"
)

// TokenType is our base type for handling tokens
//
//go:generate stringer -type=TokenType
type TokenType uint8

const (
	// TokenUndefined is a way of capturing zero-value tokens
	TokenUndefined TokenType = iota
	// TokenSpace is a space
	TokenSpace
	// TokenWord is a basic word
	TokenWord
	// TokenNumber is 0123456789
	TokenNumber
	// TokenNewLine represents \n
	TokenNewLine
	// TokenDoubleQuote represents '"'
	TokenDoubleQuote
	// TokenSingleQuote represents "'"
	TokenSingleQuote
	// TokenOpenParens is a '('
	TokenOpenParens
	// TokenCloseParens is a ')'
	TokenCloseParens
	// TokenOpenBracket is a '{'
	TokenOpenBracket
	// TokenCloseBracket is a '}'
	TokenCloseBracket
	// TokenOpenBrace is a '['
	TokenOpenBrace
	// TokenCloseBrace is a ']'
	TokenCloseBrace
	// TokenPeriod represents '.'
	TokenPeriod
	// TokenSentEnd represents any non-period end of sentences
	// so: !, ?, some assortment?
	TokenSentEnd
	// TokenEllipsis represents '...'
	TokenEllipsis
	// TokenEmoticon represents colon-based emoticons
	TokenEmoticon
	// TokenURL represents URLs in the results
	TokenURL
	// TokenCode is code in a token. For coding.
	TokenCode
	// TokenCommand is for things like !dadjoke
	TokenCommand
	// TokenSpecial is special formatting for a specific variant
	TokenSpecial
	// TokenEOL is end-of-line
	TokenEOL
	// TokenMalformed is when something went wrong
	TokenMalformed
	// TokenMentionUser is a @mention
	TokenMentionUser
	// TokenMentionChannel is a #mention
	TokenMentionChannel
	// TokenMentionRole is a @&mention (for discord)
	TokenMentionRole
	// TokenHashTag is a #hashtag (for twitter)
	TokenHashTag
)

// Token is our complicated struct for building tokens from a string
type Token struct {
	Type     TokenType
	StartPos int
	EndPos   int
	Value    []rune
	Override []rune

	hash *uint64
}

// ValueOrOverride will return a string of Token.Override if it's non-empty or of Token.Value
func (t *Token) ValueOrOverride() string {
	if len(t.Override) > 0 {
		return string(t.Override)
	}

	return string(t.Value)
}

// String takes the Token.Value and converts it to a string
func (t *Token) String() string {
	return string(t.Value)
}

// Hash returns a CRC64 hash of the Token (type, value, override)
func (t *Token) Hash() uint64 {
	if t.hash == nil {
		done := hash(t)
		t.hash = &done
	}

	return *t.hash
}

// DebugString provides a string useful for debugging
func (t *Token) DebugString() string {
	ovr := ""
	if len(t.Override) > 0 {
		ovr = fmt.Sprintf(" ovr:'%s'", string(t.Override))
	}

	return fmt.Sprintf(
		`%s (%3d, %3d) "%s"%v %v`,
		rpad(t.Type.String(), 17),
		t.StartPos,
		t.EndPos,
		string(t.Value),
		ovr,
		t.Value,
	)
}

// Parse looks at the TokenVariant and returns a new list of Token's to be
// used instead of this specific token.
func (t *Token) Parse(variant TokenVariant) []Token {

	switch t.Type {
	case TokenSpace, TokenNumber, TokenNewLine, TokenEOL, TokenEmoticon:
		return []Token{*t}
	case TokenWord:
		return t.parseTokenWord(variant)
	case TokenSpecial:
		return t.parseTokenSpecial(variant)
	case TokenMalformed:
		// We still need to provide information that it was malformed, but at least
		// we can continue on!
		return []Token{
			{Type: TokenMalformed, StartPos: 0, EndPos: 0, Value: nil, Override: nil},
		}

	default:
		log.Printf("Have not implemented '%s' token yet for token.Parse\n",
			t.Type.String())
		return []Token{*t}
	}
}

func (t *Token) parseTokenWord(variant TokenVariant) []Token {
	switch variant {
	case VariantSlack, VariantDiscord, VariantMatrix:
		return t.parseWordForChat(variant)
	case VariantXMPP:
		return t.parseTokenForXMPP()
	case VariantTwitter:
		return t.parseWordForTwitter()
	default:
		return t.parseWordForPlain()
	}
}

func (t *Token) parseTokenSpecial(variant TokenVariant) []Token {
	switch variant {
	case VariantXMPP:
		return t.parseTokenForXMPP()
	case VariantPlain, VariantBook:
		t.Type = TokenWord
		return []Token{*t}
	case VariantSlack:
		return t.parseSpecialForSlack()
	case VariantDiscord:
		return t.parseSpecialForDiscord()
	default:
		return []Token{*t}
	}
}

func isNumeric(r []rune) bool {
	var (
		numCount = 0
	)

	for k := 0; k < len(r); k++ {
		switch r[k] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '%':
			numCount++
		}
	}

	return numCount == len(r)
}

func isMaybeTimestamp(r []rune) bool {
	var (
		colCount = 0
		numCount = 0
	)

	for k := 0; k < len(r); k++ {
		switch r[k] {
		case ':':
			colCount++
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			numCount++
		}
	}

	// well okay it's only numbers and colons get the heck out
	if colCount > 0 && colCount+numCount == len(r) {
		return true
	}

	// maybe less likely? like a :4:2:
	if colCount == 2 && numCount > 2 {
		return true
	}

	// idk maybe not
	return false
}
