package dcli

import (
	"io/ioutil"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// ImportFortunes takes the incoming file and builds
// the fortunes table with the contents of that file.
func ImportFortunes(file string) {
	DBConnect()

	data, err := ioutil.ReadFile(file)
	utils.HaltingError("ImportFortunes()", err)

	fortunes := strings.Split(string(data), "\n%\n")
	imports := 0
	empty := 0

	for _, fortune := range fortunes {
		fortune = strings.TrimSpace(fortune)

		if fortune == "" {
			log.Printf("Skipping empty line\n")
			empty++
			continue
		}

		if insertFortune(fortune) {
			imports++
		}
	}

	log.Printf("Imported %d out of %d\n", imports, len(fortunes))
	log.Printf("Skipped %d out of %d because of empty line\n", empty, len(fortunes))
}

func hasFortune(fortune string) bool {
	query := `
		SELECT COUNT(*) AS cnt
		FROM fortunes
		WHERE fortune = $1
	`

	row := db.ConQueryRow(query, fortune)
	var cnt int

	if err := row.Scan(&cnt); err != nil {
		log.Fatal(err)
	}

	return cnt > 0
}

func insertFortune(fortune string) bool {
	if hasFortune(fortune) {
		return false
	}

	query := `
		INSERT INTO fortunes (fortune, added)
		VALUES($1, CURRENT_TIMESTAMP)
	`

	db.ConMustExec(query, fortune)
	return true
}

