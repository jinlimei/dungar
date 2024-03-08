package markov

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func haltingErr(loc string, err error) {
	utils.HaltingError("queries.go "+loc, err)
}

func getFragmentFromWordCount(word string) int64 {
	query := `
		SELECT COUNT(*) AS cnt
		FROM idx_word_fragment iwf
		JOIN fragment f ON f.id = iwf.fragment_id
		JOIN word     w ON w.id = iwf.word_id
		WHERE w.word = $1
	`

	row := db.ConQueryRow(query, word)

	var cnt int64

	err := row.Scan(&cnt)
	haltingErr("getFragmentFromWordCount", err)

	return cnt
}

func getFragmentFromWord(word string) *fragment {
	count := getFragmentFromWordCount(word)

	log.Println("getFragmentFromWordCount", word, count)

	if count == 0 {
		if utils.IsEmoticon(word) {
			word = strings.Trim(word, ":")
			count = getFragmentFromWordCount(word)
			log.Println("getFragmentFromWord (treating as emoji, cutting)", word, count)
		}

		// Try the word lower-case first
		wordLower := strings.ToLower(word)
		count = getFragmentFromWordCount(wordLower)

		if count > 0 {
			log.Printf("could not find upper-case of word '%s' so ended up on lower-case '%s'\n", word, wordLower)
			word = wordLower
		}

		for count == 0 {
			log.Printf("encountered unknown word '%s' (and lower '%s'), picking another one at random", word, wordLower)
			word = PickWord()

			if utils.IsURL(word) {
				log.Println("skipping url for word generation")
				continue
			}

			if word == "" {
				log.Println("Retrieved empty word, cannot generate markov")
				return nil
			}

			count = getFragmentFromWordCount(word)
		}

		log.Printf("ended up on '%s' which i guess is fine (count: %d)\n", word, count)
	}

	pos := int64(0)

	if count > 1 {
		pos = random.Int64(count-1) + 1
	}

	log.Println("getFragmentFromWord", word, count)

	query := `
		SELECT f.id, f.r_word_id, f.l_word_id, f.sentence_id, f.word_id
		FROM idx_word_fragment iwf
		JOIN fragment f ON f.id = iwf.fragment_id
		JOIN word     w ON w.id = iwf.word_id
		WHERE w.word = $1
		LIMIT 1 OFFSET %d
	`

	row := db.ConQueryRow(fmt.Sprintf(query, pos), word)

	var id, rWordID, lWordID, sentID, wordID int64
	err := row.Scan(&id, &rWordID, &lWordID, &sentID, &wordID)

	haltingErr("getFragmentFromWord", err)

	return &fragment{
		ID:         id,
		RWordID:    rWordID,
		LWordID:    lWordID,
		SentenceID: sentID,
		WordID:     wordID,
		Word:       word,
	}
}

func lookupStatesByIDs(ids []int64, con db.QueryAble) string {
	if len(ids) == 0 {
		return ""
	}

	query := `
		SELECT w.id, w.word
		FROM word w
		WHERE w.id IN (%s)
	`

	stringIds := make([]string, len(ids))

	for pos, id := range ids {
		stringIds[pos] = strconv.FormatInt(id, 10)
	}

	query = fmt.Sprintf(query, strings.Join(stringIds, ","))

	res := db.MustQuery(con, query)

	words := make(map[int64]string, 0)

	defer res.Close()
	for res.Next() {
		var id int64
		var word string

		if err := res.Scan(&id, &word); err != nil {
			haltingErr("lookupStatesByIDs", err)
		}

		words[id] = word
	}

	output := ""

	for _, id := range ids {
		output += words[id] + " "
	}

	return strings.TrimSpace(output)
}

func queryChainCount(frag *fragment, source, target string, con db.QueryAble) int64 {
	query := `
		SELECT COUNT(*) cnt
		FROM fragment f
		WHERE f.word_id = $1
		AND f.%s = $2
	`

	row := db.QueryRow(
		con,
		fmt.Sprintf(
			query,
			frag.column(source),
		),
		// getattr(fragment, t_word_id)
		frag.intAttr(target),
		// getattr(fragment, s_word_id)
		frag.WordID,
	)

	var cnt int64
	err := row.Scan(&cnt)
	haltingErr("queryChainCount", err)

	return cnt
}

