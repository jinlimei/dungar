package triggers

import "gitlab.int.magneato.site/dungar/prototype/internal/random"

const whenWillRegex = ":?\\s*(when|how long)\\b.+\\??"

// 1 hour from now
const minWhenWill = 1 * 60 * 60

// 6 years into the future
const maxWhenWill = 6 * 365 * 86400

var whenWillChoices = []weightedChoice{
	{0.42, "At 4:20 :420::350::420::350:"},
	{0.20, "Never"},
	{0.20, "Always"},
	{0.08, "Right Now"},
}

func whenWillHandler(_, _ string) string {
	if fromBasicChance("whenWillHandler--alternatives") {
		return pickWeightedChoice(whenWillChoices)
	}

	return random.MakeTime(minWhenWill, maxWhenWill)
}
