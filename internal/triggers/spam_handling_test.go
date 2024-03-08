package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestUserSpammedHandler(t *testing.T) {
	assert.Equal(t, 1, 1)

	svc := initMockServices()

	bot := svc.GetBotUser()
	msg := &core2.IncomingMessage{
		UserID: bot.ID,
	}

	assert.Equal(t, core2.EmptyRsp(), userSpammedHandler(svc, msg))

	msg.UserID = "FRED"

	assert.Equal(t, core2.EmptyRsp(), userSpammedHandler(svc, msg))

	msg.Contents = strings.Repeat("A", 1024)

	rsp := userSpammedHandler(svc, msg)
	assert.NotEqual(t, core2.EmptyRsp(), rsp)
	assert.NotNil(t, rsp)
	assert.False(t, rsp[0].ConsumedMessage)
	assert.True(t, rsp[0].HandledMessage)
}

func TestSelfSpamHandler(t *testing.T) {
	svc := initMockServices()
	inp := core2.MakeSingleRsp("AAAAAAAAAAA")

	rsp := selfSpamHandler(svc, inp)
	assert.Len(t, rsp, 0)

	inp[0].Contents = strings.Repeat("A", 1024)

	rsp = selfSpamHandler(svc, inp)
	assert.Len(t, rsp, 1)
}
