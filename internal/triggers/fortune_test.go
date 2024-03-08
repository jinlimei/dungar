package triggers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestFortuneTriggerHandler(t *testing.T) {
	t.SkipNow()

	db.TestDatabaseConnect()
	random.UseTestSeed()

	rsp := fortuneHandler(makeMessage("!b", "chalur", "#butts"))
	assert.False(t, rsp[0].HandledMessage)
	assert.Equal(t, "", rsp[0].Contents)

	rsp = fortuneHandler(makeMessage("!f", "chalur", "#butts"))
	assert.True(t, rsp[0].HandledMessage)
	assert.NotEqual(t, "", rsp[0].Contents)

	rsp = fortuneHandler(makeMessage("!frisky", "haraku", "#butts"))
	assert.True(t, rsp[0].HandledMessage)
	assert.NotEqual(t, "", rsp[0].Contents)

	fortunes = []*fortune{
		{0, "Hello, World!", time.Now().Add(-2 * time.Hour), true},
		{1, "Recently Used!", time.Now().Add(-30 * time.Minute), true},
		{2, "Never Used", time.Time{}, false},
	}

	rsp = fortuneHandler(makeMessage("!fritos", "kryptn", "#butts"))
	assert.True(t, rsp[0].HandledMessage)
	assert.NotEqual(t, "", rsp[0].Contents)

	fortunes = []*fortune{
		{0, "aaaaaaa", time.Now(), true},
	}

	rsp = fortuneHandler(makeMessage("!fritos", "kryptn", "#butts"))
	assert.True(t, rsp[0].HandledMessage)
	assert.Contains(t, rsp[0].Contents, "fortune teller")

}
