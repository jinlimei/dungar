package triggers

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"gitlab.int.magneato.site/dungar/prototype/internal/markov3"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

type pinChannelsResponse struct {
	TeamID  string       `json:"teamId"`
	Results []pinChannel `json:"results"`
	Total   int          `json:"totalResults"`
}

type pinChannel struct {
	TeamID      string `json:"teamId"`
	ChannelID   string `json:"ChannelID"`
	ChannelName string `json:"ChannelName"`
}

type pinResponse struct {
	TeamID    string `json:"teamId"`
	ChannelID string `json:"channelId"`
	Results   []pin  `json:"results"`
	Total     int    `json:"totalResults"`
}

type pin struct {
	EventID     string    `json:"eventId"`
	TeamID      string    `json:"teamId"`
	ChannelID   string    `json:"channelId"`
	ChannelName string    `json:"channelName"`
	UserID      string    `json:"userId"`
	UserDisplay string    `json:"userDisplay"`
	MessageText string    `json:"messageText"`
	MessageTime time.Time `json:"messageTime"`
}

type pinConf struct {
	URL    string
	Auth   string
	TeamID string
}

type pinChannelsAvailable struct {
	LastRequest int64
	Channels    []pinChannel
}

var pinMarkovs = make(map[markov3.MarkovSpaceID]*markov3.Markov)

func (pca *pinChannelsAvailable) ensureChannelsAvailable() {
	nowTime := time.Now().Unix()

	if (nowTime - pca.LastRequest) < 3600 {
		return
	}

	pca.Channels = requestPinChannels()
	pca.LastRequest = nowTime
}

func (pca *pinChannelsAvailable) isChannelAvailable(chanID string) bool {
	pca.ensureChannelsAvailable()
	available := false

	for _, availChan := range pca.Channels {
		if availChan.ChannelID == chanID {
			available = true
			break
		}
	}

	return available
}

var activePinConf *pinConf
var channelsAvailable = &pinChannelsAvailable{
	LastRequest: 0,
	Channels:    make([]pinChannel, 0),
}

func ensurePinConf() *pinConf {
	if activePinConf != nil {
		return activePinConf
	}

	credentials := utils.PinCredentials()

	activePinConf = &pinConf{
		URL:    credentials["url"],
		Auth:   credentials["auth"],
		TeamID: credentials["team"],
	}

	activePinConf.validate()

	return activePinConf
}

func pinRestHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !strings.HasPrefix(msg.Contents, "!pins") {
		return core2.EmptyRsp()
	}

	if !channelsAvailable.isChannelAvailable(msg.ChannelID) {
		return core2.MakeSingleRsp(pinGenerateFakeResponse(msg.ChannelID))
	}

	search := ""

	// we're searching for our long lost love?
	// "!pins" = 5 + 1 (space) + 1 (char) = 7
	if len(msg.Contents) > 7 {
		search = strings.TrimSpace(msg.Contents[5:])
	}

	pinList, ok := requestPinsForChannel(msg.ChannelID, search, "")

	// Just markov if we can't find any pins
	if !ok || len(pinList) == 0 {
		return core2.MakeSingleRsp(pinGenerateFakeResponse(msg.ChannelID))
	}

	picked := random.Int(len(pinList))
	pin := pinList[picked]

	return core2.MakeSingleRsp(fmt.Sprintf("%s\n> pin from %s", pin.MessageText, pin.UserDisplay))
}

func (pc *pinConf) buildRequestURL(endPoint string) string {
	return pc.URL + pc.TeamID + endPoint
}

func (pc *pinConf) validate() {
	if pc.TeamID == "" {
		panic("pinConf TeamID must not be blank")
	}

	if pc.Auth == "" {
		panic("pinConf Auth must not be blank")
	}

	if pc.URL == "" {
		panic("pinConf URL must not be blank")
	}
}

func pinGenerateFakeResponse(chanID string) string {
	chanMarkovID := markov3.MarkovSpaceID(chanID)

	marko, ok := pinMarkovs[chanMarkovID]

	if !ok {
		marko = markov3.MakeMarkov(chanMarkovID)

		// Let's get _all_ pins for the channel and then move it
		pins, pinsOK := requestPinsForChannel(chanID, "", "")

		if pinsOK {
			for _, pin := range pins {
				marko.LearnString(pin.MessageText, cleaner.VariantSlack)
			}
		}

		pinMarkovs[chanMarkovID] = marko
	}

	var text string
	if len(marko.Fragments) > 0 {
		text = marko.Generate(marko.GetRandomWordStr())
	} else {
		text = markovGenerate(markovPickWord())
	}

	return fmt.Sprintf("%s\n> - Actually From Dungar", text)
}

func pinRequestHeaders() map[string]string {
	ensurePinConf()

	return map[string]string{
		"Authorization": "Bearer " + activePinConf.Auth,
	}
}

func requestPinChannels() []pinChannel {
	reqURL := ensurePinConf().buildRequestURL("/channels")

	res, err := utils.MakeGetRequest(reqURL, pinRequestHeaders())

	utils.HaltingError("requestPinChannels", err)

	var channels pinChannelsResponse

	err = json.Unmarshal(res, &channels)

	if err != nil {
		utils.NonHaltingError("requestPinChannels json unmarshal", err)
		return nil
	}

	if channels.Total == 0 {
		return nil
	}

	return channels.Results
}

func requestPinsForChannel(channelID string, search string, userID string) ([]pin, bool) {
	reqURL := ensurePinConf().buildRequestURL("/channels/" + channelID + "/pins?")

	if search != "" {
		reqURL += "messageText=" + url.QueryEscape(search) + "&"
	}

	if userID != "" {
		reqURL += "userId=" + url.QueryEscape(userID)
	}

	res, err := utils.MakeGetRequest(reqURL, pinRequestHeaders())

	if err != nil {
		utils.NonHaltingError("requestPinsForChannel", err)
		return nil, false
	}

	var pinsResp pinResponse

	err = json.Unmarshal(res, &pinsResp)

	if err != nil {
		utils.NonHaltingError("unmarshal requestPinsForChannel", err)
		return nil, false
	}

	if pinsResp.Total == 0 {
		return nil, false
	}

	return pinsResp.Results, true
}
