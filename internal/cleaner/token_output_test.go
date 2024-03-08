package cleaner

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenList_GetSimpleTokenList1(t *testing.T) {
	assert.Equal(t, 1, 1)

	var (
		str = "Hey <@dungar>, check out: <https://www.google.com|Google URL>"
		res = Tokenize(str, VariantSlack)
		smp = res.GetSimpleTokenList()
	)

	spew.Dump(smp.Tokens)
}
