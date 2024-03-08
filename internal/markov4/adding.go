package markov4

import "gitlab.int.magneato.site/dungar/prototype/internal/cleaner"

// AddToken will take the incoming string and add it to Words and RevWords
func (m *Markov) AddToken(w cleaner.Token) TokenID {
	uID := w.Hash()
	tID, ok := m.InverseTokens[uID]

	if !ok {
		tID = TokenID(len(m.Tokens))
		m.Tokens = append(m.Tokens, w)
		m.InverseTokens[uID] = tID
		m.TokenWord[string(w.Value)] = tID
	}

	return tID
}

// AddFragment will add the incoming fragment to our Fragments list and
// FragmentWords list and assign it a unique ID
func (m *Markov) AddFragment(f Fragment) FragmentID {
	fID := FragmentID(len(m.Fragments))
	//hID := f.Hash()

	m.Fragments = append(m.Fragments, f)

	// Fragment Token Finding
	_, ok := m.FragmentTokens[f.CWord]

	if !ok {
		m.FragmentTokens[f.CWord] = make([]FragmentID, 0)
	}

	m.FragmentTokens[f.CWord] = append(m.FragmentTokens[f.CWord], fID)

	//// Fragment Prefix Finding
	//_, ok = m.FragmentPrefixes[hID]
	//
	//if !ok {
	//	m.FragmentPrefixes[hID] = make([]FragmentID, 0)
	//}
	//
	//m.FragmentPrefixes[hID] = append(m.FragmentPrefixes[hID], fID)

	return fID
}
