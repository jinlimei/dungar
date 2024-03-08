package markov4

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

// FindToken will attempt to clean the incoming word
// and return the MarkovID associated with it (if
// present, otherwise the MarkovID will be WordNoExist)
func (m *Markov) FindToken(w cleaner.Token) (TokenID, bool) {
	key := w.Hash()
	tID, ok := m.InverseTokens[key]

	if !ok {
		return WordNoExist, false
	}

	return tID, true
}

// WordFragmentCount will look at the word w
// and see how many fragments have that word
// as their CWord
func (m *Markov) WordFragmentCount(w MarkovID) int {
	return 0
}

// PickRandomFragment will find a random Fragment and return true
// or give up and return false
func (m *Markov) PickRandomFragment() (Fragment, bool) {
	fLen := len(m.Fragments)

	if fLen == 0 {
		return Fragment{}, false
	}

	fID := TokenID(random.Int(fLen))
	return m.Fragments[fID], true
}

// FindLCFragments will take the incoming LToken and CToken TokenID's and return a list
// of Fragment's for each fragment which matched these two TokenID's
func (m *Markov) FindLCFragments(LToken, CToken TokenID) ([]Fragment, bool) {
	ids, ok := m.FragmentTokens[CToken]

	if !ok {
		return []Fragment{}, false
	}

	found := make([]Fragment, 0, len(ids))
	for _, id := range ids {
		frag := m.Fragments[id]

		if frag.LWord == LToken {
			found = append(found, frag)
		}
	}

	return found, true
}

// FindFragment will look for a particular
// fragment associated with that word TokenID.
func (m *Markov) FindFragment(w TokenID) (Fragment, bool) {
	ids, ok := m.FragmentTokens[w]

	if ok {
		id := ids[random.Int(len(ids))]
		return m.Fragments[id], true
	}

	return Fragment{}, false
}

// WordIDsToWords takes all incoming MarkovIDs and converts them
// to their corresponding string values, including LeftMost and
// RightMost
func (m *Markov) WordIDsToWords(wordIDs []MarkovID) []cleaner.Token {
	return []cleaner.Token{}
}
