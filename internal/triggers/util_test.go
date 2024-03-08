package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixNick(t *testing.T) {
	assert.Equal(t, nil, nil)
	assert.Equal(t, 0, 0)
}

func TestPickWeightedChoice(t *testing.T) {
	choices := []weightedChoice{
		{0.01, "alfred"},
		{0.20, "fred"},
		{0.30, "bob"},
		{0.50, "george"},
	}

	results := make(map[string]int, 3)

	for i := 0; i < 10000; i++ {
		name := pickWeightedChoice(choices)

		count, there := results[name]

		if !there {
			results[name] = 1
		} else {
			results[name] = count + 1
		}
	}

	// our rng gets us within these probabilities
	assert.True(t, results["alfred"] <= 150)
	assert.True(t, results["fred"] <= 2500)
	assert.True(t, results["bob"] <= 3500)
	assert.True(t, results["george"] <= 5500)
}
