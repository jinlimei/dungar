package markov3

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/dungar/internal/db"
	"gitlab.int.magneato.site/dungar/dungar/internal/random"
)

func zTestMarkovLearnFromRawMessages(t *testing.T) {
	assert.Equal(t, 1, 1)

	db.TestDatabaseConnect()
	random.UseTestSeed()

	n := time.Now()
	m := MakeMarkov("")
	m.LearnFromRawMessagesM1()

	for k := range 50 {
		log.Printf("%02d: %s\n", k, m.Generate("genius"))
	}

	log.Printf("Finished learning, having %d words and %d fragments. Took %s\n",
		len(m.RevWords), len(m.Fragments), time.Since(n).String())
}
