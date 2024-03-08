package learning

import (
	"strings"
)

// CleanToLines takes some common file weirdness and converts things to lines.
func CleanToLines(str string) []string {
	str = cleanupSpace(str)
	str = cleanupQuotes(str)
	str = cleanupLines(str)

	return strings.Split(str, "\n")
}

func cleanupLines(str string) string {
	//log.Printf("[cleanupLines] str='%s'\n", str)

	str = strings.ReplaceAll(str, "\r", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ReplaceAll(str, "\n", " ")
	str = strings.ReplaceAll(str, "  ", " ")
	str = strings.TrimSpace(str)


	return str
}

var stuffReplacer = strings.NewReplacer(
	"\r", "",
	"-\n", "-",
	"\n", " ",
	"‘", "'",
	"’", "'",
	"_", "",
	"“", "",
	"”", "",
	"--", " ",
	"\"", "",
	" . . .", "...",
	". . .", "...",
)

func cleanupQuotes(str string) string {
	return stuffReplacer.Replace(str)
}

func cleanupSpace(str string) string {
	return strings.TrimSpace(str)
}
