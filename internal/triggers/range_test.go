package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func TestRangeHandler(t *testing.T) {
	random.UseTestSeed()

	slicedResponses := sliceFromWeightedChoice(badMinMaxResponses)

	svc := initMockServices()
	msg := makeMessage("@dungar pick a number between 1-4,910.85", "butts", "butts")
	out := rangeHandler(svc, msg)
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "number"))
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "."))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between 5 and 10", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "butt"))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between 0 and 3", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "butt"))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between 0 and 0.3", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "butt"))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between -5 and 0", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "butt"))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between -5 and 5", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, strings.Contains(out[0].Format(mockDriver, msg), "butt"))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between 0 and 0", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, utils.StringInSlice(out[0].Contents, slicedResponses))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between -0 and 0", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, utils.StringInSlice(out[0].Contents, slicedResponses))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between -0 and -0", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, utils.StringInSlice(out[0].Contents, slicedResponses))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between 0.0 and -0", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, utils.StringInSlice(out[0].Contents, slicedResponses))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between -0.0 and 0", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, utils.StringInSlice(out[0].Contents, slicedResponses))
	assert.True(t, out[0].HandledMessage)

	msg = makeMessage("@dungar pick a butt between 10 and -10", "butts", "butts")
	out = rangeHandler(svc, msg)
	assert.True(t, utils.StringInSlice(out[0].Contents, slicedResponses))
	assert.True(t, out[0].HandledMessage)
}
