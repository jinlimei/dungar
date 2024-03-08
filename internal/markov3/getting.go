package markov3

// GetRandomWordID provides a random word-associated
// MarkovID from the RevWords listing
func (m *Markov) GetRandomWordID() MarkovID {
	return randomMarkovID( MarkovID(len(m.RevWords)) )
}

// GetRandomWordStr takes in GetRandomWordID's value
// and returns the actual string word
func (m *Markov) GetRandomWordStr() string {
	id := m.GetRandomWordID()

	return m.RevWords[id]
}
