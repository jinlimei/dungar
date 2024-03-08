package markov3

import (
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

// Generate provides a bit simpler methodology for
// building a graph and turning it into a useful sentence.
func (m *Markov) Generate(w string) string {
	if w == "" || strings.TrimSpace(w) == "" {
		return ""
	}

	wID, ok := m.FindWord(w)
	if !ok {
		// MarkovID is generally needed but this will
		// have to suffice.
		wID = randomMarkovID(MarkovID(len(m.RevWords)))
	}

	resIDs := m.GenerateGraph(wID)
	words := m.WordIDsToWords(resIDs)
	return strings.TrimSpace(strings.Join(words, " "))
}

// GenerateGraphFromFragment will generate a full result from the provided Fragment.
func (m *Markov) GenerateGraphFromFragment(frag Fragment) []MarkovID {
	if len(m.Words) == 0 || len(m.Fragments) == 0 {
		return []MarkovID{}
	}

	var (
		left   = m.FollowChain(frag, LeftSide)
		center = frag.Words()
		right  = m.FollowChain(frag, RightSide)
	)

	return joinMarkovIDs(left, center, right)
}

// GenerateGraph takes an input word and then builds a Markov chain for that
// particular word. The chain is built with the word being relatively
// centered, so we generate that from center-to-left and center-to-right.
func (m *Markov) GenerateGraph(w MarkovID) []MarkovID {
	if len(m.Words) == 0 || len(m.Fragments) == 0 {
		return []MarkovID{}
	}

	frag, ok := m.GetFragmentWithWordID(w)
	if !ok {
		log.Printf("Something horribly failed")
		return []MarkovID{}
	}

	//log.Printf("GenerateGraph(w=%d) fragment decided=%v\n", w, frag.String())
	return m.GenerateGraphFromFragment(frag)
}

// FollowChain lets us determine our source & target WordTypes based off of our
// incoming destination, and then builds an entire chain (through NextChainLink)
// up to the MaxChainDistance setup for Markov
func (m *Markov) FollowChain(base Fragment, dest FragmentDestination) []MarkovID {
	var ok bool

	//log.Printf("\n\nFollowChain(base=%v), dest=%v\n\n", base.String(), dest)

	// ToRight -> source=LWord target=RWord
	// ToLeft -> source=RWord target=LWord
	source, target := m.SplitDirections(dest)

	const limit = 20

	var (
		frag    = base
		results = make([]MarkovID, 0)
		link    = 0
		runt    = 0
	)

	for link < m.MinChainDistance && runt < limit {
		for ; link <= m.MaxChainDistance; link++ {
			frag, ok = m.NextChainLink(base, frag, source, target)

			if !ok || frag.Equals(base) {
				break
			}

			results = append(results, frag.GetWord(target))
		}

		runt++
	}

	if target == LWord {
		return reverseMarkovIDs(results)
	}

	return results
}

// NextChainLink provides, given the source and target WordTypes, the next "best"
// (or available) fragment for building a chain. Chains in this particular instance
// are picked at random.
func (m *Markov) NextChainLink(orig, cur Fragment, source, target WordType) (Fragment, bool) {
	var centered []MarkovID
	var ok bool

	//log.Printf("NextChainLink: fragment=%s, source=%v, target=%v\n",
	//	cur.String(), source, target)

	if cur.GetWord(target) == LeftMost || cur.GetWord(target) == RightMost {
		//log.Printf("CWord is LeftMost or RightMost, returning not-ok\n")
		return Fragment{}, false
	}

	centered, ok = m.WordFragmentIDs(cur.GetWord(target))
	//log.Printf("wordFragmentIDs for word %d from target %s: (ok? %v) %v\n", cur.GetWord(target), target, ok, centered)

	if !ok {
		return Fragment{}, false
	}

	eligible := make([]*Fragment, 0)

	//log.Printf("source: %s, target: %s", source, target)
	//log.Printf("cur: %+v", cur.Words())

	for _, fragID := range centered {
		frag := m.Fragments[fragID]
		//log.Printf("Checking fragment %v (id: %d) for Word %d at source %s (which is %d)\n",
		//	frag, fragID, cur.CWord, source, frag.GetWord(source))
		if frag.Equals(orig) {
			//log.Printf("Fragment (%s) matches original (%s)\n",
			//	frag.String(), orig.String())
			continue
		}

		// source = LWord or RWord
		// If the "next fragment" we're looking for is
		// - in the right direction then the source should be LWord
		// cur  = [LWord, CWord, RWord]
		// frag =        [LWord, CWord, RWord]
		//                       ^^^^^
		// - in the left direction then the source should be RWord
		// cur  =        [LWord, CWord, RWord]
		// frag = [LWord, CWord, RWord]
		//                       ^^^^^

		if frag.GetWord(source) == cur.CWord {
			//log.Printf("Found Eligible Fragment: %s\n", frag.String())
			eligible = append(eligible, frag)
		}
	}

	totalEligible := len(eligible)

	var next Fragment

	randomCentered := centered[random.Int(len(centered))]

	if totalEligible == 0 {
		//log.Printf("NextChainLink could not find another eligible candidate, picking one at random!\n")
		//next = *m.Fragments[randomCentered]
		return Fragment{}, false
	} else if totalEligible == 1 && m.PeekEligibility(*eligible[0], source, target) == 1 {
		// We peaked at the eligibility count of the next link, and found it also to be one,
		// so we're going to diverge intentionally
		next = *m.Fragments[randomCentered]
	} else {
		next = *eligible[random.Int(len(eligible))]
	}

	return next, true
}

// PeekEligibility lets us get a count of the next available
// valid fragments. Useful for determining if we should pick
// a random fragment instead of an eligible fragment.
func (m *Markov) PeekEligibility(f Fragment, source, target WordType) int {
	next, ok := m.WordFragmentIDs(f.GetWord(target))

	if !ok {
		return 0
	}

	avail := 0
	for _, id := range next {
		if m.Fragments[id].GetWord(source) == f.CWord {
			avail++
		}
	}

	return avail
}

// SplitDirections lets us specify which WordType directions we should go
// for our destination.
func (m *Markov) SplitDirections(dest FragmentDestination) (WordType, WordType) {
	switch dest {
	case RightSide:
		return LWord, RWord
	case LeftSide:
		return RWord, LWord
	}

	panic("at the disco")
}
