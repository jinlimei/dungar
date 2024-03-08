package markov4

var store = make(map[MarkovSpaceID]*Markov)

// StoreMarkov stores a markov by its ID
func StoreMarkov(m *Markov) {
	store[m.ID] = m
}

// GetStoredMarkov retrieves a markov by ID
func GetStoredMarkov(id MarkovSpaceID) (*Markov, bool) {
	m, ok := store[id]

	return m, ok
}
