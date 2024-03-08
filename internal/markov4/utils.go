package markov4

import (
	"encoding/binary"
	"hash/crc64"
	"strings"
)

var repl = strings.NewReplacer(
	"<", "",
	">", "",
	"&lt;", "<",
	"&gt;", ">",
	"&amp;", "\"",
	"_", "",
	"*", "",
	"‘", "'",
	"’", "'",
	"“", "\"",
	"”", "\"",
)

func cleanWordTrimFunc(r rune) bool {
	return r == '\\' ||
			r == ',' ||
			r == ':' ||
			r == '/' ||
			r == '"' ||
			r == ' ' ||
			r == '‘' ||
			r == '’' ||
			r == '“' ||
			r == '”' ||
			r == '\''
}

func cleanWord(w string) string {
	w = repl.Replace(w)
	w = strings.TrimFunc(w, cleanWordTrimFunc)

	return w
}

func reverseTokenIDs(arr []TokenID) []TokenID {
	out := make([]TokenID, len(arr))

	fwd := 0

	for rev := len(arr) - 1; rev >= 0; rev-- {
		out[fwd] = arr[rev]
		fwd++
	}

	return out
}

func joinTokenIDs(listing ...[]TokenID) []TokenID {
	out := make([]TokenID, 0)

	for _, list := range listing {
		out = append(out, list...)
	}

	return out
}

var hashTable *crc64.Table

func hash(f Fragment) uint64 {
	if hashTable == nil {
		hashTable = crc64.MakeTable(crc64.ECMA)
	}

	cr := crc64.New(hashTable)

	binary.Write(cr, binary.LittleEndian, f.LWord)
	binary.Write(cr, binary.LittleEndian, f.CWord)
	binary.Write(cr, binary.LittleEndian, f.RWord)
	//binary.Write(cr, binary.LittleEndian, f.M4Word)

	return cr.Sum64()
}
