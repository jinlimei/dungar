package markov3

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"strings"
)

// MakeWordSequence will take the incoming string,
// use Tokenize on it, and return the TokenList.MarkovConsumable
// output as our return
func MakeWordSequence(str string, variant cleaner.TokenVariant) cleaner.TokenList {
	if strings.Contains(str, "\n") || strings.Contains(str, "\r") {
		str = strings.ReplaceAll(str, "\r", "")
		str = strings.ReplaceAll(str, "\n", " ")
	}

	for strings.Contains(str, "  ") {
		str = strings.ReplaceAll(str, "  ", " ")
	}

	// A potentially necessary sacrifice:
	// make everything lower case.
	str = strings.ToLower(str)

	return cleaner.Tokenize(str, variant)
}

// safeGetWordID lets us use a -1 and also +(over-max) and not
// have things bork due to out of range errors.
func safeGetWordID(ids []MarkovID, maxLen, pos int) MarkovID {
	if pos < 0 {
		return LeftMost
	}

	if (pos + 1) > maxLen {
		return RightMost
	}

	return ids[pos]
}
