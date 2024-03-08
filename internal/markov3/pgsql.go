package markov3

import (
	"database/sql"
	"log"

	"gitlab.int.magneato.site/dungar/prototype/internal/cleaner"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
)

// LearnFromRawMessagesM1 is a method for learning from the legacy raw messages
// table built from the markov database of it
func (m *Markov) LearnFromRawMessagesM1() int {
	if utils.InTestEnv() {
		return -1
	}

	if m.learnedFromLegacy {
		return -1
	}

	db.EnsureDatabaseConnection()

	qry := `
		SELECT sentence_id, message
		FROM raw_messages_m1
	`

	var (
		res     *sql.Rows
		id      uint64
		msg     string
		learned = 0
	)

	res = db.ConMustQuery(qry)

	for res.Next() {
		if err := res.Scan(&id, &msg); err != nil {
			log.Printf("Failed to scan: %v\n", err)
			break
		}

		spc := spaceCount(msg)

		if spc > 0 {
			m.LearnString(msg, cleaner.VariantXMPP)
			learned++
		}
	}

	m.learnedFromLegacy = true
	return learned
}

func spaceCount(s string) int {
	var (
		rs = []rune(s)
		cn = 0
	)

	for k := 0; k < len(rs); k++ {
		if rs[k] == ' ' {
			cn++
		}
	}

	return cn
}

// LearnFromRawMessages will pull in the `raw_messages`
// table from PGSQL and push them into the internal
// markov state
func (m *Markov) LearnFromRawMessages() int {
	if utils.InTestEnv() {
		return -1
	}

	var res *sql.Rows

	db.EnsureDatabaseConnection()

	if m.lastRawMessageID == 0 {
		// super simple: just learn everything nbd
		qry := `
			SELECT id, message 
			FROM raw_messages
		`
		log.Printf("[LearnFromRawMessages] lastRawMessageID=0, starting fresh")
		res = db.ConMustQuery(qry)
	} else {
		qry := `
			SELECT id, message 
			FROM raw_messages
			WHERE id > $1 
		`

		log.Printf("[LearnFromRawMessages] lastRawMessageID=%d, starting from there\n",
			m.lastRawMessageID)
		res = db.ConMustQuery(qry, m.lastRawMessageID)
	}

	var (
		id      uint64
		msg     string
		learned = 0
	)

	for res.Next() {
		if err := res.Scan(&id, &msg); err != nil {
			log.Printf("Failed to scan: %v\n", err)
			break
		}

		if len(msg) > 1 {
			m.LearnString(msg, cleaner.VariantSlack)
			m.lastRawMessageID = id
			learned++
		}
	}

	return learned
}

func keyToValues(m map[string]int) []string {
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
