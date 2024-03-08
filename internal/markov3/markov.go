package markov3

import (
	"bytes"
	"encoding/gob"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
)

// MarkovSpaceID is our "unique" id that we can pass to the MakeMarkov
// function to provide a uniqueness to it. String so we have some easier
// identification than its previous uint route.
type MarkovSpaceID string

// FragmentDestination provides our standardized destination handling
// for a particular chain builder
type FragmentDestination string

// MarkovID is our standardized method of handling int changes.
// If we want to, for example, change to int32 or int64 or uint64
// we can simply change this and a few other small spots, to minimize
// overall code impact.
type MarkovID uint32

const (
	// RightSide indicates that a chain generator should go to the right
	RightSide FragmentDestination = "ToRight"
	// LeftSide indicates that a chain generator should go to the left
	LeftSide FragmentDestination = "ToLeft"
)

const (
	// LeftMost is to identify the left-most boundary of a Markov chain
	LeftMost = MarkovID(0)
	// RightMost identifies the right-most boundary of a Markov chain
	RightMost = MarkovID(0)
	// WordNoExist is for word lookups where the lookup failed (it isn't
	// in the Markov's word list)
	WordNoExist = MarkovID(0)
)

// Markov is our internal struct. Keeping it internalized to this package,
// so that to get one you have to use MakeMarkov
type Markov struct {
	// ID is the unique identifier for this particular Markov, only useful if
	// we are working with multiple ones.
	ID MarkovSpaceID

	// Words is our storage list of words, mapped as the word being our key to
	// the ID.
	Words map[string]MarkovID

	// LowerWords is our sub-storage that points to multiple cases of the same
	// word that has been learned
	LowerWords map[string][]MarkovID

	// RevWords is the inverse to Words, and the key of the slice (a number)
	// is our MarkovID
	RevWords []string

	// Fragments is our list of valid fragments that have been generated
	// from learning sentences
	Fragments []*Fragment

	// FragmentWords is a lookup table based on the CWord in the Fragment
	// as the key MarkovID, providing a list of Fragments IDs
	FragmentWords map[MarkovID][]MarkovID

	// MaxChainDistance provides a limitation on how big a chain can be generated.
	MaxChainDistance int

	// MinChainDistance provides a lower-bounds limitation on chains
	// that can be generated.
	MinChainDistance int

	// tokenFillers is our specialized random handling of sub-tokens which
	// we parse out of incoming chat messages. Primarily useful for $URL$, but
	// also handy for everything else.
	tokenFillers map[cleaner.TokenType][]string

	// lastRawMessageID is our database-backed tracking mechanism for the
	// raw messages table. This allows us to track how much we've previously
	// learned and can then jump-start new learning going forward from that point
	lastRawMessageID uint64

	// learnedFromLegacy is declared true when we've learned from raw_messages_m1
	learnedFromLegacy bool
}

// MakeMarkov generates a Markov struct with the necessary basic details that
// can be operated with at a later time.
func MakeMarkov(id MarkovSpaceID) *Markov {
	mt := &Markov{ID: id}

	// We use Reset here because we want a centralized place to handle all
	// data initializing instead of trying to keep this (MakeMarkov) and Reset
	// synced up together.
	mt.Reset()

	return mt
}

// MarkovFromBytes takes the incoming bytes slice and then
// decodes it into a Markov struct. Useful for storing quickly
// generated data.
func MarkovFromBytes(raw []byte) (*Markov, error) {
	var (
		m   *Markov
		err error
		buf = bytes.NewBuffer(raw)
		dec = gob.NewDecoder(buf)
	)

	err = dec.Decode(&m)

	if err != nil {
		return nil, err
	}

	return m, nil
}
