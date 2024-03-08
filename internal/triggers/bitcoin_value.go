package triggers

import (
	"fmt"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var btcLastValue = float64(-1.0)
var btcMedianValue = float64(6500.00)
var btcMaxValue = float64(25000.00)

func bitCoinValueHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !strings.HasPrefix(msg.Contents, "!bitcoin") && !strings.HasPrefix(msg.Contents, "!btc") {
		return core2.EmptyRsp()
	}

	return core2.PrefixedSingleRsp(fmt.Sprintf("Bitcoin Price: $%0.2f USD", generateBitCoinValue()))
}

func generateBitCoinValue() float64 {
	if btcLastValue <= 0.0 || btcLastValue >= btcMedianValue {
		btcLastValue = random.Float64Range(1, btcMedianValue)
	} else {
		btcLastValue = random.Float64Range(btcMedianValue, btcMaxValue)
	}

	return btcLastValue
}
