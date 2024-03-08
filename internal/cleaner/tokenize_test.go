package cleaner

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertTokenSequence(t *testing.T, expectedLen int, incoming []Token, expected []TokenType) {
	var (
		incLen = len(incoming)
		expLen = len(expected)
	)

	if expectedLen != incLen {
		assert.Fail(t, fmt.Sprintf(
			"Token length (%d) did not match expected length of %d",
			incLen,
			expectedLen,
		))

		DebugPrintTokenList(incoming)
		return
	}

	if expectedLen != expLen {
		assert.Fail(t, fmt.Sprintf(
			"Broken Test: expectedLen of %d does not match expected token list of %d size",
			expectedLen,
			expLen,
		))

		panic("broken test")
		return
	}

	if incLen != expLen {
		assert.Fail(t, fmt.Sprintf(
			"Incoming Tokens (size=%d) does not match incoming expected (size=%d)",
			incLen,
			expLen,
		))

		DebugPrintTokenList(incoming)
		return
	}

	hasError := false

	for pos, token := range incoming {
		if expected[pos] != token.Type {
			log.Printf("Failed: token at pos %d expected to be '%s' but got '%s'\n",
				pos, expected[pos], token.Type)
			hasError = true
		}
	}

	if hasError {
		assert.Fail(t, "One or more tokens failed!")
		DebugPrintTokenList(incoming)
	}
}

func TestBasicTokenize(t *testing.T) {
	assert.Equal(t, 1, 1)

	var (
		str = "hello, world!"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 5, res.tokens, []TokenType{
		TokenWord,    // hello,
		TokenSpace,   //
		TokenWord,    // world
		TokenSentEnd, // !
		TokenEOL,
	})
}

func TestLargeTokenize(t *testing.T) {
	var (
		str = "we cannot stop! main stream control. We cannot stop it. I should've stopped."
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 30, res.tokens, []TokenType{
		TokenWord,    // we
		TokenSpace,   //
		TokenWord,    // cannot
		TokenSpace,   //
		TokenWord,    // stop
		TokenSentEnd, // !
		TokenSpace,   //
		TokenWord,    // main
		TokenSpace,   //
		TokenWord,    // stream
		TokenSpace,   //
		TokenWord,    // control
		TokenPeriod,  // .
		TokenSpace,   //
		TokenWord,    // We
		TokenSpace,   //
		TokenWord,    // cannot
		TokenSpace,   //
		TokenWord,    // stop
		TokenSpace,   //
		TokenWord,    // it
		TokenPeriod,  // .
		TokenSpace,   //
		TokenWord,    // I
		TokenSpace,   //
		TokenWord,    // should've
		TokenSpace,   //
		TokenWord,    // stopped
		TokenPeriod,  // .
		TokenEOL,
	})

}

func TestNumber1(t *testing.T) {
	var (
		str = "well here's about $1,000.00 okay?"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 12, res.tokens, []TokenType{
		TokenWord,    // well
		TokenSpace,   //
		TokenWord,    // here's
		TokenSpace,   //
		TokenWord,    // about
		TokenSpace,   //
		TokenWord,    // $
		TokenNumber,  // 1,000.00
		TokenSpace,   //
		TokenWord,    // okay
		TokenSentEnd, // ?
		TokenEOL,
	})
}

func TestNumber2(t *testing.T) {
	var (
		str = "350 i think"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 6, res.tokens, []TokenType{
		TokenNumber,
		TokenSpace,
		TokenWord,
		TokenSpace,
		TokenWord,
		TokenEOL,
	})

	str = "3.50 i think"
	res = Tokenize(str, VariantSlack)

	assertTokenSequence(t, 6, res.tokens, []TokenType{
		TokenNumber,
		TokenSpace,
		TokenWord,
		TokenSpace,
		TokenWord,
		TokenEOL,
	})
}

func TestSingleQuotes1(t *testing.T) {
	var (
		str = "maybe you should 'get rekt' ok."
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 15, res.tokens, []TokenType{
		TokenWord,        // maybe
		TokenSpace,       //
		TokenWord,        // you
		TokenSpace,       //
		TokenWord,        // should
		TokenSpace,       //
		TokenSingleQuote, // '
		TokenWord,        // get
		TokenSpace,       //
		TokenWord,        // rekt
		TokenSingleQuote, // '
		TokenSpace,       //
		TokenWord,        // ok
		TokenPeriod,      // .
		TokenEOL,
	})
}

func TestSingleQuotes2(t *testing.T) {
	var (
		str = "'scare quotes'"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 6, res.tokens, []TokenType{
		TokenSingleQuote, // '
		TokenWord,        // scare
		TokenSpace,       //
		TokenWord,        // quotes
		TokenSingleQuote, // '
		TokenEOL,
	})
}

