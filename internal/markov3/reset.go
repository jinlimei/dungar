package markov3

// Reset will completely remove all data for this markov
// (except its ID) and leave us with an empty state.
func (m *Markov) Reset() {
	m.lastRawMessageID = 0
	m.Words = make(map[string]MarkovID)
	m.LowerWords = make(map[string][]MarkovID)
	m.RevWords = make([]string, 0)
	m.Fragments = make([]*Fragment, 0)
	m.FragmentWords = make(map[MarkovID][]MarkovID)
	m.MaxChainDistance = 30
	m.MinChainDistance = 2

	// We must have 0 captured in the Markov word tables, otherwise
	// WordNoExist will actually point to a word.
	m.AddWord("")

	// Provide initial fillers so we don't run on empty. That would be bad.
	m.initInitialFillers()
}
