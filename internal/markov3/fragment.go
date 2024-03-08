package markov3

import "fmt"

// WordType is our type for
// the difference pieces of a Fragment
type WordType string

const (
	// LWord -> Fragment.LWord
	LWord WordType = "LWord"
	// CWord -> Fragment.CWord
	CWord WordType = "CWord"
	// RWord -> Fragment.RWord
	RWord WordType = "RWord"
)

// Opposite provides the "opposite" of the given WordType (so LWord for RWord, RWord for LWord)
func (wt WordType) Opposite() WordType {
	switch wt {
	case LWord:
		return RWord
	case CWord:
		return CWord
	case RWord:
		return LWord
	default:
		panic("Invalid WordType")
	}
}

// String coerces WordType to a string
func (wt WordType) String() string {
	return string(wt)
}

// Fragment provides our tuple for
// building up markov graphs
type Fragment struct {
	// LWord is our left-most word
	LWord MarkovID
	// CWord is our centered word
	CWord MarkovID
	// RWord is our right-most word
	RWord MarkovID
}

// String outputs a debug string for a given Fragment
func (f Fragment) String() string {
	return fmt.Sprintf("[L %d,C %d,R %d]", f.LWord, f.CWord, f.RWord)
}

// Equals lets us compare fragments (all words are compared)
func (f Fragment) Equals(o Fragment) bool {
	return f.CWord == o.CWord &&
		f.LWord == o.LWord &&
		f.RWord == o.RWord
}

// GetWord takes the incoming WordType and provides the
// value of the entry associated with it.
func (f Fragment) GetWord(t WordType) MarkovID {
	switch t {
	case LWord:
		return f.LWord
	case CWord:
		return f.CWord
	case RWord:
		return f.RWord
	}

	panic("at the disco")
}

// Words outputs the different WordType in
// an ordered slice (LWord, CWord, RWord)
func (f Fragment) Words() []MarkovID {
	return []MarkovID{
		f.LWord,
		f.CWord,
		f.RWord,
	}
}
