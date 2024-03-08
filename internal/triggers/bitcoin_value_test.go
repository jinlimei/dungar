package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestBitCoinValueHandler(t *testing.T) {
	random.UseTestSeed()

	msg := makeMessage("abcdefg", "fred", "butts")

	rsp := bitCoinValueHandler(msg)

	assert.Equal(t, core2.EmptyRsp(), rsp)

	msg.Contents = "!bitcoin"
	rsp = bitCoinValueHandler(msg)
	assert.Contains(t, rsp[0].Contents, "Bitcoin")

	msg.Contents = "!btc"
	rsp = bitCoinValueHandler(msg)
	assert.Contains(t, rsp[0].Contents, "Bitcoin")
}

func TestGenerateBitCoinValue(t *testing.T) {
	random.UseTestSeed()

	val := generateBitCoinValue()
	assert.True(t, val > 0)
	assert.True(t, val < btcMedianValue)

	val = generateBitCoinValue()
	assert.True(t, val > btcMedianValue)
	assert.True(t, val < btcMaxValue)

	val = generateBitCoinValue()
	assert.True(t, val > 0)
	assert.True(t, val < btcMedianValue)
}
