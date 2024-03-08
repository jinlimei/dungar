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
		{4, "UserSpammed", userSpammedHandler},

		// Guinea Pig Handler (who to bully)
		// TODO implement this coherently
		//{5, "GuineaPig", guineaPigHandler},

		// No You (Fuck You -> NO YOU)
		{6, "NoYou", core2.WrapNoServiceFunc(noYouHandler)},

		// GDPR
		//{7, "GDPR", gdprHandler},

		// !pins
		//{8, "!pins", core2.WrapNoServiceFunc(pinRestHandler)},
		{8, "!pins", core2.WrapNoServiceFunc(pinDBHandler)},

		// !dadjoke
		{9, "DadJoke", core2.WrapNoServiceFunc(dadJokeHandler)},

		// !bitcoin
		{10, "BitCoin", core2.WrapNoServiceFunc(bitCoinValueHandler)},

		// !alexjons and !hannity
		{11, "AlexJones", alexJonesHandler},

		// !f (fortune teller)
		{12, "Fortune", core2.WrapNoServiceFunc(fortuneHandler)},

		// Name Check (Get My Name Out of Your Mouth)
		{13, "NameCheck", nameCheckHandler},

		// Range Handler (pick a number between a and b)
		{14, "Range", rangeHandler},

		// Roll Die (1d20)
		{15, "RollDice", rollDiceHandler},

		// Coin Flip (heads/tails also haha sides)
		{16, "CoinFlip", coinFlipHandler},

		// Chances (hey go outside c/d)
		{17, "Chances", chanceGameHandler},

		// Questions answered by 8Ball
		{18, "8Ball", questionsHandler},

		{18, "Saucer", core2.WrapNoServiceFunc(sauceQuestionsHandler)},

		{18, "Sequencer", core2.WrapNoServiceFunc(sequenceHandler)},

		// Random Markov Commentary
		{19, "Commentator", core2.WrapNoServiceFunc(commentateHandler)},
	})

	messageGroup.SetResponseTriggers([]core2.ResponseEvHandler{
		{1, core2.Filter, "GarbagePrefixes", removeGarbagePrefixedHandler},
		{2, core2.Filter, "WhiteSpace", lotsOfSpaceHandler},
		{3, core2.Filter, "WeirdOutput", postWeirdOutputFilterHandler},
		{4, core2.Adder, "SelfSpamming", selfSpamHandler},
	})

	var out []*core2.Response
	svc := initMockServices()
	msg := makeMessage("mercy", "george", "main")

	for k := 0; k < 25; k++ {
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

	if out == nil || len(out) == 0 {
		assert.Fail(t, "Out is Nil")
		return
	}

	assert.True(t, out[0].HandledMessage)
	assert.NotEmpty(t, out[0].Contents)
}
