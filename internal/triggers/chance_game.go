package triggers

import (
	"regexp"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var chanceTable = []weightedChoice{
	{0.50, "c"},
	{0.60, "d"},
	{0.05, "v0v"},
	{0.05, "idk"},
}

var chanceRegex = regexp.MustCompile(".+ [cC]/[dD]\\??$")

func chanceGameHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if !chanceRegex.MatchString(msg.Contents) {
		return core2.EmptyRsp()
	}

	choice := pickWeightedChoice(chanceTable)

	if isDirectedAtDungar(svc, msg) {
		return core2.PrefixedSingleRsp(choice)
	}

	if fromBasicChance("chanceGameHandler--basic") {
		return core2.MakeSingleRsp(choice)
	}

	return core2.EmptyRsp()
}