func TestSingleQuotes3(t *testing.T) {
	var (
		str = "well... maybe you'll get better 'at it' y'know?"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 20, res.tokens, []TokenType{
		TokenWord,        // well
		TokenEllipsis,    // ...
		TokenSpace,       //
		TokenWord,        // maybe
		TokenSpace,       //
		TokenWord,        // you'll
		TokenSpace,       //
		TokenWord,        // get
		TokenSpace,       //
		TokenWord,        // better
		TokenSpace,       //
		TokenSingleQuote, // '
		TokenWord,        // at
		TokenSpace,       //
		TokenWord,        // it
		TokenSingleQuote, // '
		TokenSpace,       //
		TokenWord,        // y'know
		TokenSentEnd,     // ?
		TokenEOL,
	})
}

func TestEllipsis1(t *testing.T) {
	var (
		str = "..."
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 2, res.tokens, []TokenType{
		TokenEllipsis,
		TokenEOL,
	})
}

func TestEllipsis2(t *testing.T) {
	var (
		str = "well............. idk..... maybe?"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 9, res.tokens, []TokenType{
		TokenWord,     // well
		TokenEllipsis, // .............
		TokenSpace,    //
		TokenWord,     // idk
		TokenEllipsis, // .....
		TokenSpace,    //
		TokenWord,     // maybe
		TokenSentEnd,  // ?
		TokenEOL,
	})
}

func TestBacktickCode1(t *testing.T) {
	var (
		str = "Look: `Here is a list of things`"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 4, res.tokens, []TokenType{
		TokenWord,  // Look:
		TokenSpace, //
		TokenCode,  // `Here is a list of things`
		TokenEOL,   //
	})
}

func TestBacktickCode2(t *testing.T) {
	var (
		str = "Check this out:\n```\nHere\nIs\nSome\nCode\n```\n"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 8, res.tokens, []TokenType{
		TokenWord,    // Check
		TokenSpace,   //
		TokenWord,    // this
		TokenSpace,   //
		TokenWord,    // out:
		TokenNewLine, //
		TokenCode,    // ```\nHere\nIs\nSome\nCode\n```
		// The last newline is trim'd
		TokenEOL,
	})
}

func TestBacktickCode3(t *testing.T) {
	var (
		str = "IDK something like: ```print('foo: %s' % (1))```"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 8, res.tokens, []TokenType{
		TokenWord,  // IDK
		TokenSpace, //
		TokenWord,  // something
		TokenSpace, //
		TokenWord,  // like:
		TokenSpace, //
		TokenCode,  // ```print('foo: %s' % (1))```
		TokenEOL,
	})
}

func TestBacktickCode4(t *testing.T) {
	var (
		str = "Maybe like`this`or```this```maybe?`idk````here```"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 11, res.tokens, []TokenType{
		TokenWord,    // Maybe
		TokenSpace,   //
		TokenWord,    // like
		TokenCode,    // `this`
		TokenWord,    // or
		TokenCode,    // ```this```
		TokenWord,    // maybe
		TokenSentEnd, // ?
		TokenCode,    // `idk`
		TokenCode,    // ```here```
		TokenEOL,
	})
}

func TestBacktickCode5(t *testing.T) {
	var (
		str = "well idk `butt```"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 7, res.tokens, []TokenType{
		TokenWord,  // well
		TokenSpace, //
		TokenWord,  // idk
		TokenSpace, //
		TokenCode,  // `butt`
		TokenCode,  // ``
		TokenEOL,
	})
}

func TestBacktickCode6(t *testing.T) {
	var (
		str = "```"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 3, res.tokens, []TokenType{
		TokenCode, // ``
		TokenWord, // `
		TokenEOL,  //
	})
}

func TestBacktickCode7(t *testing.T) {
	var (
		str = "but if i have:\n```\ntry {\n  await somePromiseCall()\n}\ncatch (error) {"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 10, res.tokens, []TokenType{
		TokenWord,      // but
		TokenSpace,     //
		TokenWord,      // if
		TokenSpace,     //
		TokenWord,      // i
		TokenSpace,     //
		TokenWord,      // have:
		TokenNewLine,   //
		TokenMalformed, // ```etc.
		TokenEOL,
	})
}

func TestBacktickCode8(t *testing.T) {
	var (
		str = "gimme:\n```friendName: friend.name```\n"
		res = Tokenize(str, VariantSlack)
	)

	assert.False(t, res.IsMalformed(), "Part of message is apparently malformed")
	assertTokenSequence(t, 4, res.tokens, []TokenType{
		TokenWord,
		TokenNewLine,
		TokenCode,
		TokenEOL,
	})
}

