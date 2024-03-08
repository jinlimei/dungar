package cleaner

import "testing"

func zTestParseBasic1(t *testing.T) {
	var (
		str = "Hey check out this link at https://www.google.com/"
		res = Tokenize(str, VariantPlain)
	)

	res.DebugPrint()
}
