package cleaner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSpecialSlack1(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<@U9LDWA6QL|dungar>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenMentionUser, prs[0].Type)
	assert.Equal(t, "@U9LDWA6QL", string(prs[0].Value))
	assert.Equal(t, "dungar", string(prs[0].Override))
}

func TestParseSpecialSlack2(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<@U9LDWA6QL>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenMentionUser, prs[0].Type)
	assert.Equal(t, "@U9LDWA6QL", string(prs[0].Value))
	assert.Nil(t, prs[0].Override)

}

func TestParseSpecialSlack3(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<#C9LDWA6QL>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenMentionChannel, prs[0].Type)
	assert.Equal(t, "#C9LDWA6QL", string(prs[0].Value))
	assert.Nil(t, prs[0].Override)
}

func TestParseSpecialSlack4(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<#C9LDWA6QL|butts>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenMentionChannel, prs[0].Type)
	assert.Equal(t, "#C9LDWA6QL", string(prs[0].Value))
	assert.Equal(t, "butts", string(prs[0].Override))
}

func TestParseSpecialSlack5(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<#C9LDWA6QL|butts>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenMentionChannel, prs[0].Type)
	assert.Equal(t, "#C9LDWA6QL", string(prs[0].Value))
	assert.Equal(t, "butts", string(prs[0].Override))
}

func TestParseSpecialSlack6(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<https://www.google.com/|buttoogle yahear>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenURL.String(), prs[0].Type.String())
	assert.Equal(t, "https://www.google.com/", string(prs[0].Value))
	assert.Equal(t, "buttoogle yahear", string(prs[0].Override))
}

func TestParseSpecialSlack7(t *testing.T) {
	var (
		tok = makeToken(TokenSpecial, "<https://www.google.com/>")
		prs = tok.Parse(VariantSlack)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenURL.String(), prs[0].Type.String())
	assert.Equal(t, "https://www.google.com/", string(prs[0].Value))
	assert.Nil(t, prs[0].Override)
}

func TestParseSpecialSlack8(t *testing.T) {
	var (
		tok = Tokenize("Hey <@dungar>, check out: <https://www.google.com|Google URL>", VariantSlack)
	)

	assertTokenSequence(t, 11, tok.tokens, []TokenType{
		TokenWord,
		TokenSpace,
		TokenMentionUser,
		TokenWord,
		TokenSpace,
		TokenWord,
		TokenSpace,
		TokenWord,
		TokenSpace,
		TokenURL,
		TokenEOL,
	})
}

func TestParseSpecialSlack9(t *testing.T) {
	var (
		str = "Here is a complex URL: here<https://www.google.com/|isalink>togo<https://twitter.com|somewhere>" +
			"fun fun fun. <@U9LDWA6QL|dungar> how you doing. check out <#C7ZPDPR8W|programming> for more info"

		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 42, res.tokens, []TokenType{
		TokenWord,           // Here
		TokenSpace,          //
		TokenWord,           // is
		TokenSpace,          //
		TokenWord,           // a
		TokenSpace,          //
		TokenWord,           // complex
		TokenSpace,          //
		TokenWord,           // URL:
		TokenSpace,          //
		TokenWord,           // here
		TokenURL,            // google,isalink
		TokenWord,           // togo
		TokenURL,            // twitter,somewhere
		TokenWord,           // fun
		TokenSpace,          //
		TokenWord,           // fun
		TokenSpace,          //
		TokenWord,           // fun
		TokenPeriod,         // .
		TokenSpace,          //
		TokenMentionUser,    // @U9L,dungar
		TokenSpace,          //
		TokenWord,           // how
		TokenSpace,          //
		TokenWord,           // you
		TokenSpace,          //
		TokenWord,           // doing
		TokenPeriod,         // .
		TokenSpace,          //
		TokenWord,           // check
		TokenSpace,          //
		TokenWord,           // out
		TokenSpace,          //
		TokenMentionChannel, // #C7Z,programming
		TokenSpace,          //
		TokenWord,           // for
		TokenSpace,          //
		TokenWord,           // more
		TokenSpace,          //
		TokenWord,           // info
		TokenEOL,            //
	})
}

func TestParseSpecialSlack10(t *testing.T) {
	var (
		str = "here is|another test: <https://romeosquad.com/#|||#|a text| linke|e>"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 8, res.tokens, []TokenType{
		TokenWord,  // here
		TokenSpace, //
		TokenWord,  // is|another
		TokenSpace, //
		TokenWord,  // test:
		TokenSpace, //
		TokenURL,   // romeosquad.com
		TokenEOL,   //
	})

	assert.Equal(t, "https://romeosquad.com/#", string(res.tokens[6].Value))
	assert.Equal(t, "||#|a text| linke|e", string(res.tokens[6].Override))
}

func TestParseSpecialSlack11(t *testing.T) {
	var (
		str = "here is|another test: <https://romeosquad.com/#|a url with a `lot of text` :separate: from the link!>"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 8, res.tokens, []TokenType{
		TokenWord,  // here
		TokenSpace, //
		TokenWord,  // is|another
		TokenSpace, //
		TokenWord,  // test:
		TokenSpace, //
		TokenURL,   // romeosquad.com
		TokenEOL,   //
	})

	assert.Equal(t, "https://romeosquad.com/#", string(res.tokens[6].Value))
	assert.Equal(t, "a url with a `lot of text` :separate: from the link!", string(res.tokens[6].Override))
}

func TestParseSlack1(t *testing.T) {
	var (
		str = "Hey check out this link at www.google.com/"
		res = Tokenize(str, VariantSlack)
	)

	assertTokenSequence(t, 18, res.tokens, []TokenType{
		TokenWord,   // Hey
		TokenSpace,  //
		TokenWord,   // check
		TokenSpace,  //
		TokenWord,   // out
		TokenSpace,  //
		TokenWord,   // this
		TokenSpace,  //
		TokenWord,   // link
		TokenSpace,  //
		TokenWord,   // at
		TokenSpace,  //
		TokenWord,   // www
		TokenPeriod, // .
		TokenWord,   // google
		TokenPeriod, // .
		TokenWord,   // com/
		TokenEOL,    //
	})
}
