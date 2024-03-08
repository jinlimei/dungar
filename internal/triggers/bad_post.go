package triggers

import (
	"log"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var badWords = make([]string, 0)
var badWordCheck = int64(0)

func preBadWordsHandler(msg *core2.IncomingMessage) []*core2.Response {
	// Change the response to empty and consume it.
	if shouldTossMessage(msg.Contents) {
		return core2.CancelledRsp()
	}

	return core2.EmptyRsp()
}

func postBadWordsHandler(svc *core2.Service, responses []*core2.Response) []*core2.Response {
	for _, response := range responses {
		if shouldTossMessage(response.Contents) {
			response.Cancel()
		}
	}

	return responses
}

func shouldTossMessage(str string) bool {
	badWordListCheck()

	str = strings.ToLower(utils.CleanSpaces(str))

	if len(badWords) == 0 {
		return false
	}

	var (
		split = strings.Split(str, " ")
		toss  = false
	)

out:
	for _, word := range badWords {
		for _, piece := range split {
			if word == utils.TrimPunctuation(piece) {
				log.Printf("Found bad word '%s' in message '%s', so tossing\n", word, str)

				toss = true
				break out
			}
		}
	}

	return toss
}

func badWordListCheck() {
	now := time.Now().Unix()

	if (now - badWordCheck) > 300 {
		badWords = db.GetBadWords()
		badWordCheck = now
	}
}
