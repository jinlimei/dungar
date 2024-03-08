package triggers

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"github.com/stretchr/testify/assert"
)

func TestChanceGameHandler(t *testing.T) {
	svc := initMockServices()
	random.UseTestSeed()

	validRsps := []string{
		"c",
		"d",
		"v0v",
		"idk",
	}

	msg := makeMessage("hello", "george", "butts")
	rsp := chanceGameHandler(svc, msg)
	assert.Equal(t, core2.EmptyRsp(), rsp)

	msg.Contents = "@dungar c/d"
	rsp = chanceGameHandler(svc, msg)
	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)

	assert.True(t, strContainsAny(rsp[0].Contents, validRsps),
		fmt.Sprintf("Rsp: %s", rsp[0].Contents))

	msg.Contents = "get food c/d"
	rsp = chanceGameHandler(svc, msg)
	assert.Len(t, rsp, 1)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)

	assert.True(t, strContainsAny(rsp[0].Contents, validRsps),
		fmt.Sprintf("Rsp: %s", rsp[0].Contents))

	hasEmpty := false
	for k := 0; k < 50; k++ {
		rsp = chanceGameHandler(svc, msg)

		if rsp[0].IsEmpty() {
			hasEmpty = true
			break
		}
	}

	assert.True(t, hasEmpty)
}

func strContainsAny(str string, any []string) bool {
	for _, s := range any {
		if strings.Contains(str, s) {
			return true
		}
	}

	return false
}

func TestChanceTable(t *testing.T) {
	result := make(map[string]int, 0)

	choice := ""
	for i := 0; i < 10000; i++ {
		choice = pickWeightedChoice(chanceTable)

		resp, ok := result[choice]

		if ok {
			result[choice] = resp + 1
		} else {
			result[choice] = 1
		}
	}

	assert.True(t, result["v0v"] < 500)
	assert.True(t, result["idk"] < 500)
	assert.True(t, result["d"] >= 4000)
	assert.True(t, result["c"] >= 4000)
}

func TestChanceRegex(t *testing.T) {
	choices := []string{
		"butts c/d",
		"things c/d?",
		"idk c/d?",
		"fritos c/d",
	}

	for _, choice := range choices {
		if !chanceRegex.MatchString(choice) {
			assert.Fail(t, "string '"+choice+"' should have passed")
		}
	}
}
