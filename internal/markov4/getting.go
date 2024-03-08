package markov4

import "gitlab.int.magneato.site/dungar/prototype/internal/cleaner"

// GetRandomWordID provides a random word-associated
// MarkovID from the RevWords listing
func (m *Markov) GetRandomWordID() MarkovID {
	return randomMarkovID(MarkovID(len(m.Tokens)))
}

// GetRandomWordStr takes in GetRandomWordID's value
// and returns the actual string word
func (m *Markov) GetRandomWordStr() cleaner.Token {
	id := m.GetRandomWordID()

	return m.Tokens[id]
}

// GetSimpleWordList will take all the TokenWord's and then
// build a slice of them as a string
func (m *Markov) GetSimpleWordList() []string {
	listing := make([]string, len(m.TokenWord))

	pos := 0
	for word := range m.TokenWord {
		listing[pos] = word
		pos++
	}

	return listing
}
