package triggers

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestWhenWillRegex(t *testing.T) {
	reg := regexp.MustCompile(whenWillRegex)

	matches := []string{
		"when will the earth explode?",
		"when will mars collide",
		"when will woet be nice",
		"how long until the earth explodes?",
	}

	for _, v := range matches {
		if !reg.MatchString(v) {
			assert.Fail(t, "Regex failed to match '"+v+"'")
		}
	}
}

func TestWhenWillHandler(t *testing.T) {
	random.UseTestSeed()

	has420 := false
	hasNever := false
	hasAlways := false
	hasNow := false

	for i := 0; i < 150; i++ {
		val := whenWillHandler("", "")
		assert.NotEmpty(t, val)

		has420 = has420 || strings.Contains(val, ":420:")
		hasNever = hasNever || strings.Contains(val, "Never")
		hasAlways = hasAlways || strings.Contains(val, "Always")
		hasNow = hasNow || strings.Contains(val, "Now")
	}

	assert.True(t, has420)
	assert.True(t, hasNever)
	assert.True(t, hasAlways)
	assert.True(t, hasNow)
}
