package dcli

import (
	"database/sql"
	"encoding/json"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type config struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Data string `json:"data"`
}

type fragment struct {
	ID         uint64
	LWordID    uint64
	CWordID    uint64
	RWordID    uint64
	SentenceID uint64
}

// ReconstructFragments will take in the markov tables and reconstruct
// them into a raw_messages_m1 table
func ReconstructFragments(file string) {
	var cfg *config

	rawc, err := ioutil.ReadFile(file)
	utils.HaltingError("ReconstructFragments,readfile", err)

	err = json.Unmarshal(rawc, &cfg)
	utils.HaltingError("ReconstructFragments,unmarshal", err)

	db.InlineConnectToDatabase(
		cfg.User,
		cfg.Pass,
		cfg.Host,
		cfg.Data,
	)

	var (
		words     = retrieveAllWords()
		fragments = retrieveAllFragments()

		tx = db.BeginTransaction()

		start = time.Now()
		step  = time.Now()
		now   = time.Now()

		sentencesWritten = 0
		fragmentsWritten = 0
	)

	for sID, frags := range fragments {
		tmp := make([]string, 0, len(frags))

		for _, frag := range frags {
			tmp = append(tmp, words[frag.CWordID])
		}


		//fmt.Printf("Words: %d\n", len(tmp))
		//fmt.Println(strings.Join(tmp, " "))
		writeRawMessage(tx, sID, strings.Join(tmp, " "))

		fragmentsWritten += len(tmp)
		sentencesWritten++

		now = time.Now()
		if now.Sub(step).Seconds() > 5 {
			step = now
			log.Printf("Sentences Written: %d, Fragments Used: %d (@%s)\n",
				sentencesWritten, fragmentsWritten, now.Sub(start).String())
		}
	}

	db.CommitOrRollback(tx)
	end := time.Now()
	log.Printf("Took %s to write %d sentences with %d fragments\n",
		end.Sub(start).String(), sentencesWritten, fragmentsWritten)
}

func writeRawMessage(tx *sql.Tx, sID uint64, msg string) {
	qry := `
		INSERT INTO raw_messages_m1 ("sentence_id", "message", "created_at")
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT ("sentence_id") DO 
		UPDATE SET 
			"message" = $2,
			"created_at" = CURRENT_TIMESTAMP
	`

	//log.Printf("I suppose we'd qry here:\n%s\n", qry)
	db.MustExec(tx, qry, sID, msg)
}

func retrieveAllFragments() map[uint64][]fragment {
	var (
		// We already have an idea of how many fragments there are
		fragments = make(map[uint64][]fragment, 10_000)

		rows *sql.Rows
		err  error

		fID  uint64
		lwID uint64
		cwID uint64
		rwID uint64
		sID  uint64

		ok bool
	)

	qry := `
		SELECT id, l_word_id, word_id, r_word_id, sentence_id
		FROM fragment 
		ORDER BY sentence_id, id
	`

	rows = db.ConMustQuery(qry)
	utils.HaltingError("retrieveAllWords,rows.Err()", rows.Err())

	retrieved := 0
	for rows.Next() {
		err = rows.Scan(&fID, &lwID, &cwID, &rwID, &sID)
		utils.HaltingError("retrieveAllWords,rows.Scan()", err)

		_, ok = fragments[sID]
		if !ok {
			fragments[sID] = make([]fragment, 0, 8)
		}

		fragments[sID] = append(fragments[sID], fragment{
			ID:         fID,
			LWordID:    lwID,
			CWordID:    cwID,
			RWordID:    rwID,
			SentenceID: sID,
		})

		retrieved++
	}

	log.Printf("Retrieved %d fragments\n", retrieved)
	return fragments
}

func retrieveAllWords() map[uint64]string {
	var (
		// We can cheat here and pre-populate the # of words in the database
		// since it's mostly read-only for this particular feature
		words = make(map[uint64]string, 458_000)

		id uint64
		w  string

		rows *sql.Rows
		err  error
	)

	qry := `
		SELECT id, word 
		FROM word 
		ORDER BY id
	`

	rows = db.ConMustQuery(qry)
	utils.HaltingError("retrieveAllWords,rows.Err()", rows.Err())

	retrieved := 0
	for rows.Next() {
		err = rows.Scan(&id, &w)
		utils.HaltingError("retrieveAllWords,rows.Scan()", err)

		if id > 0 {
			retrieved++
			words[id] = w
		}
	}

	log.Printf("Retrieved '%d' words\n", retrieved)
	return words
}
