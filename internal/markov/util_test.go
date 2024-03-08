package markov

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeWordList(t *testing.T) {
	words := []string{"hello", "world", "Fred", "GEORGE"}

	assert.Equal(
		t,
		[]string{"hello", "world", "Fred", "fred", "GEORGE", "george"},
		normalizeWordList(words),
	)
}
