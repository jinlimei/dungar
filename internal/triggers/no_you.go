package triggers

import (
	"regexp"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	noYouFuckYou = regexp.MustCompile("^fu(ck|x|k) (you|ur?|off|that)")
	noYouUrMom   = regexp.MustCompile("^(your|ur|you're|yor) (mom|mum|mommy|dad|daddy|a)")

	noYouResponses = []string{
		"no you",
		"no u",
		"NO U",
		"NO YOU",
	}
)

func noYouHandler(msg *core2.IncomingMessage) []*core2.Response {
	txt := msg.Lowered()

	if fromBasicChance("noYouHandler--regexp") &&
		(noYouUrMom.MatchString(txt) || noYouFuckYou.MatchString(txt)) {
		return core2.MakeSingleRsp(random.PickString(noYouResponses))
	}

	if fromBasicChance("noYouHandler--normal") {
		return core2.MakeSingleRsp(random.PickString(noYouResponses))
	}

	return core2.EmptyRsp()
}
