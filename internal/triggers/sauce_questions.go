package triggers

import "gitlab.int.magneato.site/dungar/prototype/library/core2"

var complexTriggers = []multiTriggerCallback{
	buildMultiQuestions(bitCoinHandler, "bitcoin", "btc"),
	buildMultiQuestions(jsHandler, "javascript", "js", "node"),
}

var simpleAnywhereTriggers = []triggerChoice{
	{Trigger: "mercy", Response: "mercy main btw"},
	{Trigger: "hanzo", Response: "I for one always switch to hanzo when we don't have tanks"},
	{Trigger: "tax|taxes", Response: "if you don't pay your taxes, you steal from society", IsRegex: true},
	{Trigger: "daddy", Response: "just stop"},
}

func sauceQuestionsHandler(msg *core2.IncomingMessage) []*core2.Response {
	txt := msg.Lowered()

	if fromBasicChance("sauceQuestionsHandler--consume") {
		for _, trigger := range complexTriggers {
			if trigger.matches(txt) {
				return core2.MakeSingleRsp(trigger.Handler(txt, msg.ServerID))
			}
		}

		for _, trigger := range simpleAnywhereTriggers {
			if trigger.matches(txt) {
				return core2.MakeSingleRsp(trigger.Response)
			}
		}
	}

	return core2.EmptyRsp()
}
