package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestPercentGameSubjectRegex(t *testing.T) {
	assert.True(t, percSubjectCompiledRegex.MatchString("@dungar how butts is fred?"))
	assert.False(t, percSubjectCompiledRegex.MatchString("how butts is fred?"))
}

func TestPercentYouSubjectRegex(t *testing.T) {
	assert.True(t, percYouCompiledRegex.MatchString("@dungar how butts am i?"))
	assert.False(t, percYouCompiledRegex.MatchString("@dungar how butts is i?"))
}

func TestPercentDungarSubjectRegex(t *testing.T) {
	assert.True(t, percDungarCompiledRegex.MatchString("@dungar how butts are you?"))
	assert.True(t, percDungarCompiledRegex.MatchString("@dungar how much butts are you?"))
	assert.False(t, percDungarCompiledRegex.MatchString("@dungar how butts is you?"))

	assert.True(t, percDungarCompiledRegex2.MatchString("@dungar how much do you like butts?"))
}

func TestPercentGameSubjectHandler(t *testing.T) {
	random.UseTestSeed()

	assert.Equal(t,
		"fred [94.52% butts]",
		percGameSubjectHandler("@dungar how butts is fred?", ""),
	)
}

func TestPercGameYouHandler(t *testing.T) {
	random.UseTestSeed()

	assert.Equal(t,
		"you're 94.52% butts.",
		percGameYouHandler("@dungar how butts am i?", ""),
	)
}

func TestPercentDungarSubjectHandler(t *testing.T) {
	random.UseTestSeed()

	retVal := percGameDungarHandler("@dungar how butts are you?", "")

	assert.True(t, strings.Contains(retVal, "I'm 65.60% butts.") || strings.Contains(retVal, "i'm 24.50% butts."))
}
