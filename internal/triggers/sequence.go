package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

type sequenceMessage []string

var funSequences = []sequenceMessage{
	//{"oof", "ouch", "owie"},
	//{"haha", "yes"},
	{"me", "too", "thanks"},
	{"rip", ":rip:", ":911:", "F"},
}

var randSequences = []sequenceMessage{
	{":thinking:", ":thonkeng:", ":thonking:", ":highthonk:", ":mercythonk:", ":dodecathonk:", ":thoink:",
		":thinkio:", ":thinking_face:", ":thinkcat:", ":thinkzen:", ":thinking:", ":thinkpad:", ":ohno_think:",
		":thinkhand:", ":thinkspin:", ":thinkbear:", ":thinklego:", ":mercythink:", ":thinkskype:", ":thinkingfast:",
		":safarithink:", ":thinkferret:", ":firefoxthink:", ":thinkparrot:", ":think_different:", ":blinging:",
		":universing:", ":thinklab:", ":thinktao:", ":superthink:", ":thinkachu:", ":birdthink:", ":catgirlthink:",
		":thinkpartyblob:", ":mrthinkwide:", ":blackholethink:"},
	{":ohdear:", ":ohdearsass:", ":dodecaohdear:"},
	{":911:", ":texas:", ":canada:", ":norway:", ":mexico:", ":britain:", ":france:"},
}

func (sm sequenceMessage) next(str string) string {
	strLen := len(sm)

	for i := 0; i < strLen; i++ {
		if sm[i] == str {
			return sm[((i + 1) % strLen)]
		}
	}

	return ""
}

func (sm sequenceMessage) random(str string) string {
	strLen := len(sm)

	for i := 0; i < strLen; i++ {
		if sm[i] == str {
			return random.PickString(sm)
		}
	}

	return ""
}

func sequenceHandler(msg *core2.IncomingMessage) []*core2.Response {
	var (
		txt    = msg.Lowered()
		seqRsp string
		tmp    string
	)

	for _, seq := range funSequences {
		tmp = seq.next(txt)
		if tmp != "" {
			seqRsp = tmp
			break
		}
	}

	if seqRsp == "" {
		for _, seq := range randSequences {
			tmp = seq.random(txt)
			if tmp != "" {
				seqRsp = tmp
				break
			}
		}
	}

	if seqRsp != "" && fromBasicChance("sequenceHandler--respond") {
		return core2.MakeSingleRsp(seqRsp)
	}

	return core2.EmptyRsp()
}
