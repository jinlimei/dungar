package triggers

import (
	"fmt"
	"regexp"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

const (
	percGamePrefixRegex  = "(?:(?:@[^ ]+:?|[^ ]+:)\\s+)?"
	percGameDungarRegex2 = "^" + percGamePrefixRegex + "[Hh]ow much do you ([\\w ]+)\\??$"
	percGameDungarRegex  = "^" + percGamePrefixRegex + "[Hh]ow (?:much |)([\\w ]+) are you\\??$"
	percGameSubjectRegex = "^" + percGamePrefixRegex + "[Hh]ow ([\\w ]+) (?:is|are) ([^?]+)\\??$"
	percGameYouRegex     = "^" + percGamePrefixRegex + "[Hh]ow ([\\w ]+) am [Ii]\\??$"
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

	if hasMatch && len(matches) >= 2 {
		if fromBasicChance("percGameHandler--8ball") {
			return random.PickString(choices8Ball)
		}

		return fmt.Sprintf("%s %d%% %s.", pickWeightedChoice(percSubjectChoices), random.Int(101), strings.TrimSpace(matches[1]))
	}

	return ""
}

func percGameYouHandler(msg, _ string) string {
	if !percYouCompiledRegex.MatchString(msg) {
		return ""
	}

	matches := percYouCompiledRegex.FindStringSubmatch(msg)

	return fmt.Sprintf("you're %d%% %s.", random.Int(101), strings.TrimSpace(matches[1]))
}

func percGameSubjectHandler(msg, _ string) string {
	if !percSubjectCompiledRegex.MatchString(msg) {
		return ""
	}

	matches := percSubjectCompiledRegex.FindStringSubmatch(msg)

	return fmt.Sprintf("%s is %d%% %s.", strings.TrimSpace(matches[2]), random.Int(101), strings.TrimSpace(matches[1]))
}
