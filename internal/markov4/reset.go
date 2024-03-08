package markov4

import "gitlab.int.magneato.site/dungar/prototype/internal/cleaner"

// Reset will completely remove all data for this markov
// (except its ID) and leave us with an empty state.
func (m *Markov) Reset() {
	m.Tokens = make([]cleaner.Token, 0)
	m.TokenWord = make(map[string]TokenID)
	m.InverseTokens = make(map[uint64]TokenID)
	m.Fragments = make([]Fragment, 0)
	m.FragmentTokens = make(map[TokenID][]FragmentID)
	m.MaxChainDistance = 30
}
