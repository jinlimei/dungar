package triggers

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gopkg.in/ini.v1"
)

const channelsAvailableJSON = `{"teamId":"T031VPPUH","results":[{"teamId":"T031VPPUH","channelId":"C019C00F0M9","channelName":"factoreo"},{"teamId":"T031VPPUH","channelId":"C01A5B9CF99","channelName":"buttofthemonth"},{"teamId":"T031VPPUH","channelId":"C026E2T9F2A","channelName":"buttmobiles"},{"teamId":"T031VPPUH","channelId":"C031VPPV3","channelName":"butts"},{"teamId":"T031VPPUH","channelId":"C3R0V1AU8","channelName":"dungar-test"},{"teamId":"T031VPPUH","channelId":"C5FQJ5ZQA","channelName":"nsfw"},{"teamId":"T031VPPUH","channelId":"C6B54067K","channelName":"arpg"},{"teamId":"T031VPPUH","channelId":"C6S7J9CAW","channelName":"vegas-dreams"},{"teamId":"T031VPPUH","channelId":"C7ZPDPR8W","channelName":"programming"},{"teamId":"T031VPPUH","channelId":"CBG3YS874","channelName":"expanseme"},{"teamId":"T031VPPUH","channelId":"CCQM6G5HR","channelName":"politics"},{"teamId":"T031VPPUH","channelId":"CMF9WUF3P","channelName":"homeo"},{"teamId":"T031VPPUH","channelId":"CQFV8JKME","channelName":"politics"},{"teamId":"T031VPPUH","channelId":"G6SBC08NR","channelName":"romeo-vegas"},{"teamId":"T031VPPUH","channelId":"GBASAJ8QP","channelName":"romeo-haus"}],"totalResults":15}`

var testPinChannels []pinChannel

func ensureTestPinsSetup() {
	if activePinConf == nil {
		if utils.MustUseEnvVars() {
			creds := utils.PinCredentials()
			activePinConf = &pinConf{
				Auth:   creds["auth"],
				TeamID: creds["team"],
				URL:    creds["url"],
			}
		} else {
			cfg, err := ini.Load("../../secrets.ini")
			utils.HaltingError("requestAllPins ini", err)

			activePinConf = &pinConf{
				URL:    cfg.Section("pins").Key("url").String(),
				Auth:   cfg.Section("pins").Key("auth").String(),
				TeamID: cfg.Section("slack").Key("team_id").String(),
			}
		}
	}

	var rsp *pinChannelsResponse

	json.Unmarshal([]byte(channelsAvailableJSON), &rsp)

	channelsAvailable = &pinChannelsAvailable{
		LastRequest: time.Now().Unix(),
		Channels:    rsp.Results,
	}

	testPinChannels = rsp.Results
}

func oldTestChannelsUsed(t *testing.T) {
	random.UseTestSeed()

	ensureTestPinsSetup()

	assert.NotNil(t, testPinChannels,
		fmt.Sprintf("TestPinChannels was nil: %v\n", testPinChannels))
	assert.True(t, len(testPinChannels) > 0)
	assert.True(t, channelsAvailable.isChannelAvailable("C3R0V1AU8"), "Channel dungar-test should exist.")
}

func oldTestPinChannelRequest(t *testing.T) {
	random.UseTestSeed()
	ensureTestPinsSetup()

	resp, ok := requestPinsForChannel(testPinChannels[0].ChannelID, "", "")
	assert.True(t, ok)
	assert.True(t, len(resp) > 0)
}

func oldTestPinHandler(t *testing.T) {
	random.UseTestSeed()
	ensureTestPinsSetup()

	var (
		chanID = ""
	)

	for _, chn := range testPinChannels {
		if chn.ChannelName == "butts" {
			chanID = chn.ChannelID
			break
		}
	}

	if len(chanID) == 0 {
		assert.Fail(t, "ChanID failure: could not find #butts' channel ID")
		return
	}

	rsp := pinRestHandler(makeMessageWithChanID("!pins", "chalur", "#butts", chanID))
	assert.True(t, rsp[0].HandledMessage)
	assert.NotEqual(t, "", rsp[0].Contents)

	rsp = pinRestHandler(makeMessageWithChanID("!pin", "chalur", "#butts", chanID))
	assert.False(t, rsp[0].HandledMessage)
	assert.Equal(t, "", rsp[0].Contents)

	rsp = pinRestHandler(makeMessageWithChanID("!pins slaps", "chalur", "#butts", chanID))
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, strings.Contains(rsp[0].Contents, "buttstuff"), rsp[0].Contents)
}
