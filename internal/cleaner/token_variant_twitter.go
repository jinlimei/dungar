package cleaner

import "strings"

func (t *Token) parseWordForTwitter() []Token {
	var (
		val = string(t.Value)
	)

	if val == "RT" {
		t.Type = TokenSpecial
	} else if strings.HasPrefix(val, "https://") ||
		strings.HasPrefix(val, "http://") {
		t.Type = TokenURL
	} else if strings.HasPrefix(val, "@") {
		t.Type = TokenMentionUser
	} else if strings.HasPrefix(val, "#") {
		t.Type = TokenHashTag
	}

	return []Token{*t}
}
