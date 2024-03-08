package triggers

var masterChanceList = map[string]float64{
	"alexJonesHandler--prepositions":    0.25,
	"chanceGameHandler--basic":          0.90,
	"coinFlipHandler--singleFlip":       0.50,
	"coinFlipHandler--flip":             0.50,
	"coinFlipHandler--singleSide":       0.0001,
	"commentateScheduler--skip":         0.60,
	"commentateHandler--commentate":     0.009,
	"dadJokeHandler--markov":            0.01,
	"markovHandler--badTouch":           0.10,
	"noYouHandler--regexp":              0.0005,
	"noYouHandler--normal":              0.0005,
	"percGameHandler--8ball":            0.05,
	"pickOptionHandler--ignoreUser":     0.01,
	"pickOptionHandler--altOptions":     0.20,
	"stupidQuestionHandler--ignoreUser": 0.10,
	"questionsHandler--markov":          0.60,
	"questionsHandler--mention":         0.50,
	"repeatableHandler--literally":      0.01,
	"repeatableHandler--fingerGuns":     0.01,
	"repeatableHandler--basic":          0.01,
	"sauceQuestionsHandler--consume":    0.05,
	"sequenceHandler--respond":          0.05,
	"whenWillHandler--alternatives":     0.30,
	"simplePrefixHandler--basic":        0.40,
	"spamHandler--respond":              0.30,
}

func fromBasicChance(name string) bool {
	return basic.Execute(name, masterChanceList[name])
}

func fromDefinedChance(name string, chance float64) bool {
	return basic.Execute(name, chance)
}
