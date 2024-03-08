package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"time"
)

func commentateScheduler(svc *core2.Service) []*core2.ScheduledMessage {
	if fromBasicChance("commentateScheduler--skip") {
		return nil
	}

	msg := markovGenerate(markovPickWord())

	return []*core2.ScheduledMessage{
		core2.MakeScheduledMessage("any", msg, time.Now()),
	}
}

func commentateHandler(msg *core2.IncomingMessage) []*core2.Response {
	if fromBasicChance("commentateHandler--commentate") {
		return core2.MakeSingleRsp(
			markovGenerate(pickRandomWord(msg.Contents, true)),
		)
	}

	return core2.EmptyRsp()
}

func isQuote(msg *core2.IncomingMessage) bool {
	runes := []rune(msg.Contents)

	return len(runes) > 3 && runes[0] == '>'
}
