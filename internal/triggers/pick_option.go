package triggers

import (
	"regexp"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

/**
 * this is for the situation where someone asks dungar to choose
 * in a list of things, and we decide to ignore that list and tell them
 * something else.
 */
var alternativeOptions = []weightedChoice{
	{0.30, "let $randomNick$ decide"},
	{0.30, "how about $randomFood$ instead?"},
	{0.10, "checkbox, voted all"},
	{0.10, "no checkbox, didn't vote"},
}

const pickOptionRegex = ":?\\s*(.+\\s+or\\s+.+)\\??"

var (
	splitterRegexp = regexp.MustCompile("\\s*(?:,?\\sor\\s|,)\\s*")
	directedToRgx  = regexp.MustCompile("^(?:@[^: ][^ ]+\\s+|<@!?[^>]+(?:\\|[^>]+)?>\\s+|[^: ][^ ]+:\\s*)")
)

func pickOptionHandler(str, serverID string) string {
	if prefix := directedToRgx.FindString(str); prefix != "" {
		str = strings.TrimSpace(strings.TrimPrefix(str, prefix))
	}

	options := splitterRegexp.Split(str, -1)
	for idx := range options {
		options[idx] = strings.TrimSpace(strings.Trim(options[idx], "?!"))
	}

	if len(options) == 2 {
		return random.PickString(options)
	}

	// we decide to ignore their dumb choices and roll with our own choices.
	if fromBasicChance("pickOptionHandler--ignoreUser") {
		return buildPostMessage(pickWeightedChoice(alternativeOptions), serverID)
	}

	return strings.Trim(random.PickString(options), " ?!")
}
