package markov3

import (
	"strings"
)

// AddWord will take the incoming string and add it to Words and RevWords
func (m *Markov) AddWord(w string) MarkovID {
	var (
		lower = strings.ToLower(w)
		mID   = MarkovID(len(m.RevWords))
	)

	m.Words[w] = mID
	m.RevWords = append(m.RevWords, w)

	if _, ok := m.LowerWords[lower]; !ok {
		m.LowerWords[lower] = []MarkovID{mID}
	} else {
		m.LowerWords[lower] = append(m.LowerWords[lower], mID)
	}

	//log.Printf("AddWord(%s); WordID=%d, NextWordID=%d\n", w, m.Words[w], len(m.RevWords))

	return mID
}

// AddFragment will add the incoming fragment to our Fragments list and
// FragmentWords list and assign it a unique ID
func (m *Markov) AddFragment(f Fragment) MarkovID {
	// Get the FragmentID we'll be using for this particular fragment
	fragID := MarkovID(len(m.Fragments))

	// Simple: add to fragments
	m.Fragments = append(m.Fragments, &f)

	// Not-so-simple: We want a reverse lookup (based on the CWord of the
	// fragment, are there other fragments?) This lets us avoid an entire
	// search of `m.Fragments` when wanting to find `f.CWord`
	_, ok := m.FragmentWords[f.CWord]

	// This is the first time `f.CWord` was used, so we build its slice so we
	// can append to it easily
	if !ok {
		m.FragmentWords[f.CWord] = make([]MarkovID, 0)
	}

	// Now we add our fragID to the list of fragments that have f.CWord
	m.FragmentWords[f.CWord] = append(m.FragmentWords[f.CWord], fragID)

	return fragID
}
