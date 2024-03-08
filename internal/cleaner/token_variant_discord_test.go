package cleaner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseSpecialDiscord1(t *testing.T) {
	// here we test <:blobyes:123984019284> emoticons
	var (
		tok = makeToken(TokenSpecial, "\u003c:blobyes:900369292292403211\u003e")
		prs = tok.Parse(VariantDiscord)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenEmoticon, prs[0].Type)
	assert.Equal(t, ":blobyes:", string(prs[0].Value))
	assert.Equal(t, ":blobyes:900369292292403211", string(prs[0].Override))
}

func TestParseSpecialDiscord2(t *testing.T) {
	// here we test @user!
	var (
		tok = makeToken(TokenSpecial, "\u003c@569251658773692538\u003e")
		prs = tok.Parse(VariantDiscord)
	)

	assert.Len(t, prs, 1)
	assert.Equal(t, TokenMentionUser, prs[0].Type)
	assert.Equal(t, "@569251658773692538", string(prs[0].Value))
	assert.Equal(t, "", string(prs[0].Override))
}
