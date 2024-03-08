package utils

import (
	"regexp"
	"strings"
	"unicode"
)

var urlRegex = regexp.MustCompile(`^[a-z]+\.[a-z0-9]+(.*/?|)`)
var hasWordChar = regexp.MustCompile("\\W")
var punctuationStrip = regexp.MustCompile("^\\W*(.+?)\\W*$")

// IsEmoticon makes a determination if the incoming word is in fact of an emoticon nature
func IsEmoticon(word string) bool {
	word = strings.TrimSpace(word)
	wordLen := len(word)

	if wordLen > 2 && word[0] == ':' && word[wordLen-1] == ':' && strings.Index(word, " ") == -1 {
		return true
	}

	return false
}

// Normalize will return a word normalized (lowered, punctuation we don't like stripped)
func Normalize(s string) string {
	res := s

	if hasWordChar.MatchString(s) {
		res = punctuationStrip.ReplaceAllString(s, "$1")
	}

	return strings.ToLower(res)
}

// EndsWithPunctuation Returns whether or not a word ends in punctuation that we care about
func EndsWithPunctuation(word string) bool {
	if word == "" {
		return false
	}

	lastChar := word[len(word)-1]

	switch lastChar {
	case '!', '.', ':', ';', ',', '?', '~', ')', '}', ']', '/', '"', '\'':
		return true
	default:
		return false
	}
}

// TrimPunctuation removes any punctuation we don't
// want to particularly care about
func TrimPunctuation(s string) string {
	s = strings.Trim(s, "\t\n\r ?!.,'\"()[]{}~”")

	if len(s) > 1 && s[0] == '@' {
		return s[1:]
	}

	return s
}

// IsURL returns whether or not the string is a URL
func IsURL(str string) bool {
	if len(str) < 8 {
		return urlRegex.MatchString(str)
	}

	return str[0:7] == "http://" || str[0:8] == "https://" || urlRegex.MatchString(str)
}

// StringToWords provides a mechanism to convert a string to simple words
// and to normalize them if specified.
func StringToWords(str string, normalize bool) []string {
	if !strings.Contains(str, " ") {
		if normalize {
			str = Normalize(str)
		}

		return []string{str}
	}

	words := strings.Split(str, " ")
	for idx, word := range words {
		if normalize && !IsURL(word) {
			word = Normalize(word)
		}

		words[idx] = word
	}

	return words
}

// StringInSlice looks for str in slice
func StringInSlice(str string, slice []string) bool {
	if len(slice) == 0 {
		return false
	}

	if slice[0] == str {
		return true
	}

	for _, msg := range slice {
		if msg == str {
			return true
		}
	}

	return false
}

// CleanSpaces will remove all double spaces.
func CleanSpaces(str string) string {
	hasWeirdSpaces := strings.Contains(str, "  ")
	for hasWeirdSpaces {
		str = strings.ReplaceAll(str, "  ", " ")
		hasWeirdSpaces = strings.Contains(str, "  ")
	}

	return strings.TrimSpace(str)
}

// CoalesceStr coalesces until it finds a non-empty string.
func CoalesceStr(strings ...string) (string, bool) {
	if len(strings) == 0 {
		return "", false
	}

	if len(strings) == 1 {
		return strings[0], strings[0] != ""
	}

	ok := false
	ret := ""

	for _, str := range strings {
		if str != "" {
			ret = str
			ok = true
			break
		}
	}

	return ret, ok
}

// TitleCase is based on the current strings.Title function
// but this one isn't deprecated because we don't care about
// weird unicode word boundaries since this bot is focused on
// the English language.
func TitleCase(s string) string {
	prev := ' '
	return strings.Map(func(r rune) rune {
		if prev == ' ' {
			prev = r
			return unicode.ToTitle(r)
		}

		prev = r
		return r
	}, s)
}
