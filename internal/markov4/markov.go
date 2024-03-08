package markov4

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
type FragmentDestination uint8

// MarkovID is our standardized method of handling int changes.
// If we want to, for example, change to int32 or int64 or uint64
// we can simply change this and a few other small spots, to minimize
// overall code impact.
type MarkovID int64

// TokenID is our MarkovID for tokens (words)
type TokenID MarkovID

// FragmentID is our unique ID for fragments
type FragmentID uint32

const (
	// LeftSide indicates that a chain generator should go to the left
	LeftSide FragmentDestination = iota
	// RightSide indicates that a chain generator should go to the right
	RightSide
)

const (
	// LeftMost is to identify the left-most boundary of a Markov chain
	LeftMost = TokenID(-100)
	// RightMost identifies the right-most boundary of a Markov chain
	RightMost = TokenID(-200)
	// WordNoExist is for word lookups where the lookup failed (it isn't
	// in the Markov's word list)
	WordNoExist = TokenID(-1)
)

// Markov is our internal struct. Keeping it internalized to this package,
// so that to get one you have to use MakeMarkov
type Markov struct {
	// ID is the unique identifier for this particular Markov, only useful if
	// we are working with multiple ones.
	ID MarkovSpaceID

	Tokens []cleaner.Token

	TokenWord map[string]TokenID

	InverseTokens map[uint64]TokenID

	Fragments []Fragment

	FragmentTokens map[TokenID][]FragmentID

	MaxChainDistance int
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
