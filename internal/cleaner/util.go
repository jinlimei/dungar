package cleaner

import (
	"hash/crc64"
	"strings"
)

var bulkReplacer = strings.NewReplacer(
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
	w = bulkReplacer.Replace(w)
	w = strings.TrimFunc(w, cleanWordTrimFunc)

	return w
}

func rpad(s string, l int) string {
	start := len(s)

	for ; start < l; start++ {
		s += " "
	}

	return s
}

var hashTable *crc64.Table

func hash(t *Token) uint64 {
	if hashTable == nil {
		hashTable = crc64.MakeTable(crc64.ECMA)
	}

	cr := crc64.New(hashTable)
	cr.Write([]byte{byte(t.Type)})
	cr.Write([]byte(string(t.Value)))
	cr.Write([]byte(string(t.Override)))

	return cr.Sum64()
}
