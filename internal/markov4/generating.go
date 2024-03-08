package markov4

import (
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func isInterchangeable(t cleaner.TokenType) bool {
	return t == cleaner.TokenURL ||
		t == cleaner.TokenHashTag ||
		t == cleaner.TokenMentionChannel ||
		t == cleaner.TokenMentionUser ||
		t == cleaner.TokenMentionRole ||
		t == cleaner.TokenEmoticon ||
		// Let's just interchange numbers and not care
		t == cleaner.TokenNumber
}

// Generate provides a bit simpler methodology for
// building a graph and turning it into a useful sentence.
func (m *Markov) Generate(w string) string {
	token, ok := m.TokenWord[w]

	if !ok {
		token, ok = m.TokenWord[cleanWord(w)]
	}

	if !ok {
		return ""
	}

	tokens := m.GenerateGraph(token)
	outs := m.TokensToStrings(tokens)

	return strings.TrimSpace(strings.Join(outs, " "))
}

// GenerateNextWord will take two incoming words and provide a list of the following
// word with a count of how many times it was encountered, and a total count
func (m *Markov) GenerateNextWord(w1, w2 string) (map[string]int, int) {
	t1, ok1 := m.TokenWord[cleanWord(w1)]
	t2, ok2 := m.TokenWord[cleanWord(w2)]

	if !ok1 || !ok2 {
		return nil, 0
	}

	var (
		nexts  = m.GetNextTokenIDs(t1, t2)
		total  = 0
		output = make(map[string]int, 0)
	)

	for _, count := range nexts {
		total += count
	}

	for tokenID, count := range nexts {
		token := m.Tokens[tokenID]

		output[token.ValueOrOverride()] = count
	}

	return output, total
}

// TokensToStrings takes an incoming slice of TokenID and finds the
// string value for the tokens.
func (m *Markov) TokensToStrings(tokens []TokenID) []string {
	output := make([]string, len(tokens))

	for idx, tID := range tokens {
		if tID < 0 {
			continue
		}

		token := m.Tokens[tID]

		output[idx] = token.ValueOrOverride()
	}

	return output
}

func (m *Markov) isSpaceToken(w TokenID) bool {
	token := m.Tokens[w]
	return token.Type == cleaner.TokenSpace || token.Type == cleaner.TokenNewLine
}

// GetNextTokenIDs is same as GenerateNextWord but strictly for TokenID's
// and will not return totals
func (m *Markov) GetNextTokenIDs(w1, w2 TokenID) map[TokenID]int {
	if len(m.Tokens) == 0 || len(m.Fragments) == 0 {
		return nil
	}

	frags, ok := m.FindLCFragments(w1, w2)

	if !ok {
		return nil
	}

	stats := make(map[TokenID]int, len(frags))

	for _, frag := range frags {
		_, ok = stats[frag.RWord]

		if !ok {
			stats[frag.RWord] = 1
		} else {
			stats[frag.RWord]++
		}
	}

	return stats
}

// GenerateGraph takes an input word and then builds a Markov chain for that
// particular word. The chain is built with the word being relatively
// centered, so we generate that from center-to-left and center-to-right.
func (m *Markov) GenerateGraph(w TokenID) []TokenID {
	if len(m.Tokens) == 0 || len(m.Fragments) == 0 {
		return []TokenID{}
	}

	f, ok := m.FindFragment(w)
	if !ok {
		f, ok = m.PickRandomFragment()

		if !ok {
			panic("Something went horribly wrong")
		}
	}

	var (
		left   = m.FollowChain(f, LeftSide)
		center = f.Tokens()
		right  = m.FollowChain(f, RightSide)
	)

	return joinTokenIDs(left, center, right)
}

// FollowChain lets us determine our source & target WordTypes based off of our
// incoming destination, and then builds an entire chain (through NextChainLink)
// up to the MaxChainDistance setup for Markov
func (m *Markov) FollowChain(base Fragment, dest FragmentDestination) []TokenID {
	var (
		ok      bool
		source  WordType
		target  WordType
		results = make([]TokenID, 0)
		pos     = 0
		frag    = base
	)

	switch dest {
	case LeftSide:
		source = RWord
		target = LWord
	case RightSide:
		source = LWord
		target = RWord
	}

	for ; pos <= m.MaxChainDistance; pos++ {
		frag, ok = m.NextChainLink(base, frag, source, target)

		if !ok || frag.Equals(base) {
			break
		}

		results = append(results, frag.GetToken(target))
	}

	if target == LWord {
		return reverseTokenIDs(results)
	}

	return results
}

// NextChainLink provides, given the source and target WordTypes, the next "best"
// (or available) fragment for building a chain. Chains in this particular instance
// are picked at random.
func (m *Markov) NextChainLink(orig, curr Fragment, source, target WordType) (Fragment, bool) {
	var (
		centered []FragmentID
		ok       bool
		eligible []Fragment
		decided  Fragment
	)

	//log.Printf("NextChainLink: fragment=%s, source=%v, target=%v\n",
	//	curr.String(), source, target)

	if curr.GetToken(target) == LeftMost || curr.GetToken(target) == RightMost {
		//log.Printf("CWord is LeftMost or RightMost, returning not-ok\n")
		return Fragment{}, false
	}

	centered, ok = m.FragmentTokens[curr.GetToken(target)]
	//log.Printf("wordFragmentIDs for word %d from target %s: (ok? %v) %v\n",
	//	curr.GetToken(target), target, ok, centered)

	if !ok {
		return Fragment{}, false
	}

	for _, fragID := range centered {
		frag := m.Fragments[fragID]
		//log.Printf("Checking fragment %v (id: %d) for Word %d at source %s (which is %d)\n",
		//	frag, fragID, curr.CWord, source, frag.GetToken(source))

		if frag.Equals(orig) {
			//log.Printf("Fragment (%s) matches original (%s)\n",
			//	frag.String(), orig.String())
			continue
		}

		if frag.GetToken(source) == curr.CWord {
			//log.Printf("Eligible: %s\n", frag)
			eligible = append(eligible, frag)
		}
	}

	if len(eligible) == 0 {
		return Fragment{}, false
	} else if len(eligible) == 1 {
		decided = eligible[0]
	} else {
		decided = eligible[random.Int(len(eligible))]
	}

	return decided, true

	/**
	Starting from 'centered' and going left:
	 - Pull up L1Word's since that is our first target.
	 - Find the right L1Word to work with
	 - Pull up L2Word's since that is our next target.
	 - Find the right L2Word in the listing
	*/

	//return Fragment{}, false
}

// PeekEligibility lets us get a count of the next available
// valid fragments. Useful for determining if we should pick
// a random fragment instead of an eligible fragment.
func (m *Markov) PeekEligibility(f Fragment, source, target WordType) int {
	return 0
}
