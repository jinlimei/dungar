package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestNegativeMentionSnarkHandler(t *testing.T) {
	random.UseTestSeed()

	savedChance := masterChanceList["negativeMentionSnarkHandler--respond"]
	t.Cleanup(func() {
		masterChanceList["negativeMentionSnarkHandler--respond"] = savedChance
	})

	svc := initMockServices()
	masterChanceList["negativeMentionSnarkHandler--respond"] = 1.0

	msg := makeMessage("@Dungar you suck", "bob", "butts")
	msg.ServerID = "guild"

	rsp := negativeMentionSnarkHandler(svc, msg)

	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)
	assert.NotEmpty(t, rsp[0].Contents)
}
