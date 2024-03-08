package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestBitCoinHandler(t *testing.T) {
	random.UseTestSeed()

	assert.Equal(
		t,
		"Introducing OverwatchCoin: A new, modern, sexy cryptocurrency to disrupt the fork industry.",
		bitCoinHandler("", ""),
	)
}

func TestJavaScriptHandler(t *testing.T) {
	random.UseTestSeed()

	assert.Equal(
		t,
		"Check out Overwatch.js: A new framework that is going to revolutionize working with fork on Docker",
		jsHandler("", ""),
	)
}
