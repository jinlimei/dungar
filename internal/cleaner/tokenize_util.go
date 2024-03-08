package cleaner

func peek(set []rune, idx int) rune {
	if idx < 0 {
		return nilRune
	}

	if idx >= len(set) {
		return nilRune
	}

	return set[idx]
}

func peeped(v rune, r ...rune) bool {
	var (
		rLen = len(r)
		pos  = 0
	)

rePeep:
	if r[pos] == v {
		return true
	}

	pos++
	if pos < rLen {
		goto rePeep
	}

	return false
}

func getRuneRange(set []rune, start, end int) []rune {
	setLen := len(set)

	if start < 0 || start > setLen {
		return nil
	}

	if start >= end || end < 0 || end > setLen {
		return nil
	}

	return set[start:end]
}

func inlinePeekUntil(set []rune, start int, search rune) int {
	setLen := len(set)
	if start < 0 {
		return -1
	}

	if start >= setLen {
		return -1
	}

	pos := start

ret:
	if set[pos] == search {
		return pos
	}

	pos++
	if pos < setLen {
		goto ret
	}

	return -1
}

func peekUntil(set []rune, start int, search rune) int {
	setLen := len(set)

	if start < 0 {
		return -1
	}

	if start >= setLen {
		return -1
	}

	for pos := start; pos < setLen; pos++ {
		if set[pos] == search {
			return pos
		}
	}

	return -1
}

func peekUntilSet(set []rune, start int, searchSet ...rune) (int, rune) {
	setLen := len(set)
	searchSetLen := len(searchSet)

	if start < 0 {
		return -1, nilRune
	}

	if start >= setLen {
		return -1, nilRune
	}

	for pos := start; pos < setLen; pos++ {
		if set[pos] == searchSet[0] {
			return pos, set[pos]
		}

		for k := 1; k < searchSetLen; k++ {
			if set[pos] == searchSet[k] {
				return pos, set[pos]
			}
		}
	}

	return -1, nilRune
}
