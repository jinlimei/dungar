package markov

import (
	"log"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

var normalizedWords []string
var globalFetchTime = time.Now()
var wordListMaxAge = int64(600)
var maxChainDistance = 50
var emoticons []string
var retrievedThings = -1

// This is based off of the work done by the venerable meta toaster
// https://github.com/metatoaster/mtj.markov.git

func splitDirection(dir string) (string, string) {
	// LR -> LWordID, RWordID
	// RL -> RWordID, LWordID
	return string(dir[0]) + "WordID", string(dir[1]) + "WordID"
}

func buildGlobalWordList() {
	regenerate := false


	if retrievedThings < 0 && (normalizedWords == nil || len(normalizedWords) == 0) {
		log.Println("normalizedWords is nil or empty")
		regenerate = true
	}

	if time.Now().Unix()-globalFetchTime.Unix() > wordListMaxAge {
		log.Println("normalizedWords age is too old")
		regenerate = true
	}

	log.Printf("should regenerate? %v\n", regenerate)

	if regenerate {
		normalizedWords, emoticons = getGlobalWordEmoteLists(nil)
		globalFetchTime = time.Now()

		retrievedThings = len(normalizedWords) + len(emoticons)

		log.Printf("acquired %d normalizedWords\n", len(normalizedWords))
		log.Printf("acquired %d emoticons\n", len(emoticons))
	}
}

// PickWord takes a word from the global word list and returns it
func PickWord() string {
	buildGlobalWordList()
	return random.PickString(normalizedWords)
}

// GetRandomEmoticon takes an emoticon from the global emoticon list
func GetRandomEmoticon() string {
	buildGlobalWordList()
	return random.PickString(emoticons)
}

func pickEntryPoint(frag *fragment) *fragment {
	word := frag.Word

	if word == "" {
		word = PickWord()
	}

	return getFragmentFromWord(word)
}

// Generate builds a markov chain from a word
func Generate(word string) string {
	return utils.TrimPunctuation(GenerateFromCon(word, nil))
}

// GenerateFromCon generates a word from a given QueryAble connection
func GenerateFromCon(word string, con db.QueryAble) string {
	frag := &fragment{
		Word: word,
	}

	frag = pickEntryPoint(frag)

	if frag == nil {
		// something went horribly wrong
		log.Println("Fragment was <nil> when it shouldn't be.")
		//utils.SentryMessage(
		//	"fragment <nil> when shouldn't be",
		//	map[string]string{"loc": "GenerateFromCon"},
		//)

		return ""
	}

	lhs := followFragmentChain(frag, "RL")
	c := frag.getWordIds()
	rhs := followFragmentChain(frag, "LR")

	return lookupStatesByIDs(utils.JoinInt64s(lhs, c, rhs), con)
}

func followFragmentChain(fragment *fragment, destination string) []int64 {
	results := make([]int64, 0)

	frag := fragment

	source, target := splitDirection(destination)

	for i := 0; i < maxChainDistance; i++ {
		frag = queryChain(frag, source, target, nil)

		if frag == nil {
			break
		}

		// append(getattr(fragment, t_word_id))
		results = append(results, frag.intAttr(target))
	}

	if target == "LWordID" {
		return utils.ReverseInt64s(results)
	}

	return results
}

func messageToWords(message string) []string {
	words := strings.Split(message, " ")

	if len(words) == 0 || words[0][0] == '!' {
		return nil
	}

	return words
}

func joinerWords() []string {
	return []string{
		"a",
		"if",
		"its",
		"it's",
		"and",
		"or",
		"because",
		"with",
		"when",
		"like",
		"then",
		"than",
		"after",
		"also",
		"before",
	}
}

func isJoinerWord(word string) bool {
	if len(word) > 1 && utils.EndsWithPunctuation(word) {
		return true
	}

	word = strings.ToLower(word)

	for _, val := range joinerWords() {
		if val == word {
			return true
		}
	}

	return false
}
