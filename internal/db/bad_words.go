package db

import (
	"database/sql"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// GetBadWords retrieves a list of bad words (or empty if there's no results)
func GetBadWords() []string {
	if utils.InTestEnv() {
		return []string{
			"fizz",
			"buzz",
		}
	}

	query := `
		SELECT word
		FROM bad_words
		WHERE active = 1
	`

	rows := ConMustQuery(query)

	var (
		tmp string
		out = make([]string, 0)
	)

	for rows.Next() {
		if err := rows.Scan(&tmp); err != nil {

			if sql.ErrNoRows == err {
				break
			}

			utils.HaltingError("GetBadWords", err)
		}

		out = append(out, tmp)
	}

	return out
}
