package triggers

import (
	"strconv"
	"strings"
	"testing"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"github.com/stretchr/testify/assert"
)

type questionResult struct {
	Result   bool
	Question string
}

func TestIsQuestionRegex(t *testing.T) {
	questions := []questionResult{
		{true, "dungar: what do you think about butts?"},
		{true, "dungar: fritos or butts"},
		{true, "dungar: hello?"},
		{true, "@dungar what do you think about butts?"},
		{true, "@dungar fritos or butts"},
		{true, "@dungar hello?"},
		{false, "what do you think about butts?"},
		{false, "fritos or butts?"},
		{false, "hello?"},
	}

	for _, qr := range questions {
		if questionStartRegex.MatchString(qr.Question) != qr.Result {
			assert.Fail(t, "Question '"+qr.Question+"' should have been "+strconv.FormatBool(qr.Result))
		}

		if qr.Result {
			matches := questionStartRegex.FindStringSubmatch(qr.Question)
			assert.Equal(t, "dungar", strings.Trim(matches[1], "@:"))
		}
	}
}

func TestQuestionsHandler(t *testing.T) {
	random.UseTestSeed()

	svc := initMockServices()
	msg := makeMessage("hello", "bob", "arena")
	rsp := questionsHandler(svc, msg)

	assert.Equal(t, core2.EmptyRsp(), rsp)

	msg.Contents = "@fred How are you?"
	rsp = questionsHandler(svc, msg)

	assert.Equal(t, core2.EmptyRsp(), rsp)

	msg.Contents = "@dungar How are you?"
	rsp = questionsHandler(svc, msg)

	assert.True(t, rsp[0].ConsumedMessage)
	assert.True(t, rsp[0].HandledMessage)

}
