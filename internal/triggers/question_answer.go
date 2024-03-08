package triggers

import (
	"regexp"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

var stupidAnswerResponses = []weightedChoice{
	{0.60, "idk do you?"},
	{0.20, ":jiggled:"},
	{0.01, "like and subscribe to find out"},
	{0.01, ":lobster:"},
}

const questionAnswerRegex = ":?\\s*(?:how do|wh?[ua]t do|why do|do)" +
	"\\s+(?:you|u)\\s+" +
	"(?:hate|like|think [^ ]+|love|[^ ]+ about)" +
	"\\s+(.+)\\??"

const stupidQuestionAnswerRegex = ":?\\s*do you (.+?)\\??$"

var questionAnswerCompiled = regexp.MustCompile(questionAnswerRegex)
var stupidQuestionAnswerCompiled = regexp.MustCompile(stupidQuestionAnswerRegex)

func questionAnswerHandler(question, _ string) string {
	question = utils.TrimPunctuation(question)

	sub := questionAnswerCompiled.FindStringSubmatch(question)

	if len(sub) <= 1 {
		return markovGenerate(markovPickWord())
	}

	return markovGenerate(pickRandomWord(sub[1], false))
}

func stupidQuestionAnswerHandler(question, _ string) string {
	question = utils.TrimPunctuation(question)

	sub := stupidQuestionAnswerCompiled.FindStringSubmatch(question)

	if fromBasicChance("stupidQuestionHandler--ignoreUser") {
		return pickWeightedChoice(stupidAnswerResponses)
	}

	if len(sub) <= 1 {
		return markovGenerate(markovPickWord())
	}

	return markovGenerate(pickRandomWord(sub[1], false))
}