func TestBacktickCode9(t *testing.T) {
	var (
		str = "i turned\n```\tinput := `blockType {\n    L1 word word\n\tL2 word word word\n\tL3 wor\n\n\tL4 word \n}\n````\ninto\n```(language.BlockDefinition) {\n BlockType: (string) (len=9) \"blockType\",\n KeyValuePairs: ([]language.KeyValuePair) &lt;nil&gt;,\n SubBlocks: ([]language.BlockDefinition) &lt;nil&gt;,\n Words: ([]string) (len=11 cap=16) {\n  (string) (len=2) \"L1\",\n  (string) (len=4) \"word\",\n  (string) (len=4) \"word\",\n  (string) (len=2) \"L2\",\n  (string) (len=4) \"word\",\n  (string) (len=4) \"word\",\n  (string) (len=4) \"word\",\n  (string) (len=2) \"L3\",\n  (string) (len=3) \"wor\",\n  (string) (len=2) \"L4\",\n  (string) (len=4) \"word\"\n }\n}```\n:toot:"
		res = Tokenize(str, VariantSlack)
	)

	assert.False(t, res.IsMalformed(), "Part of message is apparently malformed")
	assertTokenSequence(t, 13, res.tokens, []TokenType{
		TokenWord,     // i
		TokenSpace,    //
		TokenWord,     // turned
		TokenNewLine,  //
		TokenCode,     // codestuff
		TokenWord,     // `
		TokenNewLine,  //
		TokenWord,     // into
		TokenNewLine,  //
		TokenCode,     // morecode
		TokenNewLine,  //
		TokenEmoticon, // toot
		TokenEOL,      //
	})
}

func TestEmoticon1(t *testing.T) {
	var (
		str = "Here's what I think about your dollar: :condi:"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 16, res.tokens, []TokenType{
		TokenWord,     // Here's
		TokenSpace,    //
		TokenWord,     // what
		TokenSpace,    //
		TokenWord,     // I
		TokenSpace,    //
		TokenWord,     // think
		TokenSpace,    //
		TokenWord,     // about
		TokenSpace,    //
		TokenWord,     // your
		TokenSpace,    //
		TokenWord,     // dollar:
		TokenSpace,    //
		TokenEmoticon, // :condi:
		TokenEOL,
	})
}

func TestEmoticon2(t *testing.T) {
	var (
		str = ":condi::condi::condi::condi:"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 5, res.tokens, []TokenType{
		TokenEmoticon,
		TokenEmoticon,
		TokenEmoticon,
		TokenEmoticon,
		TokenEOL,
	})
}

func TestEmoticon3(t *testing.T) {
	var (
		str = ":condi:::condi:"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 4, res.tokens, []TokenType{
		TokenEmoticon,
		TokenWord,
		TokenEmoticon,
		TokenEOL,
	})

	assert.Equal(t, []rune{':'}, res.tokens[1].Value)
}

func TestEmoticon4(t *testing.T) {
	var (
		str = ":condi::::condi:"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 4, res.tokens, []TokenType{
		TokenEmoticon,
		TokenWord,
		TokenEmoticon,
		TokenEOL,
	})

	assert.Equal(t, []rune("::"), res.tokens[1].Value)
}

func TestInfinity1(t *testing.T) {
	var (
		str = "@Dungar how are you?"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 9, res.tokens, []TokenType{
		TokenMentionUser,
		TokenSpace,
		TokenWord,
		TokenSpace,
		TokenWord,
		TokenSpace,
		TokenWord,
		TokenSentEnd,
		TokenEOL,
	})
}

func TestSentence1(t *testing.T) {
	var (
		str = ".@JebBush was terrible on Face The Nation today. Being at 2% and falling seems to have totally affected his confidence. A basket case!"
		res = Tokenize(str, VariantPlain)
	)

	log.Println(str)
	res.DebugPrint()
}

func TestSentence2(t *testing.T) {
	var (
		str = "We need a #POTUS with great strength &amp; stamina. Hillary does not have that. #Trump2016 https://www.facebook.com/DonaldTrump/posts/10156436262810725 http://pbs.twimg.com/media/CWsuCOCW4AE5N-l.jpg"
		res = Tokenize(str, VariantPlain)
	)

	log.Println(str)
	res.DebugPrint()
}

func TestSentence3(t *testing.T) {
	var (
		str = "Will be on MediaBuzz on @FoxNews  at 11:00 A.M."
		res = Tokenize(str, VariantPlain)
	)

	log.Println(str)
	res.DebugPrint()
}
