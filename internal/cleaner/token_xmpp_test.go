package cleaner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXMPP1(t *testing.T) {
	assert.Equal(t, 1, 1)

	var (
		msg = "(20:46:56) ***Dungarmatic readies vuvuzela"
		res = Tokenize(msg, VariantXMPP)
	)

	res.DebugPrint()
}

func TestXMPP2(t *testing.T) {
	assert.Equal(t, 1, 1)

	var (
		msg = "Dungarmatic: buttstuf​f linked that 43 days ago. It has been relinked 2 times."
		res = Tokenize(msg, VariantXMPP)
	)

	res.DebugPrint()
}

func TestXMPP3(t *testing.T) {
	assert.Equal(t, 1, 1)

	var (
		msg = "<Madakad> https://www.youtube.com/watch?v=IcrbM1l_BoIembed: Avicii - Wake Me Up (Official Video) - Get the new EP here: https://lnk.to/AviciEPAvicii - https://www.youtube.com/watch?v=IcrbM1l_BoI"
		res = Tokenize(msg, VariantXMPP)
	)

	res.DebugPrint()
}
