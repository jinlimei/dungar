package triggers

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var rollDiceRegex = regexp.MustCompile("\\s+((\\d+)[Dd](\\d+)([hlHLoO](\\d+)|))")

var overDiceRollsLimitResponses = []weightedChoice{
	{0.80, "perhaps try using 5 or less different rolls"},
	{0.10, "your outlook doesn't look too good fam"},
	{0.10, "instead of dice rolls why not try bitcoin?"},
}

func rollDiceHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if !isDirectedAtDungar(svc, msg) {
		return core2.EmptyRsp()
	}

	if !rollDiceRegex.MatchString(msg.Contents) {
		return core2.EmptyRsp()
	}

	matches := rollDiceRegex.FindAllStringSubmatch(msg.Contents, -1)
	output := ""

	max := len(matches)

	if max > 5 {
		return core2.PrefixedSingleRsp(pickWeightedChoice(overDiceRollsLimitResponses))
	}

	//spew.Dump(matches)
	for pos, matchGroup := range matches {
		diceRolls, _ := strconv.Atoi(matchGroup[2])
		diceSize, _ := strconv.Atoi(matchGroup[3])

		if diceSize <= 1 || diceRolls <= 0 || diceSize > 10000 || diceRolls > 10 {
			output += matchGroup[1] + ": nah"
		} else {
			output += generateDiceRolls(matchGroup, diceRolls, diceSize)
		}

		if (pos + 1) < max {
			output += "; "
		}
	}

	output = strings.TrimSpace(output)

	return core2.PrefixedSingleRsp(output)
}

func generateDiceRolls(matchGroup []string, diceRolls, diceSize int) string {
	output := matchGroup[1] + ": "

	//spew.Dump(matchGroup)
	rolls := make([]int, 0)

	for k := 0; k < diceRolls; k++ {
		val := random.Int(diceSize) + 1
		rolls = append(rolls, val)
	}

	if matchGroup[4] != "" {
		lim, _ := strconv.Atoi(matchGroup[5])

		if lim <= 0 || lim >= diceRolls {
			return output + "nah"
		}

		sort.Ints(rolls)

		nRolls := make([]int, 0)

		if matchGroup[4][0] == 'h' || matchGroup[4][0] == 'H' {
			pos := len(rolls) - 1
			for x := 0; x < lim; x++ {
				nRolls = append(nRolls, rolls[pos])
				pos--
			}
		} else if matchGroup[4][0] == 'l' || matchGroup[4][0] == 'L' {
			for x := 0; x < lim; x++ {
				nRolls = append(nRolls, rolls[x])
			}
		} else {
			midStart := int(len(rolls)/2) - lim
			if midStart < 0 {
				midStart = 0
			}

			for x := 0; x < lim; x++ {
				nRolls = append(nRolls, rolls[midStart+x])
			}
		}

		rolls = nRolls
	}

	for k, roll := range rolls {
		output += strconv.Itoa(roll)

		if (k + 1) < len(rolls) {
			output += ", "
		}
	}

	return output
}
