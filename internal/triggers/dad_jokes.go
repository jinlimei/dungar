package triggers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/shttp"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var dadJokeClient shttp.HTTPRequester
var lastDadJoke = time.Now()

type dadJokeResponse struct {
	ID     string `json:"id"`
	Joke   string `json:"joke"`
	Status int64  `json:"status"`
}

func (djr *dadJokeResponse) format() {
	djr.Joke = strings.TrimSpace(djr.Joke)
}

func (djr *dadJokeResponse) checksum() string {
	h := sha256.New()

	h.Write([]byte(
		strings.ToLower(djr.Joke),
	))

	rs := h.Sum(nil)

	return hex.EncodeToString(rs)
}

func dadJokeHandler(msg *core2.IncomingMessage) []*core2.Response {
	now := time.Now()
	dadJokeWaitLimit := (now.Unix() - lastDadJoke.Unix()) > 60

	if strings.HasPrefix(msg.Contents, "!dadjoke") && dadJokeWaitLimit {
		lastDadJoke = now

		if fromBasicChance("dadJokeHandler--markov") {
			return core2.MakeSingleRsp(markovGenerate("dad"))
		}

		joke := getRemoteDadJoke()

		if joke.ID != "" {
			storeDadJoke(joke)
		}

		return core2.MakeSingleRsp(joke.Joke)
	}

	return core2.EmptyRsp()
}

func dadJokeStr(str string) dadJokeResponse {
	return dadJokeResponse{
		ID:     "",
		Joke:   str,
		Status: 500,
	}
}

func getRemoteDadJoke() dadJokeResponse {
	if dadJokeClient == nil {
		dadJokeClient = shttp.New()
	}

	req := &shttp.SimpleRequest{
		Method: http.MethodGet,
		URL:    "https://icanhazdadjoke.com/",
		Headers: map[string]string{
			"Accept":     "application/json",
			"User-Agent": "discord bot (jinli.mei.eve@gmail.com)",
		},
	}

	rsp, err := dadJokeClient.Request(req)

	if err != nil {
		log.Printf("Failed to make a DadJoke API Request: %v\n", err)
		return dadJokeStr("dad jokes exploded :smith:")
	}

	if rsp.Code != http.StatusOK {
		log.Printf("Received non-200 from DadJoke API: %d\n", rsp.Code)
		return dadJokeStr("dad jokes kinda exploded :smith:")
	}

	var joke dadJokeResponse

	err = json.Unmarshal(rsp.Body, &joke)

	if err != nil {
		log.Printf("Failed to decode DadJoke API Body: %v\n", err)
		return dadJokeStr("dad jokes like, overwhelmingly exploded :fgf:")
	}

	if joke.Status != http.StatusOK {
		log.Printf("Inline status is not-200: %d\n", joke.Status)
		out := dadJokeStr(fmt.Sprintf("something weird happened (djs: %d)", joke.Status))
		out.Status = joke.Status
		return out
	}

	joke.format()

	return joke
}

func storeDadJoke(rsp dadJokeResponse) {
	if utils.InTestEnv() {
		return
	}

	chk := rsp.checksum()

	sel := `
		SELECT id, joke_id
		FROM dad_jokes
		WHERE joke_checksum = $1
	`

	var (
		id     uint64
		jokeID string
	)

	row := db.ConQueryRow(sel, chk)
	err := row.Err()

	if err != nil {
		log.Printf("Failed to retrieve Dad Jokes: %v\n", err)
		return
	}

	err = row.Scan(&id, &jokeID)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Failed to Scan: %v\n", err)
	}

	if err == sql.ErrNoRows {
		ins := `
			INSERT INTO dad_jokes (the_joke, joke_checksum, joke_id, created_at, last_used)
			VALUES($1, $2, $3, CURRENT_TIMESTAMP, NULL)
		`

		_, err = db.ConExec(ins, rsp.Joke, chk, rsp.ID)
		if err != nil {
			log.Printf("Failed to store dad joke: %v\n", err)
		}
	} else {
		if jokeID != rsp.ID {
			log.Printf("Something weird: Internal JokeID (%s) does not match remote ID (%s)\n",
				jokeID, rsp.ID)
		}

		upd := `
			UPDATE dad_jokes
			SET last_used = CURRENT_TIMESTAMP
			WHERE id = $1
		`

		_, err = db.ConExec(upd, id)
	}
}