func queryChain(frag *fragment, source, target string, con db.QueryAble) *fragment {
	count := queryChainCount(frag, source, target, con)

	if count <= 0 {
		return nil
	}

	var pos int64

	if count > 1 {
		pos = 1 + random.Int64(count-1)
	}

	query := `
		SELECT *
		FROM fragment f
		WHERE f.word_id = $1
		AND f.%s = $2
		LIMIT 1 OFFSET %d
	`

	row := db.QueryRow(
		con,
		fmt.Sprintf(
			query,
			// getattr(fragment, s_word_id)
			frag.column(source),
			pos,
		),
		// getattr(fragment, t_word_id)
		frag.intAttr(target),
		// fragment.word_id
		frag.WordID,
	)

	var id, rWordID, lWordID, sentID, wordID int64

	if err := row.Scan(&id, &rWordID, &lWordID, &sentID, &wordID); err != nil {
		haltingErr("queryChain", err)
	}

	return &fragment{
		ID:         id,
		RWordID:    rWordID,
		LWordID:    lWordID,
		SentenceID: sentID,
		WordID:     wordID,
	}
}

func getGlobalWordEmoteLists(con db.QueryAble) ([]string, []string) {
	query := "SELECT DISTINCT word FROM word WHERE word != '' "
	res := db.MustQuery(con, query)

	emoticons := make(map[string]int, 0)
	normalized := make(map[string]int, 0)

	defer res.Close()
	for res.Next() {
		var word string
		haltingErr("getGlobalWordEmoteLists", res.Scan(&word))

		if utils.IsEmoticon(word) {
			word = strings.ToLower(word)

			_, emoteOk := emoticons[word]
			if !emoteOk {
				emoticons[word] = 1
			}
		} else {
			word = utils.Normalize(word)

			_, wordOk := normalized[word]
			if !wordOk {
				normalized[word] = 1
			}
		}
	}

	normOut := make([]string, 0)
	emoteOut := make([]string, 0)

	for norm := range normalized {
		normOut = append(normOut, norm)
	}

	for emote := range emoticons {
		emoteOut = append(emoteOut, emote)
	}

	return normOut, emoteOut
}

func lookupWordID(word string, con db.QueryAble) int64 {
	query := `
		SELECT w.id
		FROM word w
		WHERE w.word = $1
	`

	row := db.QueryRow(con, query, word)

	var id int64

	err := row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		}

		log.Printf("Error in lookupWordID('%s', con)=> %v,, %v\n", word, id, err)
	}

	return id
}

func lookupWordsByWords(words []string, con db.QueryAble) (map[int64]string, []string) {
	output := make(map[int64]string, 0)
	unknown := make([]string, 0)

	for _, word := range words {
		res := lookupWordID(word, con)

		if res <= 0 {
			unknown = append(unknown, word)
		} else {
			output[res] = word
		}
	}

	return output, unknown
}

func mergeWord(word string, con db.QueryAble) int64 {
	lookup := lookupWordID(word, con)
	if lookup > 0 {
		return lookup
	}

	query := `
		INSERT INTO word (word)
		VALUES ($1)
		RETURNING id
	`

	res := db.QueryRow(con, query, word)

	var id int64
	err := res.Scan(&id)
	haltingErr("mergeWord ID", err)

	return id
}

func declareSentence(con db.QueryAble) int64 {
	query := `
		INSERT INTO sentence ( "timestamp" )
		VALUES ( extract(epoch from now() at time zone 'utc') )
		RETURNING id
	`

	res := db.QueryRow(con, query)

	var id int64
	err := res.Scan(&id)
	haltingErr("declareSentence ID", err)

	return id
}

func makeFragment(wordID, leftWord, rightWord, sentenceID int64, con db.QueryAble) int64 {
	query := `
		INSERT INTO fragment (r_word_id, l_word_id, sentence_id, word_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	res := db.QueryRow(con, query, rightWord, leftWord, sentenceID, wordID)

	var id int64
	err := res.Scan(&id)
	haltingErr("makeFragment ID", err)

	return id
}

func makeIdxFragment(wordID, fragmentID int64, con db.QueryAble) int64 {
	query := `
		INSERT INTO idx_word_fragment (word_id, fragment_id)
		VALUES ($1, $2)
		RETURNING id
	`

	res := db.QueryRow(con, query, wordID, fragmentID)

	var id int64
	err := res.Scan(&id)
	haltingErr("makeIdxFragment ID", err)

	return id
}
