package triggers

import (
	"regexp"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

type weightedChoice struct {
	Chance   float64
	Response string
}

type triggerChoice struct {
	Trigger  string
	IsRegex  bool
	Response string

	cache *regexp.Regexp
}

type triggerCallback struct {
	Trigger string
	IsRegex bool
	Handler func(str, serverID string) string

	cache *regexp.Regexp
}

type multiTriggerCallback struct {
	Triggers []string
	IsRegex  bool
	Handler  func(str, serverID string) string

	cache []*regexp.Regexp
}

func (tc *triggerCallback) matches(str string) bool {
	str = strings.ToLower(str)

	if !tc.IsRegex {
		return strings.Contains(str, " "+tc.Trigger+" ") ||
			strings.Contains(str, " "+tc.Trigger) ||
			strings.Contains(str, tc.Trigger+" ")
	}

	if tc.cache == nil {
		tc.cache = regexp.MustCompile(tc.Trigger)
	}

	return tc.cache.MatchString(str)
}

func (mtc *multiTriggerCallback) matches(str string) bool {
	str = strings.ToLower(str)

	if !mtc.IsRegex {
		for _, trigger := range mtc.Triggers {
			if strings.Contains(str, " "+trigger+" ") ||
				strings.Contains(str, " "+trigger) ||
				strings.Contains(str, trigger+" ") {

				return true
			}
		}

		return false
	}

	ensureMtcCache(mtc)
	for _, re := range mtc.cache {
		if re.MatchString(str) {
			return true
		}
	}

	return false
}

func ensureMtcCache(mtc *multiTriggerCallback) {
	if mtc.cache == nil || len(mtc.cache) == 0 {
		mtc.cache = make([]*regexp.Regexp, len(mtc.Triggers))

		for i, k := range mtc.Triggers {
			mtc.cache[i] = regexp.MustCompile(k)
		}
	}
}

func (tc triggerChoice) matches(str string) bool {
	if tc.IsRegex {
		if tc.cache == nil {
			tc.cache = regexp.MustCompile(tc.Trigger)
		}

		return tc.cache.MatchString(str)
	}

	return strings.Contains(str, tc.Trigger)
}

func randomWeightChoice(choice string) weightedChoice {
	return weightedChoice{
		random.Float64(),
		choice,
	}
}

func (wc weightedChoice) randomizeChance() {
	wc.Chance = random.Float64()
}

func sliceFromWeightedChoice(choices []weightedChoice) []string {
	output := make([]string, 0)

	for _, value := range choices {
		output = append(output, value.Response)
	}

	return output
}

// Bmt = BuildMultiTrigger (I know)
func buildMultiQuestions(callback func(trigger, serverID string) string, questions ...string) multiTriggerCallback {
	return multiTriggerCallback{
		Triggers: questions,
		IsRegex:  false,
		Handler:  callback,
	}
}

func regexQuestion(regex string, handler func(trigger, serverID string) string) *triggerCallback {
	return &triggerCallback{
		Trigger: regex,
		IsRegex: true,
		Handler: handler,
	}
}

func buildMultiRegexQuestions(callback func(trigger, serverID string) string, questions ...string) multiTriggerCallback {
	return multiTriggerCallback{
		Triggers: questions,
		IsRegex:  true,
		Handler:  callback,
	}
}
