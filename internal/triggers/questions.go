package triggers

import (
	"log"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var questionHandlers = []*triggerCallback{
	regexQuestion(percGameDungarRegex, percGameDungarHandler),
	regexQuestion(percGameYouRegex, percGameYouHandler),
	regexQuestion(percGameSubjectRegex, percGameSubjectHandler),
	regexQuestion(questionAnswerRegex, questionAnswerHandler),
	regexQuestion(stupidQuestionAnswerRegex, stupidQuestionAnswerHandler),
	regexQuestion(pickOptionRegex, pickOptionHandler),
	regexQuestion(whenWillRegex, whenWillHandler),
	regexQuestion(bestInLifeRegex, bestInLifeHandler),
}

var choices8Ball = []string{
	"As I see it, Yes",
	"Ask again later",
	"Cannot predict now",
	"Don't count on it",
	"It is certain",
	"It is decidedly so",
	"Most likely",
	"My reply is no",
	"My sources say no",
	"Outlook good",
	"Outlook not so good",
	"no bueno",
	"Reply hazy",
	"Signs point to........ maybe",
	"Honestly I don't know why you are asking me this. Like, am I supposed to make the big decisions now?",
	"REALLY?",
	"Very doubtful",
	"Without a doubt",
	"Yes",
	"Yes - Definitely",
	"You may rely on it",
	"As I see it, possibly. Furthermore,",
	":heysexy:",
	"Analyzing Answer....",
	"You really don't want me to answer this",
	"My sources say: you are a banana",
	"My sources say: you are a meat popsicle",
	"Outlook better than Thunderbird :xd:",
	"What do I look like, an 8ball? :colbert:",
	"heh why would you ask that? :/",
	":heavy_check_mark:",
	"When 2020 stops being 2020",
	"When 2021 stops being 2020",
	"When 2022 stops being 2020",
	"When 2023 stops being 2020",
	"When",
	"Outlook not so good, mlyp",
	"Literally no. Never. Not even ever never.",
	"Signs point to your mothers house, because that is the worlds most common destination :iceburn:",
	"When CCP fixes lag.",
	"When CCP burns eve to the ground.",
	"Cute question.",
	"idk",
	"idk lol",
	"WHY",
	"When slack stays up for more than 2 consecutive weeks",
	"When reddit does something good",
	"When bitcoin stops crashing",
	"When bitcoin starts crashing",
	"When Blizzard makes good games again",
	"When you can attach two monitors to an M2 Apple",
	"When I get verified on Twitter",
	"When Elon faces consequences for his actions",
	"When Chris Wilson buffs fun",
	"As soon as people stop mining bitcoin and switch to mining DANK MEMES",
	":thinking:",
	"When Google stops killing successful products",
	"When Google brings back Google Reader",
}

func questionsHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	contents := msg.Contents
	process := false

	// Must be directed at dungar in order to answer a question
	if isDirectedAtDungar(svc, msg) {
		process = true
	}

	if isMentioningDungar(svc, msg) && fromBasicChance("questionsHandler--mention") {
		process = true
	}

	if !process {
		return core2.EmptyRsp()
	}

	//log.Println("question: " + msg.Contents)

	for _, handler := range questionHandlers {
		if handler.matches(contents) {
			log.Println("question handling " + utils.GetFunctionName(handler.Handler))
			msg := handler.Handler(contents, msg.ServerID)

			if msg == "" {
				continue
			}

			return core2.PrefixedSingleRsp(msg)
		}
	}

	// ends with a question mark
	if contents[len(contents)-1] == '?' {
		return core2.PrefixedSingleRsp(random.PickString(choices8Ball))
	}

	if fromBasicChance("questionsHandler--markov") {
		return core2.PrefixedSingleRsp(
			markovGenerate(pickRandomWord(contents, true)),
		)
	}

	return core2.EmptyRsp()
}
