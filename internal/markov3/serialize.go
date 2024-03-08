package markov3

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Serialize takes in the existing struct data and pushes it into a
// gob encoded byte slice
func (m *Markov) Serialize() []byte {
	var (
		buf = &bytes.Buffer{}
		enc = gob.NewEncoder(buf)
		err = enc.Encode(m)
	)

	if err != nil {
		log.Fatalf("Could not serialize Markov: %v\n", err)
	}

	return buf.Bytes()
}
