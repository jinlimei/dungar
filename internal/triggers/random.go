package triggers

import (
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	basic = random.NewBasicControlGroup()
	adv   = random.NewAdvControlGroup()

	lastRanAnalytics time.Time
)

func randomHandler(msg *core2.IncomingMessage) []*core2.Response {
	now := time.Now()

	if lastRanAnalytics.Add(5 * time.Minute).Before(now) {
		return core2.EmptyRsp()
	}

	if strings.HasPrefix(msg.Lowered(), "!analytics") {
		basicAnalytics := basic.Analytics()

		rsp := ""
		for _, data := range basicAnalytics {
			rsp += data.String() + "\n"
		}

		return core2.MakeSingleRsp(rsp)
	}

	return core2.EmptyRsp()
}
