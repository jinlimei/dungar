package cleaner

import (
	"log"
)

// DebugPrintTokenList will take an incoming Token slice and
// use Token.DebugString to print it out
func DebugPrintTokenList(l []Token) {
	for i, t := range l {
		log.Printf("%03d: %v\n", i, t.DebugString())
	}
}

func makeToken(tokType TokenType, value string) Token {
	runes := []rune(value)

	return Token{
		Type:     tokType,
		StartPos: 0,
		EndPos:   len(runes),
		Value:    runes,
		Override: nil,
	}
}

func intAbs(i int) int {
	if i < 0 {
		return i * -1
	}

	return i
}

func isDigit(r rune) bool {
	return r == '0' ||
		r == '1' ||
		r == '2' ||
		r == '3' ||
		r == '4' ||
		r == '5' ||
		r == '6' ||
		r == '7' ||
		r == '8' ||
		r == '9'
}
