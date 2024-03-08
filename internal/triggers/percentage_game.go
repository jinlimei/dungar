package triggers

import (
	"fmt"
	"regexp"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

const (
	percGameDungarRegex2 = "^(@[^ ]+:?|[^ ]+:) [Hh]ow much do you ([\\w ]+)\\??$"
	percGameDungarRegex  = "^(@[^ ]+:?|[^ ]+:) [Hh]ow (?:much |)([\\w ]+) are you\\??$"
	percGameSubjectRegex = "^(@[^ ]+:?|[^ ]+:) [Hh]ow ([\\w ]+) (?:is|are) ([^?]+)\\??$"
	percGameYouRegex     = "^(@[^ ]+:?|[^ ]+:) [Hh]ow ([\\w ]+) am [Ii]\\??$"
)

var (
	percDungarCompiledRegex2 = regexp.MustCompile(percGameDungarRegex2)
	percDungarCompiledRegex  = regexp.MustCompile(percGameDungarRegex)
	percSubjectCompiledRegex = regexp.MustCompile(percGameSubjectRegex)
	percYouCompiledRegex     = regexp.MustCompile(percGameYouRegex)
)

var percSubjectChoices = []weightedChoice{
	{0.70, "I'm"},
	{0.14, "i'm"},
	{0.14, "im"},
	{0.02, "am"},
}

func percGameDungarHandler(msg, _ string) string {
	var matches []string

	hasMatch := false

	if percDungarCompiledRegex.MatchString(msg) {
		hasMatch = true
		matches = percDungarCompiledRegex.FindStringSubmatch(msg)
	} else if percDungarCompiledRegex2.MatchString(msg) {
		hasMatch = true
		matches = percDungarCompiledRegex2.FindStringSubmatch(msg)
	}

	if hasMatch && len(matches) >= 3 {
		if fromBasicChance("percGameHandler--8ball") {
			return random.PickString(choices8Ball)
		}

		return fmt.Sprintf("%s %0.2f%% %s.", pickWeightedChoice(percSubjectChoices), random.Float64()*100, matches[2])
	}

	return ""
}

func percGameYouHandler(msg, _ string) string {
	if !percYouCompiledRegex.MatchString(msg) {
		return ""
	}

	matches := percYouCompiledRegex.FindStringSubmatch(msg)

	return fmt.Sprintf("you're %0.2f%% %s.", random.Float64()*100, matches[2])
}

func percGameSubjectHandler(msg, _ string) string {
	if !percSubjectCompiledRegex.MatchString(msg) {
		return ""
	}

	matches := percSubjectCompiledRegex.FindStringSubmatch(msg)

	return fmt.Sprintf("%s [%0.2f%% %s]", matches[3], random.Float64()*100, matches[2])
}
