package triggers

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestSequenceMessage_Next(t *testing.T) {
	random.UseTestSeed()
	seq := &sequenceMessage{"oof", "ouch", "owie"}

	assert.Equal(t, "ouch", seq.next("oof"))
	assert.Equal(t, "owie", seq.next("ouch"))
	assert.Equal(t, "oof", seq.next("owie"))

	assert.Equal(t, ":thonkeng:", randSequences[0].next(":thinking:"))
	assert.Equal(t, ":universing:", randSequences[0].random(":thinking:"))
}

func TestHandlerOtherTriggersInChain(t *testing.T) {
	random.UseTestSeed()

	messageGroup := &core2.MessageTriggerGroup{}

	messageGroup.SetMessageTriggers([]core2.MessageEvHandler{
		// User Spamming (stop spamming dungar, nerd!)
		{Order: 4, Name: "UserSpammed", Func: userSpammedHandler},

		// Guinea Pig Handler (who to bully)
		// TODO implement this coherently
		// {5, "GuineaPig", guineaPigHandler},

		// No You (Fuck You -> NO YOU)
		{Order: 6, Name: "NoYou", Func: core2.WrapNoServiceFunc(noYouHandler)},

		// GDPR
		// {7, "GDPR", gdprHandler},

		// !pins
		// {8, "!pins", core2.WrapNoServiceFunc(pinRestHandler)},
		{Order: 8, Name: "!pins", Func: core2.WrapNoServiceFunc(pinDBHandler)},

		// !dadjoke
		{Order: 9, Name: "DadJoke", Func: core2.WrapNoServiceFunc(dadJokeHandler)},

		// !bitcoin
		{Order: 10, Name: "BitCoin", Func: core2.WrapNoServiceFunc(bitCoinValueHandler)},

		// !alexjons and !hannity
		{Order: 11, Name: "AlexJones", Func: alexJonesHandler},

		// !f (fortune teller)
		{Order: 12, Name: "Fortune", Func: core2.WrapNoServiceFunc(fortuneHandler)},

		// Name Check (Get My Name Out of Your Mouth)
		{Order: 13, Name: "NameCheck", Func: nameCheckHandler},

		// Range Handler (pick a number between a and b)
		{Order: 14, Name: "Range", Func: rangeHandler},

		// Roll Die (1d20)
		{Order: 15, Name: "RollDice", Func: rollDiceHandler},

		// Coin Flip (heads/tails also haha sides)
		{Order: 16, Name: "CoinFlip", Func: coinFlipHandler},

		// Chances (hey go outside c/d)
		{Order: 17, Name: "Chances", Func: chanceGameHandler},

		// Questions answered by 8Ball
		{Order: 18, Name: "8Ball", Func: questionsHandler},

		{Order: 18, Name: "Saucer", Func: core2.WrapNoServiceFunc(sauceQuestionsHandler)},

		{Order: 18, Name: "Sequencer", Func: core2.WrapNoServiceFunc(sequenceHandler)},

		// Random Markov Commentary
		{Order: 19, Name: "Commentator", Func: core2.WrapNoServiceFunc(commentateHandler)},
	})

	messageGroup.SetResponseTriggers([]core2.ResponseEvHandler{
		{Order: 1, Type: core2.Filter, Name: "GarbagePrefixes", Func: removeGarbagePrefixedHandler},
		{Order: 2, Type: core2.Filter, Name: "WhiteSpace", Func: lotsOfSpaceHandler},
		{Order: 3, Type: core2.Filter, Name: "WeirdOutput", Func: postWeirdOutputFilterHandler},
		{Order: 4, Type: core2.Adder, Name: "SelfSpamming", Func: selfSpamHandler},
	})

	var out []*core2.Response
	svc := initMockServices()
	msg := makeMessage("mercy", "george", "main")

	for range 25 {
		random.UseTimeBasedSeed()

		iters := 0
		limit := 10000
		found := false

		for {
			messageGroup.Process(svc, msg)
			out = messageGroup.Responses()

			if len(out) == 1 && strings.Contains(out[0].Contents, "mercy main btw") {
				found = true
				break
			}

			if iters >= limit {
				break
			}

			iters++
		}

		log.Printf("Found 'mercy main btw' after %d iterations\n", iters)
		assert.Truef(t, found, "After %d iterations, could not find mercy main btw", iters)
	}

	random.UseTestSeed()
}

func TestIsQuote(t *testing.T) {
	random.UseTestSeed()

	msg := makeMessage("> implying", "woet", "#butts")
	assert.True(t, isQuote(msg))
	msg = makeMessage(">implying", "woet", "#butts")
	assert.True(t, isQuote(msg))
	msg = makeMessage("things >things", "woet", "#butts")
	assert.False(t, isQuote(msg))
}

func TestCommentateHandler(t *testing.T) {
	useAliceInWonderland()
	random.UseTestSeed()

	msg := makeMessage(":thinking:", "jinli", "#butts")

	var out []*core2.Response

	masterChanceList["commentateHandler--commentate"] = 0.50

	for i := 0; i <= 50; i++ {
		out = commentateHandler(msg)

		if len(out) > 0 && out[0].HandledMessage {
			break
		}
	}

	if len(out) == 0 {
		assert.Fail(t, "Out is Nil")
		return
	}

	assert.True(t, out[0].HandledMessage)
	assert.NotEmpty(t, out[0].Contents)
}
