package triggers

import "gitlab.int.magneato.site/dungar/prototype/library/core2"

var userSpammed = []weightedChoice{
	{0.02, "behold! A rare performance in which $nick$ attempts intelligent conversation through the use of word vomit"},
	{0.08, "We get it, bitcoin crashed. it's fine."},
	{0.30, "tl;dr"},
	{0.30, "stop trying to spam me, nerds"},
	{0.30, ":words:"},
}

var botSpammed = []weightedChoice{
	{0.03, "MY WALL OF SPAM CRITS YOU FOR 999,999,999 DAMAGE YOU DIE"},
	{0.07, "Wow I'll shut up now"},
	{0.10, "GOTTEM"},
	{0.10, "GET REKT"},
	{0.10, "UGA WAKA TEKA TAKA SPAGHETTI"},
	{0.60, "oops wrong window"},
}

func userSpammedHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if msg.UserID == svc.GetBotUser().ID {
		return core2.EmptyRsp()
	}

	if !fromBasicChance("spamHandler--respond") {
		return core2.EmptyRsp()
	}

	symbols := []rune(msg.Contents)

	if len(symbols) >= 768 {
		rsp := core2.MakeRsp(pickWeightedChoice(userSpammed))
		rsp.ConsumedMessage = false

		return core2.SingleRsp(rsp)
	}

	return core2.EmptyRsp()
}

func selfSpamHandler(svc *core2.Service, rsps []*core2.Response) []*core2.Response {
	additions := make([]*core2.Response, 0)

	for _, rsp := range rsps {
		symbols := []rune(rsp.Contents)
		if len(symbols) < 1024 {
			continue
		}

		//if basic.Execute("selfSpamHandler--skip", 0.10) {
		//	continue
		//}

		additions = append(additions, &core2.Response{
			ConsumedMessage:  true,
			HandledMessage:   true,
			PrefixUsername:   false,
			CancelledMessage: false,
			Contents:         pickWeightedChoice(botSpammed),
			ResponseType:     core2.ResponseTypeBasic,
		})
	}

	return additions
}
