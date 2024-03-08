package dcli

import (
	"database/sql"
	"fmt"
	"gitlab.int.magneato.site/dungar/prototype/internal/learning"
	"io/ioutil"
	"log"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/markov"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

const m1LearnLength = 500

// M1LearnFile is our file-reading wrapper to import a text file
// into the MarkovV1 structure & tables.
func M1LearnFile(file string) {
	var (
		bytes []byte
		err   error
		str   string
		lines []string
		step  time.Time
		stop  time.Time
		max   int

		start = time.Now()
	)

	bytes, err = ioutil.ReadFile(file)

	utils.HaltingError("M1LearnFile", err)

	str = string(bytes)
	lines = learning.ReadABook(str)

	max = len(lines) - 1

	if max > m1LearnLength {
		log.Printf("Setting a hard limit to max of %d\n", m1LearnLength)
		max = m1LearnLength
	}

	DBConnect()

	if isFileAlreadyLoaded() {
		fmt.Println("File is already loaded in markov database... skipping")
		return
	}

	start = time.Now()
	stop = time.Now()
	step = time.Now()

	tx := db.BeginTransaction()

	for pos, line := range lines {
		if (pos % 50) == 0 {
			stop = time.Now()

			fmt.Printf("Learning line %d of %d lines (%v)\n", pos, max,
				stop.Sub(step).String())

			step = stop
		}

		//db.LegacyRecordRawMessage(line, source)
		markov.ActuallyLearnSentence(line)

		if pos >= max {
			break
		}
	}

	db.CommitOrRollback(tx)

	stop = time.Now()

	fmt.Printf("Learned %d of %d lines (%v, total %v)\n", max, max,
		stop.Sub(step).String(), stop.Sub(start).String())
}

func isFileAlreadyLoaded() bool {
	var (
		cnt int64
		qry string
		row *sql.Row
		err error
	)

	qry = `
		SELECT COUNT(*) AS cnt
		FROM word
	`

	row = db.ConQueryRow(qry)

	if row.Err() != nil {
		log.Printf("Failed to query row: %v\n", row.Err())
		panic(row.Err())
	}

	err = row.Scan(&cnt)

	if err != nil {
		log.Printf("Failed to scan row: %v\n", err)
		panic(err)
	}

	log.Printf("COUNT: %d\n", cnt)
	return cnt > 1
}
