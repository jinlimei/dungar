package core2

import "sync"

// OutgoingResponses is a (supposedly) goroutine-safe method for working with a set of
// ResponseEnvelope's in a slice. This struct is basically just a sync.RWMutex around the
// various operations necessary for a slice to work with multiple goroutines
type OutgoingResponses struct {
	envelopes []*ResponseEnvelope
	mutex     sync.RWMutex
}

// NewOutgoingResponses creates a new OutgoingResponses struct with the envelopes slice
// already initialized correctly
func NewOutgoingResponses() *OutgoingResponses {
	return &OutgoingResponses{
		envelopes: make([]*ResponseEnvelope, 0),
		mutex:     sync.RWMutex{},
	}
}

// AddEnvelope tries to safely add an envelope to the set of envelopes present.
// This will not interact with the mutex or envelopes slice if env is nil
func (ow *OutgoingResponses) AddEnvelope(env *ResponseEnvelope) {
	if env == nil {
		return
	}

	ow.mutex.Lock()
	defer ow.mutex.Unlock()

	ow.envelopes = append(ow.envelopes, env)
}

// Reset will attempt to safely clear out the envelopes slice
func (ow *OutgoingResponses) Reset() {
	ow.mutex.Lock()
	defer ow.mutex.Unlock()

	ow.envelopes = make([]*ResponseEnvelope, 0)
}

// Length tries to safely return the length of available envelopes
func (ow *OutgoingResponses) Length() int {
	ow.mutex.RLock()
	defer ow.mutex.RUnlock()

	return len(ow.envelopes)
}

// HasEnvelopes tries to safely check if there are any envelopes present
func (ow *OutgoingResponses) HasEnvelopes() bool {
	ow.mutex.RLock()
	defer ow.mutex.RUnlock()

	return len(ow.envelopes) > 0
}

// RetrieveEnvelopes tries to safely return a slice of ResponseEnvelope's and also
// reset envelopes to an empty value
func (ow *OutgoingResponses) RetrieveEnvelopes() []*ResponseEnvelope {
	ow.mutex.Lock()
	defer ow.mutex.Unlock()

	var (
		upd = make([]*ResponseEnvelope, 0, len(ow.envelopes))
		out = ow.envelopes
	)

	ow.envelopes = upd

	return out
}
