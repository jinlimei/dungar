package markov4

import "fmt"

// WordType is our type for
// the difference pieces of a Fragment
type WordType uint8

const (
	// LWord -> Fragment.LWord
	LWord WordType = iota
	// CWord -> Fragment.CWord
	CWord
	// RWord -> Fragment.RWord
	RWord
)

// String converts WordType into a string value
func (wt WordType) String() string {
	switch wt {
	case LWord: return "LWord"
	case CWord: return "CWord"
	case RWord: return "RWord"
	}

	return ""
}

// Fragment provides our tuple for
// building up markov graphs
type Fragment struct {
	// LWord is our left-most word
	LWord TokenID
	// CWord is our center word
	CWord TokenID
	// RWord is our right-most word
	RWord TokenID
	hash  *uint64
}

// String outputs a debug string for a given Fragment
func (f Fragment) String() string {
	return fmt.Sprintf("[L %d,C %d,R %d]", f.LWord, f.CWord, f.RWord)
}

// Hash
//func (f Fragment) Hash() uint64 {
//	if f.hash == nil {
//		h := hash(f)
//		f.hash = &h
//	}
//
//	return *f.hash
//}

// Equals lets us compare fragments (all words are compared)
func (f Fragment) Equals(co Fragment) bool {
	return f.LWord == co.LWord &&
		f.CWord == co.CWord &&
		f.RWord == co.RWord
}

// GetToken takes the incoming WordType and provides the
// value of the entry associated with it.
func (f Fragment) GetToken(t WordType) TokenID {
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

// Tokens outputs the different WordType in
// an ordered slice (L1Word, L2Word, CWord, M1Word, M2Word)
func (f Fragment) Tokens() []TokenID {
	return []TokenID{
		f.LWord,
		f.CWord,
		f.RWord,
	}
}
