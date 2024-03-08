package learning

import "strings"

// ReadABible reads the bible (formatted a special way)
func ReadABible(str string) []string {
	var (
		split = strings.Split(str, "\n\n")
		out   = make([]string, 0, len(split))
	)

	for _, s := range split {
		s = strings.TrimSpace(s)
		s = strings.ReplaceAll(s, "\n", " ")

		if s == "" {
			continue
		}

		out = append(out, s)
	}

	return out
}

func peekUntil(rs []rune, r rune, s int) int {
	rl := len(rs)

	if s < 0 || s+1 >= rl {
		return -1
	}

	for ; s < rl; s++ {
		if rs[s] == r {
			return s
		}
	}

	return -1
}

func isNumber(x rune) bool {
	return x == '0' || x == '1' || x == '2' || x == '3' ||
		x == '4' || x == '5' || x == '6' || x == '7' ||
		x == '8' || x == '9' || x == ':'
}
