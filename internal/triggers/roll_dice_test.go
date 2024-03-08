package triggers

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRollDiceHandler(t *testing.T) {
	svc := initMockServices()
	msg := makeMessage("@dungar 1d20 2d20", "kinkouin", "butts")
	rsp := rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("1d20: \\d+; 2d20: \\d+, \\d+").MatchString(rsp[0].Contents),
		"Match of 1d20 2d20 failed: "+rsp[0].Contents)

	msg.Contents = "@dungar 10d100"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("10d100: (\\d+,\\s*)+").MatchString(rsp[0].Contents),
		"Match of 10d100 failed: "+rsp[0].Contents)

	msg.Contents = "@dungar 10d10000"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("10d10000: (\\d+,\\s*)+").MatchString(rsp[0].Contents),
		"Match of 10d10000 failed: "+rsp[0].Contents)

	msg.Contents = "@dungar 100d20"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, strings.Contains(rsp[0].Contents, "100d20: nah"),
		"Match of 100d20 failed: "+rsp[0].Contents)

	msg.Contents = "@dungar 1d20 1d20 1d20 1d20 1d20 1d20 1d20"
	rsp = rollDiceHandler(svc, msg)

	assert.False(t, strings.Contains(rsp[0].Contents, "1d20"),
		"Match of '1d20 1d20 1d20 1d20 1d20 1d20 1d20' failed: "+rsp[0].Contents)
}

func TestHighLowMiddleRollDceHandler(t *testing.T) {
	svc := initMockServices()
	msg := makeMessage("@dungar 4d20h2", "kinkouin", "butts")
	rsp := rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("4d20h2: (\\d+[,\\s]*){2}$").MatchString(rsp[0].Contents),
		"Match of 4d20h2 failed: "+rsp[0].Contents)

	msg.Contents = "@dungar 6d50l3 4d20h2"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("6d50l3: (\\d+[,\\s]*){3}; 4d20h2: (\\d+[,\\s]*){2}$").MatchString(rsp[0].Contents),
		fmt.Sprintf("Match of 6d50l3 4d20h2 failed: %v\n", rsp[0].Contents))

	msg.Contents = "@dungar 2d20h2"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("2d20h2: nah").MatchString(rsp[0].Contents),
		fmt.Sprintf("Match of 2d20h2 failed: %v\n", rsp[0].Contents))

	msg.Contents = "@dungar 2d20l2"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("2d20l2: nah").MatchString(rsp[0].Contents),
		fmt.Sprintf("Match of 2d20l2 failed: %v\n", rsp[0].Contents))

	msg.Contents = "@dungar 2d20o2"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("2d20o2: nah").MatchString(rsp[0].Contents),
		fmt.Sprintf("Match of 2d20o2 failed: %v\n", rsp[0].Contents))

	msg.Contents = "@dungar 5d20o3"
	rsp = rollDiceHandler(svc, msg)

	assert.True(t, regexp.MustCompile("5d20o3: (\\d+[,\\s]*){3}$").MatchString(rsp[0].Contents),
		fmt.Sprintf("Match of 5d20o3 failed: %v\n", rsp[0].Contents))

}
