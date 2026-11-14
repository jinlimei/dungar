package triggers

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var msgGroup *core2.MessageTriggerGroup

// incomingMessageHandler Handles all assorted messageTriggers and
// returns a slice of sentences
func incomingMessageHandler(msg *core2.IncomingMessage) *core2.ResponseEnvelope {
	log.Printf(
		"[incomingMessageHandler serverID='%s' userID='%s'] %s\n",
		msg.ServerID,
		msg.UserID,
		msg.String(),
	)

	if msgGroup == nil {
		initMessageGroup()
	}

	msgGroup.Process(core, msg)

	var (
		responses = msgGroup.Responses()
		out       = make([]*core2.Response, 0)
	)

	log.Printf("[incomingMessageHandler] Responses from messageGroup: (%d)\n",
		len(responses))
	if len(responses) > 0 {
		log.Print(spew.Sdump(responses))
	}

	var (
		empty     = 0
		cancelled = 0
	)

	for _, response := range responses {
		if !response.IsEmpty() && !response.IsCancelled() {
			out = append(out, response)
		} else if response.IsCancelled() {
			cancelled++
		} else if response.IsEmpty() {
			empty++
		}
	}

	if empty > 0 || cancelled > 0 {
		log.Printf("[handleMessage] Output: %d, Cancelled: %d, Empty: %d\n", len(out), cancelled, empty)
	}

	return &core2.ResponseEnvelope{
		Message:   msg,
		Responses: out,
	}
}

func initMessageGroup() {
	msgGroup = &core2.MessageTriggerGroup{}

	msgGroup.SetMessageTriggers([]core2.MessageEvHandler{
		// Pre-filtering of Bad Words
		{Order: 0, Name: "PreBadWords", Func: core2.WrapNoServiceFunc(preBadWordsHandler)},
		// {1, "TestMode", testModeHandler},

		// Raw Message Recorder
		{Order: 1, Name: "Recorder", Func: rawMessageRecorder},

		// User Tracking
		{Order: 2, Name: "UserTracker", Func: userTrackingHandler},

		// Debugger!
		{Order: 2, Name: "Debugger", Func: debugHandler},

		// Test Mode

		// Markov Commands (version, analytics)
		{Order: 3, Name: "MarkovSetting", Func: core2.WrapNoServiceFunc(markovSettingHandler)},

		// Random analytics (for things that are setup w/ it)
		{Order: 3, Name: "RandomAnalytics", Func: core2.WrapNoServiceFunc(randomHandler)},

		// Main Control (to disable/change wild triggers)
		{Order: 3, Name: "MainControl", Func: core2.WrapNoServiceFunc(mainControlHandler)},

		// User Spamming (stop spamming dungar, nerd!)
		{Order: 4, Name: "UserSpammed", Func: userSpammedHandler},

		// Guinea Pig Handler (who to bully)
		// TODO implement this coherently
		// {Order:5, Name: "GuineaPig", Func: guineaPigHandler},

		// GDPR
		// {Order:6, Name: "GDPR", Func: gdprHandler},

		// !pins
		// {Order:7, Name: "!pins", Func: core2.WrapNoServiceFunc(pinRestHandler)},
		{Order: 7, Name: "!pins", Func: core2.WrapNoServiceFunc(pinDBHandler)},

		// !dadjoke
		{Order: 8, Name: "DadJoke", Func: core2.WrapNoServiceFunc(dadJokeHandler)},

		// !bitcoin
		{Order: 9, Name: "BitCoin", Func: core2.WrapNoServiceFunc(bitCoinValueHandler)},

		// !alexjons and !hannity
		{Order: 10, Name: "AlexJones", Func: alexJonesHandler},

		// !f (fortune teller)
		{Order: 11, Name: "Fortune", Func: core2.WrapNoServiceFunc(fortuneHandler)},

		// No You (Fuck You -> NO YOU)
		{Order: 12, Name: "NoYou", Func: core2.WrapNoServiceFunc(noYouHandler)},

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

		// Markov Extended ("Tell me about shakespeare")
		{Order: 17, Name: "Markov Extended", Func: markovExtendedHandler},

		// Questions answered by 8Ball
		{Order: 18, Name: "8Ball", Func: questionsHandler},

		// Sequence handler
		{Order: 19, Name: "Sequences", Func: core2.WrapNoServiceFunc(sequenceHandler)},

		// Prefix triggers (dovibus)
		{Order: 20, Name: "Simple Prefix", Func: core2.WrapNoServiceFunc(simplePrefixHandler)},

		// Random triggers
		{Order: 20, Name: "Sauce Questions", Func: core2.WrapNoServiceFunc(sauceQuestionsHandler)},

		// Repeatable messaging
		{Order: 21, Name: "Repeatable", Func: core2.WrapNoServiceFunc(repeatableHandler)},

		// Random Markov Commentary
		{Order: 22, Name: "Commentator", Func: core2.WrapNoServiceFunc(commentateHandler)},
	})

	msgGroup.SetResponseTriggers([]core2.ResponseEvHandler{
		{Order: 0, Type: core2.Filter, Name: "BadWords", Func: postBadWordsHandler},
		{Order: 1, Type: core2.Filter, Name: "GarbagePrefixes", Func: removeGarbagePrefixedHandler},
		{Order: 2, Type: core2.Filter, Name: "WhiteSpace", Func: lotsOfSpaceHandler},
		// {Order:3, Type: core2.Filter, Name: "WeirdOutput", Func: postWeirdOutputFilterHandler},
		{Order: 4, Type: core2.Adder, Name: "SelfSpamming", Func: selfSpamHandler},
	})
}
