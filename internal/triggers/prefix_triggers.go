package triggers

import "gitlab.int.magneato.site/dungar/prototype/library/core2"

var simplePrefixTriggers = []triggerChoice{
	{Trigger: "!dogho", Response: "!DOGHO"},
	{Trigger: "fuck you", Response: "no u"},
	{Trigger: "shazbot", Response: "vgs"},
	{Trigger: "vgs", Response: "shazbot"},
	{Trigger: "vgtg", Response: "i am the greatest!"},
	{Trigger: "i am the greatest", Response: "vgtg"},
	{Trigger: "jinlitar", Response: "afk lol"},
	{Trigger: "dovibus", Response: "dead lol also fuck cancer"},
	{Trigger: "icosabus", Response: "let me tell you about selling your erebus for $int$ ragnaroks and $int$ hels"},
	{Trigger: "dodecabus", Response: "you do mean icosabus right? :colbert:"},
	{Trigger: "mad dovi", Response: ":mad:"},
	{Trigger: "dubstep", Response: "no"},
}

func simplePrefixHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !fromBasicChance("simplePrefixHandler--basic") {
		return core2.EmptyRsp()
	}

	for _, trigger := range simplePrefixTriggers {
		if trigger.matches(msg.Contents) {
			return core2.MakeSingleRsp(trigger.Response)
		}
	}

	return core2.EmptyRsp()
}
