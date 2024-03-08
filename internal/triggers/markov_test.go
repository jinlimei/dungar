package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestMarkovPickWord(t *testing.T) {
	useAliceInWonderland()
	random.UseTestSeed()

	var word string

	//activeMarkovVersion = mV1
	//
	//word := markovPickWord()
	//assert.NotEmpty(t, word)

	activeMarkovVersion = mV3
	word = markovPickWord()
	assert.NotEmpty(t, word)
}

func TestMarkovGenerate(t *testing.T) {
	useAliceInWonderland()
	random.UseTestSeed()

	var msg string

	//activeMarkovVersion = mV1
	//
	//msg := markovGenerate("alice")
	//assert.NotEmpty(t, msg)

	activeMarkovVersion = mV3

	msg = markovGenerate("alice")
	assert.NotEmpty(t, msg)
	assert.False(t, lastLoadedM3.IsZero())
}

func TestMarkovHandler(t *testing.T) {
	initMockServices()
	useAliceInWonderland()
	random.UseTestSeed()
	db.TestDatabaseConnect()

	msg := makeMessage("hello", "U3Q9ZPR32", "wonderland")
	rsp := markovSettingHandler(msg)
	assert.Equal(t, core2.EmptyRsp(), rsp)

	msg.Contents = "!m"
	rsp = markovSettingHandler(msg)
	assert.Equal(t, core2.EmptyRsp(), rsp)

	//msg.Contents = "!m1"
	//rsp = markovSettingHandler(msg)
	//assert.True(t, rsp[0].HandledMessage)
	//assert.True(t, rsp[0].ConsumedMessage)
	//
	//msg.Contents = "!m1 alice"
	//rsp = markovSettingHandler(msg)
	//assert.True(t, rsp[0].HandledMessage)
	//assert.True(t, rsp[0].ConsumedMessage)

	msg.Contents = "!m3"
	rsp = markovSettingHandler(msg)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)

	msg.Contents = "!m3 alice"
	rsp = markovSettingHandler(msg)
	assert.True(t, rsp[0].HandledMessage)
	assert.True(t, rsp[0].ConsumedMessage)

	msg.UserID = "kinkouin"
	msg.Contents = "!markov get-version"

	hasNotYou := false
	for k := 0; k < 50; k++ {
		rsp = markovSettingHandler(msg)

		if strings.Contains(rsp[0].Contents, "TOUCH") {
			hasNotYou = true
			break
		}
	}

	assert.True(t, hasNotYou, "Had a NOT YOU experience")
}

func TestMarkovSettingHandlerAsJinli(t *testing.T) {
	initMockServices()
	useAliceInWonderland()
	random.UseTestSeed()
	db.TestDatabaseConnect()

	//log.Println("TestMarkovSettingHandlerAsJinli")

	msg := makeMessage("!markov get-version", "U3Q9ZPR32", "dungar-test")

	//log.Printf("msg '%+v'", msg)

	setMarkovVersion(mV1)
	//log.Printf("[SHOULD REPORT V1] Active Markov Version %d", activeMarkovVersion)
	rsp := markovSettingHandler(msg)
	//log.Printf("[mV1] RESPONSE: %+v", rsp)
	assert.Contains(t, rsp[0].Contents, "Actively")
	assert.Contains(t, rsp[0].Contents, "1")

	setMarkovVersion(mV3)
	//log.Printf("[SHOULD REPORT V3] Active Markov Version %d", activeMarkovVersion)
	rsp = markovSettingHandler(msg)
	//log.Printf("[mV3] RESPONSE: %+v", rsp)
	assert.Contains(t, rsp[0].Contents, "Actively")
	assert.Contains(t, rsp[0].Contents, "3")

	msg.Contents = "!markov set-version 1"
	rsp = markovSettingHandler(msg)
	//log.Printf("[set-version 1] RESPONSE: %+v", rsp)
	assert.Equal(t, mV1, activeMarkovVersion)

	msg.Contents = "!markov set-version 3"
	rsp = markovSettingHandler(msg)
	//log.Printf("[set-version 2] RESPONSE: %+v", rsp)
	assert.Equal(t, mV3, activeMarkovVersion)

	msg.Contents = "!markov set-version"
	rsp = markovSettingHandler(msg)
	//log.Printf("[set-version none] RESPONSE: %+v", rsp)
	assert.Contains(t, rsp[0].Contents, "WHAT")

	msg.Contents = "!markov set-version blue"
	rsp = markovSettingHandler(msg)
	//log.Printf("[set-version blue] RESPONSE: %+v", rsp)
	assert.Contains(t, rsp[0].Contents, "fucked up")

	msg.Contents = "!markov set-version 4"
	rsp = markovSettingHandler(msg)
	//log.Printf("[set-version 4] RESPONSE: %+v", rsp)
	assert.Equal(t, core2.EmptyRsp(), rsp)

	// causes a data race :(
	//msg.Contents = "!markov load-m3"
	//rsp = markovSettingHandler(msg)
	//assert.Contains(t, rsp[0].Contents, "whatever")

	msg.Contents = "!markov m3-stats"
	rsp = markovSettingHandler(msg)
	//log.Printf("[m3-stats] RESPONSE: %+v", rsp)
	assert.Contains(t, rsp[0].Contents, "Words")
	assert.Contains(t, rsp[0].Contents, "Fragments")
}
