package triggers

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var msgGroup *core2.MessageTriggerGroup

// incomingMessageHandler Handles all assorted messageTriggers and
// returns an slice of sentences
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
		{0, "PreBadWords", core2.WrapNoServiceFunc(preBadWordsHandler)},
		//{1, "TestMode", testModeHandler},

		// Raw Message Recorder
		{1, "Recorder", rawMessageRecorder},

		// User Tracking
		{2, "UserTracker", userTrackingHandler},

		// Debugger!
		{2, "Debugger", debugHandler},

		// Test Mode

		// Markov Commands (version, analytics)
		{3, "MarkovSetting", core2.WrapNoServiceFunc(markovSettingHandler)},

		// Random analytics (for things that are setup w/ it)
		{3, "RandomAnalytics", core2.WrapNoServiceFunc(randomHandler)},

		// Main Control (to disable/change wild triggers)
		{3, "MainControl", core2.WrapNoServiceFunc(mainControlHandler)},

		// User Spamming (stop spamming dungar, nerd!)
		{4, "UserSpammed", userSpammedHandler},

		// Guinea Pig Handler (who to bully)
		// TODO implement this coherently
		//{5, "GuineaPig", guineaPigHandler},

		// GDPR
		//{6, "GDPR", gdprHandler},

		// !pins
		//{7, "!pins", core2.WrapNoServiceFunc(pinRestHandler)},
		{7, "!pins", core2.WrapNoServiceFunc(pinDBHandler)},

		// !dadjoke
		{8, "DadJoke", core2.WrapNoServiceFunc(dadJokeHandler)},

		// !bitcoin
		{9, "BitCoin", core2.WrapNoServiceFunc(bitCoinValueHandler)},

		// !alexjons and !hannity
		{10, "AlexJones", alexJonesHandler},

		// !f (fortune teller)
		{11, "Fortune", core2.WrapNoServiceFunc(fortuneHandler)},

		// No You (Fuck You -> NO YOU)
		{12, "NoYou", core2.WrapNoServiceFunc(noYouHandler)},

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

		// Markov Extended ("Tell me about shakespeare")
		{17, "Markov Extended", markovExtendedHandler},

		// Questions answered by 8Ball
		{18, "8Ball", questionsHandler},

		// Sequence handler
		{19, "Sequences", core2.WrapNoServiceFunc(sequenceHandler)},

		// Prefix triggers (dovibus)
		{20, "Simple Prefix", core2.WrapNoServiceFunc(simplePrefixHandler)},

		// Random triggers
		{20, "Sauce Questions", core2.WrapNoServiceFunc(sauceQuestionsHandler)},

		// Repeatable messaging
		{21, "Repeatable", core2.WrapNoServiceFunc(repeatableHandler)},

		// Random Markov Commentary
		{22, "Commentator", core2.WrapNoServiceFunc(commentateHandler)},
	})

	msgGroup.SetResponseTriggers([]core2.ResponseEvHandler{
		{0, core2.Filter, "BadWords", postBadWordsHandler},
		{1, core2.Filter, "GarbagePrefixes", removeGarbagePrefixedHandler},
		{2, core2.Filter, "WhiteSpace", lotsOfSpaceHandler},
		//{3, core2.Filter, "WeirdOutput", postWeirdOutputFilterHandler},
		{4, core2.Adder, "SelfSpamming", selfSpamHandler},
	})
}
