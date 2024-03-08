package triggers

var bestInLife = []weightedChoice{
	{0.1, "To crush your biscuits, see them dunked in your tea, and hear the conversations of your auntie"},
	{0.2, "To post your memes, see them reposted in other subreddits, and feel the joys of updoots across the redditverse"},
	{0.3, "To mine your bitcoins, see them stolen in an exchange, and hodling for a brighter day"},
	{1.0, "To crush your enemies, see them driven before you, and to hear the lamentation of their women"},
}

const bestInLifeRegex = "(the )?best in life"

func bestInLifeHandler(_, _ string) string {
	return pickWeightedChoice(bestInLife)
}
