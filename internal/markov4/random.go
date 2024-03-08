package markov4

import "gitlab.int.magneato.site/dungar/prototype/internal/random"

// randomMarkovID lets us box in the usage of the underlying
// data type of MarkovID to a single function so we can then
// increase it from current (uint32) to maybe something more
// in the future (int64? uint64?) if necessary. Or lower if
// desired as well (welcome back int32)
func randomMarkovID(n MarkovID) MarkovID {
	return MarkovID(random.Int64(int64(n)))
}
