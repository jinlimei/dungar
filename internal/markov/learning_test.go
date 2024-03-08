package markov

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLearnSentence(t *testing.T) {
	t.SkipNow()

	connect()

	tx := LearnSentence("hello world, how are you today?", false)

	assert.NotNil(t, tx)
	tx.Rollback()
}
