package markov3

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func zTestMarkovLearnFromRawMessages(t *testing.T) {
	assert.Equal(t, 1, 1)

	db.TestDatabaseConnect()
	random.UseTestSeed()

	n := time.Now()
	m := MakeMarkov("")
	m.LearnFromRawMessagesM1()

	for k := 0; k < 50; k++ {
		log.Printf("%02d: %s\n", k, m.Generate("genius"))
	}

	log.Printf("Finished learning, having %d words and %d fragments. Took %s\n",
		len(m.RevWords), len(m.Fragments), time.Now().Sub(n).String())
}
