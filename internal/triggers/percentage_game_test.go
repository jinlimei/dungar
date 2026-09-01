package triggers

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestPercentGameSubjectRegex(t *testing.T) {
	assert.True(t, percSubjectCompiledRegex.MatchString("@dungar how butts is fred?"))
	assert.True(t, percSubjectCompiledRegex.MatchString("how butts is fred?"))
}

func TestPercentYouSubjectRegex(t *testing.T) {
	assert.True(t, percYouCompiledRegex.MatchString("@dungar how butts am i?"))
	assert.True(t, percYouCompiledRegex.MatchString("how butts am i?"))
	assert.False(t, percYouCompiledRegex.MatchString("@dungar how butts is i?"))
}

func TestPercentDungarSubjectRegex(t *testing.T) {
	assert.True(t, percDungarCompiledRegex.MatchString("@dungar how butts are you?"))
	assert.True(t, percDungarCompiledRegex.MatchString("@dungar how much butts are you?"))
	assert.True(t, percDungarCompiledRegex.MatchString("how butts are you?"))
	assert.False(t, percDungarCompiledRegex.MatchString("@dungar how butts is you?"))

	assert.True(t, percDungarCompiledRegex2.MatchString("@dungar how much do you like butts?"))
	assert.True(t, percDungarCompiledRegex2.MatchString("how much do you like butts?"))
}

func TestPercentGameSubjectHandler(t *testing.T) {
	random.UseTestSeed()

	assert.Regexp(t,
		regexp.MustCompile(`^fred is \d+% butts\.$`),
		percGameSubjectHandler("how butts is fred?", ""),
	)
}

func TestPercGameYouHandler(t *testing.T) {
	random.UseTestSeed()

	assert.Regexp(t,
		regexp.MustCompile(`^you're \d+% butts\.$`),
		percGameYouHandler("how butts am i?", ""),
	)
}

func TestPercentDungarSubjectHandler(t *testing.T) {
	random.UseTestSeed()

	retVal := percGameDungarHandler("how butts are you?", "")

	assert.True(t, strings.Contains(retVal, "I'm ") || strings.Contains(retVal, "i'm ") || strings.Contains(retVal, "im ") || strings.Contains(retVal, "am "))
	assert.True(t, strings.HasSuffix(retVal, "% butts."))
}
