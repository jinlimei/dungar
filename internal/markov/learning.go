package markov

import (
	"database/sql"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// LearnMultipleSentences learns stuff very quickly I think hopefully?
func LearnMultipleSentences(strs []string) {
	var (
		tx         = db.BeginTransaction()

		normed string
		fragmentID int64
		sentenceID int64
		wordID     int64
		words      []string
		known      map[int64]string
		unknown    []string
	)

	for _, str := range strs {
		str = strings.ToLower(str)
		sentenceID = declareSentence(tx)
		words = strings.Split(str, " ")

		known, unknown = lookupWordsByWords(words, tx)

		for _, word := range unknown {
			wordID = mergeWord(word, tx)
			known[wordID] = word

			normed = utils.Normalize(word)
			wordID = mergeWord(normed, tx)
			known[wordID] = normed
		} // end of for range unknown

		wordToID := func(word string) int64 {
			if word == "" {
				return 1
			}

			for id, knownWord := range known {
				if word == knownWord {
					return id
				}
			}

			return 1
		}

		chains := makeChains(words)

		for _, chain := range chains {
			fragmentID = makeFragment(
				wordToID(chain[1]),
				wordToID(chain[0]),
				wordToID(chain[2]),
				sentenceID,
				tx,
			)

			makeIdxFragment(
				wordToID(utils.Normalize(chain[1])),
				fragmentID,
				tx,
			)
		} // end of for range chains
	} // end of for range strs

	db.CommitOrRollback(tx)
}

// ActuallyLearnSentence is a shortcut for LearnSentence(sentence, true)
func ActuallyLearnSentence(sentence string) {
	LearnSentence(sentence, true)
}

// LearnSentence learns a sentence by recording it in the Database
func LearnSentence(sentence string, shouldCommit bool) *sql.Tx {
	sentence = strings.TrimSpace(sentence)

	if len(sentence) == 0 {
		return nil
	}

	tx := db.BeginTransaction()

	if shouldCommit {
		defer db.CommitOrRollback(tx)
	}

	sentenceID := declareSentence(tx)

	words := strings.Split(sentence, " ")

	known, unknown := lookupWordsByWords(words, tx)

	for _, word := range unknown {
		// save raw
		savedID := mergeWord(word, tx)
		known[savedID] = word

		// save normalized
		normed := utils.Normalize(word)
		savedID = mergeWord(normed, tx)
		known[savedID] = normed
	}

	wordToID := func(word string) int64 {
		// quick optimization: if it's empty we just use id=1
		if word == "" {
			return 1
		}

		for id, knownWord := range known {
			if word == knownWord {
				return id
			}
		}

		return 1
	}

	chains := makeChains(words)

	for _, chain := range chains {
		fragmentID := makeFragment(
			wordToID(chain[1]),
			wordToID(chain[0]),
			wordToID(chain[2]),
			sentenceID,
			tx,
		)

		makeIdxFragment(
			wordToID(utils.Normalize(chain[1])),
			fragmentID,
			tx,
		)
	}

	return tx
}
