package learning

import (
	"log"
	"strings"
)

// ReadABook converts a book written in a string into a line of sentences
// that we can easily operate on.
func ReadABook(str string) []string {
	str = cleanupSpace(str)
	str = cleanupQuotes(str)
	str = strings.ToLower(str)

	lines := bookStartCleanup(str)
	lines = bookFinalCleanup(lines)

	return lines
}

func getRuneRange(runes []rune, start, end int) string {
	if start < 0 || start >= end {
		return ""
	}

	if end > len(runes) {
		log.Printf("getRuneRange end (%d) > len (%d)\n", end, len(runes))
		end = len(runes)
	}

	var (
		//rLen = len(runes)
		out = make([]rune, 0, end-start+1)
		pos = start
		idx = 0
	)

appender:

	out = append(out, runes[pos])
	idx++
	pos++

	if pos < end {
		goto appender
	}

	return string(out)
}

func bookStartCleanup(str string) []string {
	var (
		runes  = []rune(str)
		output = make([]string, 0)
		max    = len(runes)
		step   = 0
		pos    = 0

		chr  rune
		peek rune
	)

	for ; pos <= max; pos++ {
		if pos >= max {
			chr = rune(0x00)
		} else {
			chr = runes[pos]
		}

		if pos+1 >= max {
			peek = rune(0x00)
		} else {
			peek = runes[pos+1]
		}

		//log.Printf("%04d: %v (%v)\n", pos, string(chr), int64(chr))
		if (chr == '!' || chr == '?' || chr == '.') && (peek == ' ' || peek == rune(0x00)) {
			r := getRuneRange(runes, step, pos+1)
			if !isAbbreviation(r) {
				output = append(output, cleanupLines(r))
				step = pos + 2
				pos += 2
			}

		} else if chr == rune(0x00) {
			output = append(output, cleanupLines(getRuneRange(runes, step, pos)))
		}

	} // end of for

	return output
}

func getLastWord(s string) []rune {
	var (
		srs = []rune(s)
		pos = len(srs) - 1
	)

	hasSpace := false
	for ; pos >= 0; pos-- {
		if srs[pos] == ' ' {
			hasSpace = true
			break
		}
	}

	if !hasSpace {
		out := make([]rune, 0, len(srs))

		out = append(out, ' ')
		return append(out, srs...)
	}

	return srs[pos:]
}

var abbreviations = [][]rune{
	{' ', 's', 't', '.'},
	{' ', 'm', 'r', '.'},
	{' ', 'd', 'r', '.'},
	{' ', 'm', 's', '.'},
	{' ', 'm', 'r', 's', '.'},
	{' ', 'l', 'n', '.'},
	{' ', 'a', 'v', 'e', '.'},
	{' ', 'e', 's', 't', '.'},
	{' ', 'd', 'e', 'p', 't', '.'},
	{' ', 'n', 'o', '.'},
	{' ', 'v', 's', '.'},
	{' ', 'v', '.'},
}

func isAbbreviation(s string) bool {
	word := getLastWord(s)

	for _, abbr := range abbreviations {
		if len(abbr) == len(word) {
			matched := true
			for k := 0; k < len(word); k++ {
				if word[k] != abbr[k] {
					matched = false
					break
				}
			}

			if matched {
				return true
			}
		}
	}

	return false
}

func bookFinalCleanup(lines []string) []string {
	output := make([]string, 0)

	for _, line := range lines {
		for strings.Contains(line, "  ") {
			line = strings.ReplaceAll(line, "  ", " ")
		}

		if line == "" {
			continue
		}

		if strings.Contains(line, "* ") {
			continue
		}

		output = append(output, line)
	}

	return output
}
