package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestMainControlHandler(t *testing.T) {
	msg := &core2.IncomingMessage{
		UserID:        admins[0],
		Contents:      "!change bananas 0",
		LowerContents: "!change bananas 0",
	}

	rsp := mainControlHandler(msg)
	assert.NotNil(t, rsp)
	assert.True(t, len(rsp) > 0)

	assert.True(t, strings.Contains(rsp[0].Contents, "not find value"))

	msg.Contents = "!change alexJonesHandler--prepositions 0.54"
	msg.LowerContents = "!change alexjoneshandler--prepositions 0.54"

	rsp = mainControlHandler(msg)
	assert.NotNil(t, rsp)
	assert.True(t, len(rsp) > 0)

	assert.Equal(t, "Changed val 'alexJonesHandler--prepositions' to be from '0.2500' to '0.5400'",
		rsp[0].Contents)

	msg.Contents = "!change alexJonesHandler--prepositions banana"
	msg.LowerContents = "!change alexjoneshandler--prepositions banana"

	rsp = mainControlHandler(msg)
	assert.NotNil(t, rsp)
	assert.True(t, len(rsp) > 0)

	assert.Equal(t,
		"Could not figure out next value 'banana': strconv.ParseFloat: parsing \"banana\": invalid syntax",
		rsp[0].Contents,
	)

	msg.Contents = "!disable alexJonesHandler--prepositions"
	msg.LowerContents = "!disable alexjoneshandler--prepositions"

	rsp = mainControlHandler(msg)
	assert.NotNil(t, rsp)
	assert.True(t, len(rsp) > 0)

	assert.Equal(t,
		"Changed val 'alexJonesHandler--prepositions' to be from '0.5400' to '0.0000'",
		rsp[0].Contents,
	)
}
