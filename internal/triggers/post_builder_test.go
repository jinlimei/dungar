package triggers

import (
	"fmt"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestBuildPostMessage(t *testing.T) {
	random.UseTestSeed()

	assert.Equal(t, "", buildPostMessage("", ""))

	assert.Equal(t, "soylent fork", buildPostMessage("$randomFood$ $randomNoun$", ""))

	valids := []string{
		"ethereum bitcoin romeo",
		"overwatch pome fork",
	}

	built := buildPostMessage("$randomNoun$ $randomNoun$ $randomNoun$", "")

	assert.True(
		t,
		utils.StringInSlice(built, valids),
		fmt.Sprintf("Built '%v' not in list of valid strings\n", built),
	)
}
