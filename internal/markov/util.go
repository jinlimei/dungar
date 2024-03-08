package markov

import "gitlab.int.magneato.site/dungar/prototype/internal/utils"

func normalizeWordList(words []string) []string {
	occurrences := make(map[string]int, 0)
	output := make([]string, 0)

	handleWord := func(w string) int {
		_, ok := occurrences[w]

		if ! ok {
			occurrences[w] = 1
			return 0
		}

		occurrences[w]++
		return 1
	}

	for _, word := range words {
		if handleWord(word) == 0 {
			output = append(output, word)
		}

		normed := utils.Normalize(word)
		if handleWord(normed) == 0 {
			output = append(output, normed)
		}
	}

	return output
}

func makeChains(words []string) [][]string {
	wordLen := len(words)

	wordAt := func(pos int) string {
		if pos < 0 || pos >= wordLen {
			return ""
		}

		return words[pos]
	}

	output := make([][]string, 0)

	for i := 0; i < wordLen; i++ {
		output = append(output, []string{
			wordAt(i - 1),
			words[i],
			wordAt(i + 1),
		})
	}

	return output
}
