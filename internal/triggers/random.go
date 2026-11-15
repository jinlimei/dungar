package triggers

import (
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/dungar/internal/random"
	"gitlab.int.magneato.site/dungar/dungar/library/core2"
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

		var rsp strings.Builder
		for _, data := range basicAnalytics {
			rsp.WriteString(data.String() + "\n")
		}

		return core2.MakeSingleRsp(rsp.String())
	}

	return core2.EmptyRsp()
}
