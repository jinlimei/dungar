package triggers

import (
	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"strings"
	"testing"
)

func TestIsRepeatableMessage(t *testing.T) {
	random.UseTestSeed()

	assert.False(t, isRepeatableMessage(""))
	assert.False(t, isRepeatableMessage("aa"))
	assert.False(t, isRepeatableMessage("hello world"))
	assert.False(t, isRepeatableMessage("https://www"))
	assert.False(t, isRepeatableMessage("<thing>"))
	assert.False(t, isRepeatableMessage("<@USER>"))
	assert.False(t, isRepeatableMessage("<#CHAN|NAME>"))
	assert.False(t, isRepeatableMessage(":hello world:"))

	assert.True(t, isRepeatableMessage(":thinking:"))
	// disabled these
	//assert.True(t, isRepeatableMessage("hello"))
	//assert.True(t, isRepeatableMessage("ohgod"))
	//assert.True(t, isRepeatableMessage("ohgod"))
}

func TestHandleRepeatable(t *testing.T) {
	msg := makeMessage(":thinking:", "jinli", "#butts")

	var (
		hasLiterally  = false
		hasFingerGuns = false
		hasNormal     = false
	)

	masterChanceList["repeatableHandler--literally"] = 0.80
	masterChanceList["repeatableHandler--fingerGuns"] = 0.80
	masterChanceList["repeatableHandler--basic"] = 0.80

	for k := 0; k < 50; k++ {
		rsp := repeatableHandler(msg)
		if len(rsp) == 0 || rsp[0].IsEmpty() {
			continue
		}

		val := rsp[0].Contents

		if strings.Contains(val, "literally") {
			hasLiterally = true
		} else if strings.Contains(val, ":point_right:") {
			hasFingerGuns = true
		} else if val == ":thinking:" {
			hasNormal = true
		}
	}

	assert.True(t, hasLiterally, "Could not find a 'literally <x>' instance")
	assert.True(t, hasFingerGuns, "Could not find a ':point_right: <x>' instance")
	assert.True(t, hasNormal, "Could not find a normal instance")

}
