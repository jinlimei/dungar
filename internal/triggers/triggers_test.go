package triggers

import (
	"io"
	"net/http"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/learning"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func makeMessageWithChanID(msg, nick, channel, chanID string) *core2.IncomingMessage {
	out := makeMessage(msg, nick, channel)
	out.ChannelID = chanID

	return out
}

func makeMessage(msg, userID, chanID string) *core2.IncomingMessage {
	return &core2.IncomingMessage{
		UserID:        userID,
		ChannelID:     chanID,
		Contents:      msg,
		LowerContents: strings.ToLower(msg),
		ParsedContents: &core2.ParsedMessage{
			Tokens:    make([]core2.MessageToken, 0),
			Raw:       msg,
			Converted: msg,
		},
	}
}

var learnedAliceInWonderland = false

func retrieveAliceInWonderland() ([]byte, error) {
	rsp, err := http.Get("https://tabhome.app/alice_in_wonderland.txt")
	if err != nil {
		return nil, err
	}

	bod, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return bod, nil
}

func useAliceInWonderland() {
	if learnedAliceInWonderland {
		return
	}

	markovUsingV3()
	// Avoid a database hit, let's just use alice forever ok? ok
	lastLoadedM3 = time.Now()
	lastSpokeM3 = time.Now()

	res, err := retrieveAliceInWonderland()

	if err != nil {
		panic(err)
	}

	lines := learning.ReadABook(string(res))

	for _, line := range lines {
		m3.LearnString(line, cleaner.VariantBook)
	}

	learnedAliceInWonderland = true
}

func isEmptyRsp(rsp []*core2.Response) bool {
	if len(rsp) == 0 {
		return true
	}

	return !rsp[0].IsHandled() && !rsp[0].IsConsumed() && !rsp[0].IsCancelled()
}
