package markov3

import (
	"slices"
	"strings"
)

var repl = strings.NewReplacer(
	"<", "",
	">", "",
	"&lt;", "<",
	"&gt;", ">",
	"&amp;", "\"",
	"_", "",
	"*", "",
	"‘", "'",
	"’", "'",
	"“", "\"",
	"”", "\"",
)

func cleanWordTrimFunc(r rune) bool {
	return r == '\\' ||
		r == ',' ||
		r == ':' ||
		r == '/' ||
		r == '"' ||
		r == ' ' ||
		r == '‘' ||
		r == '’' ||
		r == '“' ||
		r == '”' ||
		r == '\''
}

func cleanWord(w string) string {
	w = repl.Replace(w)
	w = strings.TrimFunc(w, cleanWordTrimFunc)

	return w
}

func reverseMarkovIDs(arr []MarkovID) []MarkovID {
	out := make([]MarkovID, len(arr))

	fwd := 0

	for _, a := range slices.Backward(arr) {
		out[fwd] = a
		fwd++
	}

	return out
}

func joinMarkovIDs(listing ...[]MarkovID) []MarkovID {
	out := make([]MarkovID, 0)

	for _, list := range listing {
		out = append(out, list...)
	}

	return out
}
