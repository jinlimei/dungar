package triggers

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestPickOptionRegex(t *testing.T) {
	regex, err := regexp.Compile(pickOptionRegex)

	assert.Nil(t, err)
	assert.NotNil(t, regex)

	matches := []string{
		"fritos or butts",
		"fish or fred",
		"sticks, cults, fishes, or freds",
		"0,3,4, or 5",
		"0 or 3 or 4 or 5",
	}

	for _, match := range matches {
		if !regex.MatchString(match) {
			assert.Fail(t, "did not match '"+match+"'")
		}
	}
}

func TestPickOption2(t *testing.T) {
	random.UseTestSeed()

	quickTest := func(possibilities []string) {
		result := pickOptionHandler("@dungar food or sleep?", "")

		passed := false
		for _, possible := range possibilities {
			if result != possible {
				passed = true
				break
			}
		}

		if !passed {
			assert.Fail(t,
				fmt.Sprintf("failed to find '%s' in '%v'", result, possibilities))
		}
	}

	quickTest([]string{"checkbox, voted all", "fred"})

	assert.Equal(t, "food", pickOptionHandler("@dungar food or sleep?", ""))
	assert.Equal(t, "sleep", pickOptionHandler("@dungar food or sleep?", ""))
	assert.Equal(t, "food", pickOptionHandler("@dungar food or sleep?", ""))
}

func TestPickOption3(t *testing.T) {
	splits := splitterRegexp.Split("thing1, thing2, order thing three?", -1)
	assert.Len(t, splits, 3)

	splits = splitterRegexp.Split("a,b,c or d?", -1)
	assert.Len(t, splits, 4)

	splits = splitterRegexp.Split("0,3,4, or 5?", -1)
	assert.Len(t, splits, 4)

	splits = splitterRegexp.Split("0 or 3 or 4 or 5", -1)
	assert.Len(t, splits, 4)
}

func TestPickOption4(t *testing.T) {
	result := pickOptionHandler("<@U9LDWA6QL|dungar> mizuho, kamoi, hayasui, comma, or gotland?", "")

	assert.NotEqual(t, "", result)

	result = pickOptionHandler("<@U9LDWA6QL|dungar> 0,3,4, or 5?", "")

	assert.NotEqual(t, "", result)

	result = pickOptionHandler("<@U9LDWA6QL|dungar> 0 or 3 or 4 or 5?", "")

	assert.NotEqual(t, "", result)
}

func TestPickOption(t *testing.T) {
	results := make(map[string]int, 0)

	for i := 0; i < 50000; i++ {
		opt := pickOptionHandler("a or b?", "")

		count, ok := results[opt]

		if ok {
			results[opt] = count + 1
		} else {
			results[opt] = 1
		}
	}

	assert.True(t, results["a"] > 0 && results["a"] <= 25000)
	assert.True(t, results["b"] > 0 && results["b"] <= 25000)

	nickDeciders := 0

	for key, val := range results {
		if strings.Contains(key, "let") && strings.Contains(key, "decide") {
			nickDeciders += val
		}
	}

	assert.True(t, nickDeciders > 0 && nickDeciders <= 1000)
}
