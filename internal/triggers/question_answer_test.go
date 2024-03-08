package triggers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

var testQuestionMatches = []string{
	"do you hate pizza",
	"what do you think about pizza?",
	"how do you like bananas?",
	"do you wonder about bananas?",
}

func TestQuestionAnswerRegex(t *testing.T) {
	for _, match := range testQuestionMatches {
		if !questionAnswerCompiled.MatchString(match) {
			assert.Fail(t, "did not match '"+match+"'")
		}
	}

	assert.False(t, questionAnswerCompiled.MatchString("do you x?"))
	assert.False(t, questionAnswerCompiled.MatchString("do you batman.jpg?"))
	// last word match
	assert.Equal(t, "pizza?", questionAnswerCompiled.FindStringSubmatch("do you hate pizza?")[1])
}

func TestStupidQuestionAnswerRegex(t *testing.T) {
	for _, match := range testQuestionMatches {
		if !stupidQuestionAnswerCompiled.MatchString(match) {
			assert.Fail(t, "did not match '"+match+"'")
		}
	}

	assert.True(t, stupidQuestionAnswerCompiled.MatchString("do you batman.jpg?"))
	assert.Equal(t, []string{"do you batman.jpg?", "batman.jpg"}, stupidQuestionAnswerCompiled.FindStringSubmatch("do you batman.jpg?"))
}

func TestQuestionAnswer(t *testing.T) {
	useAliceInWonderland()

	//db.TestDatabaseConnect()

	random.UseTestSeed()

	assert.NotNil(t, questionAnswerHandler("what do you think about ))?", ""))

	res := strings.ToLower(questionAnswerHandler("what do you think about alice wonderland?", ""))

	assert.True(t, strings.Contains(res, "alice") || strings.Contains(res, "wonderland"),
		"could not find alice or wonderland in the string")
	assert.True(t, strings.Contains(res, "alice") || strings.Contains(res, "wonderland"),
		"could not find alice or wonderland in the string")

	askQuestion := func() string {
		return strings.ToLower(
			questionAnswerHandler("What do you think about alice?", ""),
		)
	}

	assert.Contains(t, askQuestion(), "alice")
	assert.Contains(t, askQuestion(), "alice")
	assert.Contains(t, askQuestion(), "alice")
}

func TestWeirdQnAMarkov(t *testing.T) {
	//db.TestDatabaseConnect()
	useAliceInWonderland()

	seeds := [4]int64{0, 1337, 31337, 42069}

	for _, seed := range seeds {
		random.UseSeed(seed)

		val := questionAnswerHandler("what do you think about duchess?", "")
		assert.True(t, strings.Contains(val, "Duchess") || strings.Contains(val, "duchess"),
			fmt.Sprintf("Failed: answer='%v'\n", val))
	}
}

func TestStupidQuestionAnswerHandler(t *testing.T) {
	//db.TestDatabaseConnect()
	useAliceInWonderland()

	assert.NotNil(t, stupidQuestionAnswerHandler("do you batman.jpg?", ""))
}
