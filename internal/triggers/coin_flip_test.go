package triggers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestCoinFlipHandler(t *testing.T) {
	svc := initMockServices()
	random.UseTestSeed()

	msg := makeMessage("@dungar flip a coin", "fred", "butts")
	rsp := coinFlipHandler(svc, msg)
	assert.True(t, isCoinFlip(rsp[0].Contents))

	msg.Contents = "@dungar flip 40.3 coins"
	rsp = coinFlipHandler(svc, msg)

	msg.Contents = "@dungar flip 1000 coins"
	rsp = coinFlipHandler(svc, msg)

	assert.True(t, strings.Contains(rsp[0].Contents, "flipped 1000 coins:"),
		"failed to find coins: "+rsp[0].Contents)

	rgx := regexp.MustCompile("coins: (\\d+) heads, (\\d+) tails")
	m := rgx.FindAllStringSubmatch(rsp[0].Contents, 1)

	heads, _ := strconv.Atoi(m[0][1])
	tails, _ := strconv.Atoi(m[0][1])

	assert.True(t, isAround(heads, 500, 60),
		fmt.Sprintf("heads failed to be around 500 +/- 60: %d", heads))
	assert.True(t, isAround(tails, 500, 60),
		fmt.Sprintf("tails failed to be around 500 +/- 60: %d", tails))
}

func isCoinFlip(msg string) bool {
	return strings.Contains(msg, "heads") || strings.Contains(msg, "tails")
}

func TestCoinFlipHandlerOdds(t *testing.T) {
	svc := initMockServices()
	random.UseTestSeed()

	msg := makeMessage("@dungar flip a coin", "fred", "butts")

	heads := 0
	tails := 0
	sides := 0

	for i := 0; i < 1000; i++ {
		rsp := coinFlipHandler(svc, msg)

		if strings.Contains(rsp[0].Contents, "heads") {
			heads++
		} else if strings.Contains(rsp[0].Contents, "tails") {
			tails++
		} else if strings.Contains(rsp[0].Contents, "oh shit") {
			sides++
		}
	}

	assert.True(t, isAround(heads, 500, 70),
		fmt.Sprintf("failed to check if heads is around ~500 +/- 70; heads=%d", heads))
	assert.True(t, isAround(tails, 500, 70),
		fmt.Sprintf("failed to check if tails is around ~500 +/- 70; tails=%d", tails))
	assert.True(t, isAround(sides, 0, 70),
		fmt.Sprintf("failed to check if sides is around ~50 +/- 15; sides=%d", sides))
}
