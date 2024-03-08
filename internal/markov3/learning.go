package markov3

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
)

// LearnString takes the incoming string and builds
// out its sequence (via MakeWordSequence) and then
// calls LearnWordSequence
func (m *Markov) LearnString(str string, variant cleaner.TokenVariant) {
	//str = strings.ToLower(str)
	m.LearnTokenList(cleaner.Tokenize(str, variant))
}

// LearnTokenList uses the cleaner.TokenList to build
// a learning sequence. this will also split up sequences
// by cleaner.TokenSentEnd
func (m *Markov) LearnTokenList(tl cleaner.TokenList) {
	var (
		simple   = tl.GetSimpleTokenList()
		sequence = make([]cleaner.SimpleToken, 0)
	)

	for _, token := range simple.Tokens {
		switch token.Type {
		case cleaner.TokenSentEnd:
			m.LearnSimpleTokens(sequence)
			sequence = make([]cleaner.SimpleToken, 0)
		default:
			if m.needsFiller(token.Type) {
				m.AddFiller(token.Type, token.Value)
				token.Value = m.fillerConversion(token.Type)
			}

			sequence = append(sequence, token)
		}
	}

	if len(sequence) > 0 {
		m.LearnSimpleTokens(sequence)
	}
}

// LearnSimpleTokens takes an incoming cleaner.SimpleToken slice
// and just learns it as it would a slice of strings
func (m *Markov) LearnSimpleTokens(tl []cleaner.SimpleToken) {
	wordIDs := make([]MarkovID, len(tl))

	for pos, st := range tl {
		wordID, ok := m.Words[st.Value]

		if !ok {
			wordID = m.AddWord(st.Value)
		}

		wordIDs[pos] = wordID
	}

	chain := m.CreateFragments(wordIDs)

	for _, link := range chain {
		m.AddFragment(link)
	}
}

// LearnWordSequence is our primary means for learning
// new chains and input them into the Markov's various
// data entries
func (m *Markov) LearnWordSequence(words []string) {
	// We're guaranteed to always have the same length
	// of words to word IDs so this lets us plan that.
	wordIDs := make([]MarkovID, len(words))

	for idx, word := range words {
		wordID, ok := m.FindWord(word)

		if !ok {
			wordID = m.AddWord(word)
		}

		//fmt.Printf("word=%s, wordID=%d\n", word, wordID)
		wordIDs[idx] = wordID
	}

	// Now that we have our list of WordIDs, we are going to
	// build a list of fragments and throw them in!
	//log.Printf("creating chain (of wordIDs): %v\n", wordIDs)
	chain := m.CreateFragments(wordIDs)
	//spew.Dump(chain)

	// Fragments for this particular chain are completed, but we're
	// *only* storing them in `chain`, so we want them to go to
	// our `m.Fragments` but also `m.FragmentWords`, so here we go
	for _, link := range chain {
		m.AddFragment(link)
	}
}

// CreateFragments will take the list of wordIDs incoming and
// turn them into a set of tuples:
//
// "for example hello world"
// (LeftMost, "for", "example")
// ("for", "example", "hello")
// ("example", "hello", "world")
// ("hello", "world", RightMost)
func (m *Markov) CreateFragments(wordIDs []MarkovID) []Fragment {
	maxLen := len(wordIDs)
	output := make([]Fragment, 0)

	for pos := 0; pos <= maxLen; pos++ {
		// Even if there are existing fragments that are exactly like this,
		// we still want to create another fragment because it expands the
		// likelihood that this particular Fragment (L/C/R)Words will be selected.
		output = append(output, Fragment{
			LWord: safeGetWordID(wordIDs, maxLen, pos-1),
			CWord: safeGetWordID(wordIDs, maxLen, pos),
			RWord: safeGetWordID(wordIDs, maxLen, pos+1),
		})
	}

	return output
}
