package cleaner

import "strings"

// TokenList is our listing struct of tokens
type TokenList struct {
	tokens []Token
	simple *SimpleTokenList
}

// String provides the raw values back at the user (rebuilt)
func (tl *TokenList) String() string {
	out := ""
	for _, t := range tl.tokens {
		out += string(t.Value)
	}

	return out
}

// TokenStr returns the list of tokens in a string (delimited by space)
func (tl *TokenList) TokenStr() string {
	out := make([]string, len(tl.tokens))
	for p, t := range tl.tokens {
		out[p] = t.Type.String()
	}

	return strings.Join(out, " ")
}

// GetTokens returns the internal list of tokens
func (tl *TokenList) GetTokens() []Token {
	return tl.tokens
}

// DebugPrint prints the list of tokens to STDERR (via log.Printf)
func (tl *TokenList) DebugPrint() {
	DebugPrintTokenList(tl.tokens)
}

// IsMalformed returns whether or not we found a token that appears
// malformed.
func (tl *TokenList) IsMalformed() bool {
	for _, token := range tl.tokens {
		if token.Type == TokenMalformed {
			return true
		}
	}

	return false
}
