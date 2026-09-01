package triggers

import (
	"regexp"
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

	assert.Contains(t, []string{"food", "sleep"}, pickOptionHandler("@dungar food or sleep?", ""))
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

	assert.True(t, results["a"] > 0)
	assert.True(t, results["b"] > 0)
	assert.Equal(t, 2, len(results))
}
