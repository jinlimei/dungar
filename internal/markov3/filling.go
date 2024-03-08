package markov3

import "gitlab.int.magneato.site/dungar/prototype/internal/cleaner"

var fillerMap = map[cleaner.TokenType]string{
	cleaner.TokenURL:            "\u0000URL",
	cleaner.TokenEmoticon:       "\u0000Emoticon",
	cleaner.TokenMentionUser:    "\u0000User",
	cleaner.TokenMentionChannel: "\u0000Channel",
	cleaner.TokenMentionRole:    "\u0000Role",
	cleaner.TokenHashTag:        "\u0000HashTag",
}

var revFillerMap = map[string]cleaner.TokenType{
	"\u0000URL":      cleaner.TokenURL,
	"\u0000Emoticon": cleaner.TokenEmoticon,
	"\u0000User":     cleaner.TokenMentionUser,
	"\u0000Channel":  cleaner.TokenMentionChannel,
	"\u0000Role":     cleaner.TokenMentionRole,
	"\u0000HashTag":  cleaner.TokenHashTag,
}

func (m *Markov) initInitialFillers() {
	m.tokenFillers = map[cleaner.TokenType][]string{
		cleaner.TokenURL:            {"tedcruzforhumanpresident.com"},
		cleaner.TokenEmoticon:       {":parrot:"},
		cleaner.TokenMentionChannel: {"butts"},
		cleaner.TokenMentionUser:    {"chalur"},
		cleaner.TokenMentionRole:    {"president"},
		cleaner.TokenHashTag:        {"#notmypresident"},
	}
}

func (m *Markov) needsFiller(tt cleaner.TokenType) bool {
	return tt == cleaner.TokenURL ||
		tt == cleaner.TokenEmoticon ||
		tt == cleaner.TokenMentionChannel ||
		tt == cleaner.TokenMentionUser ||
		tt == cleaner.TokenMentionRole ||
		tt == cleaner.TokenHashTag
}

func (m *Markov) fillerConversion(tt cleaner.TokenType) string {
	out, ok := fillerMap[tt]
	if !ok {
		return ""
	}

	return out
}

// SetFiller replaces all filler strings for ChatSubTokenType
func (m *Markov) SetFiller(t cleaner.TokenType, s []string) {
	m.tokenFillers[t] = s
}

// AddFiller adds a filler string to the filler list
func (m *Markov) AddFiller(t cleaner.TokenType, s string) {
	m.tokenFillers[t] = append(m.tokenFillers[t], s)
}
