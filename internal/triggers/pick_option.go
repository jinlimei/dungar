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

var alternativeOptionsForTwo = []weightedChoice{
	{0.30, "por que no los dos?"},
	{0.30, "yes"},
	{0.30, "i'm leaning towards the former but honestly the latter tho"},
	{0.10, "have you considered picking up chess instead?"},
}

const pickOptionRegex = ":?\\s*(.+\\s+or\\s+.+)\\??"

var (
	splitterRegexp = regexp.MustCompile("\\s*(?:,?\\sor\\s|,)\\s*")
	directedToRgx  = regexp.MustCompile("^(@[^:][^ ]+|[^:][^ ]+):?\\s*")
)

func pickOptionHandler(str, serverID string) string {
	if directedToRgx.MatchString(str) {
		sub := directedToRgx.FindStringSubmatch(str)
		str = strings.TrimSpace(strings.Replace(str, sub[0], "", -1))
	}

	options := splitterRegexp.Split(str, -1)

	// we decide to ignore their dumb choices and roll with our own choices.
	if fromBasicChance("pickOptionHandler--ignoreUser") {
		// If it's only two options we can build off of a different list.
		if len(options) == 2 && fromBasicChance("pickOptionHandler--altOptions") {
			return buildPostMessage(pickWeightedChoice(alternativeOptionsForTwo), serverID)
		}

		return buildPostMessage(pickWeightedChoice(alternativeOptions), serverID)
	}

	return strings.Trim(random.PickString(options), " ?!")
}
