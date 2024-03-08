package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBadWords(t *testing.T) {
	// Really simple test
	TestDatabaseConnect()

	words := GetBadWords()

	assert.NotNil(t, words)
	assert.True(t, len(words) >= 0)
}
