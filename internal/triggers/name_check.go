package triggers

import "gitlab.int.magneato.site/dungar/prototype/library/core2"

var nameOutMouthChoices = []weightedChoice{
	{0.90, "get my name out of your mouth"},
	{0.10, "get my name out ur mouf"},
}

func nameCheckHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if isMentioningDungar(svc, msg) && basic.Execute("nameCheckHandler", 0.005) {
		return core2.PrefixedSingleRsp(pickWeightedChoice(nameOutMouthChoices))
	}

	return core2.EmptyRsp()
}
