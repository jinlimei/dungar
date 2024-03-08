package cleaner

import "testing"

func TestParseSpecialTwitter1(t *testing.T) {
	var (
		str = "RT @bob: Here's what it looks like to eat bananas: https://www.google.com/"
		res = Tokenize(str, VariantTwitter)
	)

	res.DebugPrint()
}
