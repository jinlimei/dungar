package markov3

import (
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

// FindWordsWithVariants will take the incoming words and return their given variants via
// FindWordVariants
func (m *Markov) FindWordsWithVariants(w string) map[string][]MarkovID {
	var (
		words = strings.Split(w, " ")
		out   = make(map[string][]MarkovID, len(words))
	)

	for _, word := range words {
		out[word] = m.FindWordVariants(word)
	}

	return out
}

// FindWordVariants will attempt to look at all the variants (casing, etc.) of
// a word.
func (m *Markov) FindWordVariants(w string) []MarkovID {
	w = strings.ToLower(w)
	cW := cleanWord(w)

	if w == cW {
		return m.LowerWords[w]
	}

	ids1, _ := m.LowerWords[w]
	ids2, _ := m.LowerWords[cW]

	out := make([]MarkovID, 0, len(ids1)+len(ids2))
	out = append(out, ids1...)
	out = append(out, ids2...)

	return out
}

// FindWord will attempt to clean the incoming word
// and return the MarkovID associated with it (if
// present, otherwise the MarkovID will be WordNoExist)
func (m *Markov) FindWord(w string) (MarkovID, bool) {
	// lower-case lets us normalize the output
	// in such a manner that we can provide a better
	// Markov experience
	id, ok := m.Words[w]

	if !ok {
		w = cleanWord(w)
	}

	id, ok = m.Words[w]

	if !ok {
		return WordNoExist, false
	}

	return id, true
}

// GetFragmentWithWordIDs will attempt to retrieve a fragment with n MarkovID's, from
// 1 to 3.
// For single word IDs it will use GetFragmentWithWordID
// For two word IDs it will use findLCFragments and findCRFragments
// For three word IDs it will use findLCRFragments
func (m *Markov) GetFragmentWithWordIDs(ids []MarkovID) (Fragment, bool) {
	idLen := len(ids)

	if idLen == 1 {
		return m.GetFragmentWithWordID(ids[0])
	}

	// We can only establish a fragment with 2 or 3 words (LWord, CWord, RWord)
	if idLen != 2 && idLen != 3 {
		return Fragment{}, false
	}

	var choices []Fragment

	switch idLen {
	case 1:
		return m.GetFragmentWithWordID(ids[0])
	case 2:
		f1, _ := m.findLCFragments(ids[0], ids[1])
		f2, _ := m.findCRFragments(ids[0], ids[1])

		choices = make([]Fragment, 0, len(f1)+len(f2))
		choices = append(choices, f1...)
		choices = append(choices, f2...)

	case 3:
		choices, _ = m.findLCRFragments(ids[0], ids[1], ids[2])

	default:
		return Fragment{}, false
	}

	return choices[random.Int(len(choices))], true
}

func (m *Markov) findLCFragments(lID, cID MarkovID) ([]Fragment, bool) {
	ids, ok := m.FragmentWords[cID]

	if !ok {
		return []Fragment{}, false
	}

	found := make([]Fragment, 0, len(ids))
	for _, id := range ids {
		frag := *m.Fragments[id]

		if frag.LWord == lID {
			found = append(found, frag)
		}
	}

	return found, true
}

func (m *Markov) findCRFragments(cID, rID MarkovID) ([]Fragment, bool) {
	ids, ok := m.FragmentWords[cID]

	if !ok {
		return []Fragment{}, false
	}

	found := make([]Fragment, 0, len(ids))
	for _, id := range ids {
		frag := *m.Fragments[id]

		if frag.RWord == rID {
			found = append(found, frag)
		}
	}

	return found, true
}

func (m *Markov) findLCRFragments(lID, cID, rID MarkovID) ([]Fragment, bool) {
	ids, ok := m.FragmentWords[cID]

	if !ok {
		return []Fragment{}, false
	}

	found := make([]Fragment, 0, len(ids))
	for _, id := range ids {
		frag := *m.Fragments[id]

		if frag.LWord == lID && frag.RWord == rID {
			found = append(found, frag)
		}
	}

	return found, true
}

// WordFragmentCount will look at the word w
// and see how many fragments have that word
// as their CWord
func (m *Markov) WordFragmentCount(w MarkovID) int {
	fragments, ok := m.FragmentWords[w]
	if !ok {
		return 0
	}

	return len(fragments)
}

// WordFragmentIDs will look at the word w
// for fragments, and return that list of
// fragments.
func (m *Markov) WordFragmentIDs(w MarkovID) ([]MarkovID, bool) {
	if w < 0 {
		return []MarkovID{}, false
	}

	fragments, ok := m.FragmentWords[w]
	//log.Printf("WordFragmentIDs(w=%d), ok=%v, fragments=%v\n", w, ok, fragments)

	if !ok {
		//log.Fatalf("WordFragmentIDs failed to capture any with word %d\n", w)
		return []MarkovID{}, false
	}

	return fragments, true
}

// GetFragmentWithWordID will look for a particular
// fragment associated with that word MarkovID.
// if none can be found we pick a word from
// random and return that instead.
func (m *Markov) GetFragmentWithWordID(w MarkovID) (Fragment, bool) {
	fragIDs, ok := m.FragmentWords[w]

	if ok {
		fragID := fragIDs[random.Int(len(fragIDs))]
		return *m.Fragments[fragID], true
	}

	// we failed to find a fragment for the word provided (perhaps
	// it does not exist or any of the fragments were pulled? who
	// knows, really)

	iter := 0
	revLen := MarkovID(len(m.RevWords))

	for !ok {
		// Pick a word from random:
		w = randomMarkovID(revLen)
		fragIDs, ok = m.FragmentWords[w]
		iter++

		if !ok && iter > 10 {
			log.Printf("ERROR: Iteration Count for finding a word is really shit, fix please")
			return Fragment{}, false
		}
	}

	fragID := fragIDs[random.Int(len(fragIDs))]
	return *m.Fragments[fragID], true
}

// WordIDsToWords takes all incoming MarkovIDs and converts them
// to their corresponding string values, including LeftMost and
// RightMost
func (m *Markov) WordIDsToWords(wordIDs []MarkovID) []string {
	//log.Printf("WordIDsToWords: %v\n", wordIDs)
	var (
		words = make([]string, len(wordIDs))
		pos   = 0

		w  string
		t  cleaner.TokenType
		ok bool
	)

	for _, wordID := range wordIDs {
		if wordID < 0 {
			continue
		}

		w = m.RevWords[wordID]
		t, ok = revFillerMap[w]

		if ok {
			w = random.PickString(m.tokenFillers[t])
		}

		words[pos] = w
		pos++
	}

	return words
}
