package cleaner

import "gitlab.int.magneato.site/dungar/prototype/internal/utils"

func (t *Token) parseTokenForXMPP() []Token {
	// Don't need that shit in the logs
	if len(t.Value) > 1 && t.Value[0] == '<' && t.Value[len(t.Value)-1] == '>' {
		return nil
	}

	if len(t.Value) > 1 && (t.Value[0] == '!' || t.Value[0] == '$') {
		t.Type = TokenCommand
		return []Token{*t}
	}

	if len(t.Value) > 1 && t.Value[0] == '@' {
		t.Type = TokenMentionUser
		return []Token{*t}
	}

	if isNumeric(t.Value) {
		t.Type = TokenNumber
		return []Token{*t}
	}

	if isMaybeTimestamp(t.Value) {
		return nil
	}

	str := string(t.Value)

	if utils.IsURL(str) {
		t.Type = TokenURL
		return []Token{*t}
	}

	return []Token{*t}
}
