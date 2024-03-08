package markov

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsJoinerWord(t *testing.T) {
	assert.True(t, isJoinerWord("Hello!"))
	assert.True(t, isJoinerWord("because"))
}

func TestGetWordList(t *testing.T) {
	t.SkipNow()

	connect()

	a, b := getGlobalWordEmoteLists(nil)

	assert.NotNil(t, a)
	assert.NotNil(t, b)
}

func TestGetRandomWord(t *testing.T) {
	t.SkipNow()

	connect()

	assert.NotEqual(t, "", PickWord())
}

func TestGetRandomEmoticon(t *testing.T) {
	t.SkipNow()

	connect()

	// We don't have any emoticons available in the database with the markov learning Alice in Wonderland
	assert.Equal(t, "", GetRandomEmoticon())
}

func TestMessageToWords(t *testing.T) {
	assert.Equal(
		t,
		[]string{"hello", "world!"},
		messageToWords("hello world!"),
	)
}
