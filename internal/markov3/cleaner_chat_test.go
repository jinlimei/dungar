package markov3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	msg := "Hey friends :sun: what's going on? I just went to " +
		"https://www.google.com or maybe <@thing> or maybe @thing might like it idk " +
		"but also check out <https://www.reddit.com> yeah ok cool bye :sun:"

	res := Tokenize(msg)

	assert.True(t, res.Tokens[4].TokenType == ChatTokenEmoticon)
	//spew.Dump(res.Tokens)
	//spew.Dump(res.MarkovConsumable())

	msg = "@dungar dungar or executive order for a month as long as your regular pop like ten days"

	res = Tokenize(msg)

	assert.True(t, res.Tokens[0].TokenType == ChatTokenSpecial)

	//spew.Dump(res.Tokens)
	//spew.Dump(res.MarkovConsumable())
}

func TestTokenizeCode(t *testing.T) {
	msg := "``` ```"
	res := Tokenize(msg)

	assert.True(t, res.Tokens[0].TokenType == ChatTokenCode)

	msg = "`0.9 6.3 8271900 1169016`"

	res = Tokenize(msg)

	assert.True(t, res.Tokens[0].TokenType == ChatTokenCode)
}

func TestTokenizeCodePart2(t *testing.T) {
	//msg := ":3: `creation date: 2005-12-25t17:06:21z`"
	msg := "idk how to get raw `` through tho. any of those in your message get encoded to `&lgt;` from slack. only making urls or @'ing people makes the actual chars"
	res := Tokenize(msg)

	assert.True(t, res.Tokens[10].TokenType == ChatTokenCode)
}

func TestTokenLongLine(t *testing.T) {
	msg := "2 of these in it <https://ark.intel.com/content/www/us/en/ark/products/120495/intel-xeon-gold-6154-processor-24-75m-cache-3-00-ghz.html>"
	res := Tokenize(msg)

	assert.True(t, res.Tokens[10].TokenType == ChatTokenURL)
}

func TestTokenAnalyzer(t *testing.T) {
	var (
		tok = ChatToken{
			TokenType: ChatTokenEmoticon,
			Value:     []rune(":abc::"),
		}
	)

	tokRes := tok.Analyze()
	assert.Equal(t, 2, len(tokRes))
	assert.Equal(t, ChatTokenEmoticon, tokRes[0].TokenType)
	assert.Equal(t, ":abc:", string(tokRes[0].Value))
	assert.Equal(t, ChatTokenWord, tokRes[1].TokenType)
	assert.Equal(t, ":", string(tokRes[1].Value))
}

func TestTokenEmoticon(t *testing.T) {
	var (
		msg = "here :sun: and an example: foo"
		res = Tokenize(msg)
	)

	assert.True(t, res.Tokens[8].TokenType == ChatTokenWord,
		"the word 'example:' should not be considered an imaginary emote")

	//msg = "uh hello: :rrrr:: ::: :"
	msg = "uh hello: :rrrr::"
	res = Tokenize(msg)

	assert.True(t, res.Tokens[4].TokenType == ChatTokenWord,
		fmt.Sprintf("Token isn't Word: %+v\n", res.Tokens[4].String()))

	msg = "uh hello: :rrrr:: ::: :"
	//msg = "uh hello: :rrrr::"
	res = Tokenize(msg)

	assert.True(t, res.Tokens[4].TokenType == ChatTokenEmoticon,
		fmt.Sprintf("Token isn't Emoticon: %+v\n", res.Tokens[4].String()))

	assert.Equal(t, ":rrrr:", string(res.Tokens[4].Value))

	assert.True(t, res.Tokens[7].TokenType == ChatTokenWord,
		fmt.Sprintf("Token isn't Word: %+v\n", res.Tokens[7].String()))

	assert.True(t, res.Tokens[9].TokenType == ChatTokenWord,
		fmt.Sprintf("Token isn't Word: %+v\n", res.Tokens[9].String()))

}

func TestTokenizeDiscord(t *testing.T) {
	var (
		msg = "Hey there <:durrhuREEE:546076030666801177>"
		res = Tokenize(msg)
	)

	assert.True(t, len(res.Tokens) > 0)
	assert.True(t, res.Tokens[4].TokenType == ChatTokenEmoticon)

	msg = "<@&766488443248836609> this is pinging the example role"
	res = Tokenize(msg)

	assert.True(t, len(res.Tokens) > 0)
	assert.True(t, res.Tokens[0].TokenType == ChatTokenSpecial)
}

func TestTokenUnknown(t *testing.T) {
	//msg := "$URL$ $EMOTICON$ $UNKNOWN|we're|$ excited to see it as a regular basis"
	var (
		msg = "free market conservatives: we're going to de-regulate the markets and let them solve their own problems!\nalso free market conservatives: the market is fixing itself in a way I don't like! let's regulate them"
		res = Tokenize(msg)
	)

	assert.True(t, res.Tokens[4].TokenType == ChatTokenWord)
	assert.True(t, res.Tokens[5].TokenType == ChatTokenSpace)
	assert.True(t, res.Tokens[6].TokenType == ChatTokenWord)

	msg = ":musical_note: we're gonna lose the house :musical_note:"
	res = Tokenize(msg)

	assert.True(t, res.Tokens[0].TokenType == ChatTokenEmoticon)
	assert.True(t, res.Tokens[1].TokenType == ChatTokenSpace)
	assert.True(t, res.Tokens[2].TokenType == ChatTokenWord)

	assert.True(t, res.Tokens[10].TokenType == ChatTokenWord)
	assert.True(t, res.Tokens[11].TokenType == ChatTokenSpace)
	assert.True(t, res.Tokens[12].TokenType == ChatTokenEmoticon,
		fmt.Sprintf("Token: %+v\n", res.Tokens[12].String()))

	msg = "been a good morning... realized my previous landlord is in violation of rental laws and is withholding our security deposit past the duration he is legally allowed to which means we're entitled to like 3x the amount withheld, plus $100, plus attorney fees  :notbadobama:"
	res = Tokenize(msg)

	assert.Equal(t, msg, res.String(),
		"Reconstructing string from tokens failed!")
}
