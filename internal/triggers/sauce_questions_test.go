package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/dungar/internal/random"
)

func TestSauceQuestionsHandler(t *testing.T) {
	initMockServices()
	random.UseTestSeed()

	msg := makeMessage("hey bitcoin", "fred", "butts")

	masterChanceList["sauceQuestionsHandler--consume"] = 0.60

	handled := false
	for range 30 {
		out := sauceQuestionsHandler(msg)
		if len(out) > 0 && out[0].HandledMessage {
			handled = true
			break
		}
	}

	assert.True(t, handled)
}
