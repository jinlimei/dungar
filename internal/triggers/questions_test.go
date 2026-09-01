package triggers

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"github.com/stretchr/testify/assert"
)

type questionResult struct {
	Result   bool
	Question string
}

func useQuestionsTestMarkov(tb testing.TB) {
	tb.Helper()

	savedM3 := m3
	savedLoaded := lastLoadedM3
	savedActive := activeMarkovVersion

	tb.Cleanup(func() {
		m3 = savedM3
		lastLoadedM3 = savedLoaded
		activeMarkovVersion = savedActive
	})

	m3 = markov3.MakeMarkov("questions-test")
	m3.LearnString("hello skill friend", cleaner.VariantPlain)
	m3.LearnString("hello there friend", cleaner.VariantPlain)
	m3.LearnString("how are you friend", cleaner.VariantPlain)
	lastLoadedM3 = time.Now()
	activeMarkovVersion = mV3
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
	useQuestionsTestMarkov(t)

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

func TestNormalizeDirectedContents(t *testing.T) {
	svc := initMockServices()
	bot := core2.BotUser{
		ID:   "454311604847378442",
		Name: "Furrymatic",
	}

	assert.Equal(t, "do you skill?", normalizeDirectedContents(svc, "arena", "@Furrymatic do you skill?", bot))
	assert.Equal(t, "do you skill?", normalizeDirectedContents(svc, "arena", "Furrymatic: do you skill?", bot))
	assert.Equal(t, "do you skill?", normalizeDirectedContents(svc, "arena", "<@454311604847378442> do you skill?", bot))
	assert.Equal(t, "do you skill?", normalizeDirectedContents(svc, "arena", "<@!454311604847378442> do you skill?", bot))
}

func TestNormalizeDirectedContentsUsesGuildResolvedBotName(t *testing.T) {
	svc := initMockServices()
	mockDriver.SetBotUser(core2.BotUser{
		ID:      "dungar",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	})
	mockDriver.Users["dungar"] = core2.User{
		ID:      "dungar",
		Name:    "Furrymatic",
		IsBot:   true,
		IsAdmin: false,
	}

	assert.Equal(
		t,
		"do you skill?",
		normalizeDirectedContents(svc, "arena", "@Furrymatic do you skill?", svc.GetBotUser()),
	)
}

func TestQuestionsHandlerDirectedQuestionUsesMarkovFallback(t *testing.T) {
	random.UseTestSeed()
	useQuestionsTestMarkov(t)

	svc := initMockServices()
	mockDriver.SetBotUser(core2.BotUser{
		ID:      "454311604847378442",
		Name:    "Furrymatic",
		IsBot:   true,
		IsAdmin: false,
	})

	msg := makeMessage("@Furrymatic hello?", "bob", "arena")
	rsp := questionsHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.True(t, rsp[0].HandledMessage)
	assert.NotEmpty(t, rsp[0].Contents)
	assert.NotContains(t, choices8Ball, rsp[0].Contents)
}

func TestQuestionsHandlerDirectedTwoOptionQuestionPicksProvidedOption(t *testing.T) {
	random.UseTestSeed()

	svc := initMockServices()
	mockDriver.SetBotUser(core2.BotUser{
		ID:      "454311604847378442",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	})

	msg := makeMessage("@Dungarmatic tacos or burgers?", "bob", "arena")
	rsp := questionsHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.True(t, rsp[0].HandledMessage)
	assert.Contains(t, []string{"tacos", "burgers"}, strings.ToLower(strings.TrimSpace(rsp[0].Contents)))
}

func TestQuestionsHandlerDirectedPercentageSubjectRepliesDirectly(t *testing.T) {
	random.UseTestSeed()

	svc := initMockServices()
	mockDriver.SetBotUser(core2.BotUser{
		ID:      "454311604847378442",
		Name:    "Dungarmatic",
		IsBot:   true,
		IsAdmin: false,
	})

	msg := makeMessage("@Dungarmatic how fried is xyz?", "bob", "arena")
	rsp := questionsHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.True(t, rsp[0].HandledMessage)
	assert.False(t, rsp[0].PrefixUsername)
	assert.True(t, rsp[0].MakeSubChannelIfPossible)
	assert.Regexp(t, regexp.MustCompile(`^xyz is \d+% fried\.$`), rsp[0].Contents)
}
