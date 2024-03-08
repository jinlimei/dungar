package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestSauceQuestionsHandler(t *testing.T) {
	initMockServices()
	random.UseTestSeed()

	msg := makeMessage("hey bitcoin", "fred", "butts")

	masterChanceList["sauceQuestionsHandler--consume"] = 0.60

	handled := false
	for k := 0; k < 30; k++ {
		out := sauceQuestionsHandler(msg)
		if len(out) > 0 && out[0].HandledMessage {
			handled = true
			break
		}
	}

	assert.True(t, handled)
}
