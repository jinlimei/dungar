package triggers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/shttp"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

//
//import (
//	"fmt"
//	"strings"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//)
//
//func TestDadJokesPuller (t *testing.T) {
//	words := getRemoteDadJoke()
//	assert.True(t, !strings.Contains(words, "explodes"),
//		fmt.Sprintf("Message came back as '%s'", words))
//}

const validDadJoke = `{"id":"vPmy5EtPKuc","joke":"Do you want a brief explanation of what an acorn is? In a nutshell, it's an oak tree.","status":200}`

var dadJokeMockClient *shttp.HTTPMock

func TestDadJokeStruct(t *testing.T) {
	djr := dadJokeResponse{
		ID:     "",
		Joke:   "  hello   world  ",
		Status: 0,
	}

	djr.format()
	assert.Equal(t, "hello   world", djr.Joke)

	chk := djr.checksum()
	assert.Equal(t, "f06343efbbc9bcac096b71e819118b94deb6f1d8d351ce75ddac0c7d2d86ca18", chk,
		fmt.Sprintf("DJR Checksum(%s) did not match expected", chk))
}

func TestGetRemoteDadJoke(t *testing.T) {
	dadJokeMockClient = shttp.NewMock()
	dadJokeClient = dadJokeMockClient

	dadJokeMockClient.SetResponse(&shttp.SimpleResponse{
		Code: http.StatusOK,
		Body: []byte(validDadJoke),
	})

	rsp := getRemoteDadJoke()
	assert.Equal(t, "vPmy5EtPKuc", rsp.ID)
	assert.Equal(t, int64(200), rsp.Status)

	dadJokeMockClient.SetError(errors.New("ohgod"))

	rsp = getRemoteDadJoke()
	assert.Equal(t, "", rsp.ID)
	assert.Equal(t, int64(500), rsp.Status)
	assert.Contains(t, rsp.Joke, "exploded")

	dadJokeMockClient.SetError(nil)
	dadJokeMockClient.SetResponse(&shttp.SimpleResponse{
		Code: http.StatusBadGateway,
		Body: nil,
	})

	rsp = getRemoteDadJoke()
	assert.Equal(t, "", rsp.ID)
	assert.Equal(t, int64(500), rsp.Status)
	assert.Contains(t, rsp.Joke, "kinda exploded")

	dadJokeMockClient.SetResponse(&shttp.SimpleResponse{
		Code: http.StatusOK,
		Body: []byte(`{"id":"","AAAAAAAAAAAAAAAAAAAAAAAA`),
	})

	rsp = getRemoteDadJoke()
	assert.Equal(t, "", rsp.ID)
	assert.Equal(t, int64(500), rsp.Status)
	assert.Contains(t, rsp.Joke, "overwhelmingly exploded")

	dadJokeMockClient.SetResponse(&shttp.SimpleResponse{
		Code: http.StatusOK,
		Body: []byte(`{"id":"","joke":"","status":400}`),
	})

	rsp = getRemoteDadJoke()
	assert.Equal(t, "", rsp.ID)
	assert.Equal(t, int64(400), rsp.Status)
	assert.Contains(t, rsp.Joke, "weird happened")
}

func TestDadJokeHandler(t *testing.T) {
	random.UseTestSeed()
	db.TestDatabaseConnect()
	log.Printf("MarkovUsingV3!\n")
	markovUsingV3()

	msg := &core2.IncomingMessage{
		Contents: "aaaaaaaaaaaaaaaaa",
	}

	rsp := dadJokeHandler(msg)

	assert.Equal(t, core2.EmptyRsp(), rsp)

	msg.Contents = "!dadjoke"

	lastDadJoke = time.Now()
	rsp = dadJokeHandler(msg)
	assert.Equal(t, core2.EmptyRsp(), rsp)

	dadJokeMockClient = shttp.NewMock()
	dadJokeClient = dadJokeMockClient

	dadJokeMockClient.SetResponse(&shttp.SimpleResponse{
		Code: http.StatusOK,
		Body: []byte(validDadJoke),
	})

	nutshellCount := 0
	for k := 0; k < 10; k++ {
		lastDadJoke = time.Now().Add(-24 * time.Hour)
		rsp = dadJokeHandler(msg)

		if strings.Contains(rsp[0].Contents, "nutshell") {
			nutshellCount++
		} else {
			log.Printf("dadJoke rsp: %s\n", rsp[0].Contents)
		}
	}

	assert.True(t, nutshellCount > 0,
		fmt.Sprintf("failed nutshell count: %d", nutshellCount))

}
