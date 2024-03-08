package markov

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

func connect() {
	//if db.GetDatabase() != nil {
	//	return
	//}
	//
	//utils.LoadSettingsAndSecrets()
	//db.ConnectToDatabase()
}

func TestGetFragmentFromWordCount(t *testing.T) {
	t.SkipNow()

	random.UseTestSeed()
	connect()

	testWord := func (word, condName string, cond func(i int64) bool) {
		c := getFragmentFromWordCount(word)
		assert.True(t, cond(c), "Word '" + word + "' should pass condition " + condName)
	}

	shouldBeZero := func (c int64) bool { return c == 0 }
	shouldBeMany := func (c int64) bool { return c > 0 }

	testWord("butt", "shouldBeZero", shouldBeZero)
	testWord("https://gist.github.com/metatoaster/e20394f6377f5aa05272f21f259fc5a8", "shouldBeZero", shouldBeZero)
	testWord("https://made.up.url", "shouldBeZero", shouldBeZero)
	testWord(":condi:", "shouldBeZero", shouldBeZero)
	testWord(":thinking:", "shouldBeZero", shouldBeZero)
	testWord("english", "shouldBeMany", shouldBeMany)
	testWord("complained", "shouldBeMany", shouldBeMany)
}

func TestGetFragmentFromWord(t *testing.T) {
	t.SkipNow()

	random.UseTestSeed()
	connect()

	assertValidFragment := func (frag *fragment) {
		assert.NotNil(t, frag)

		if frag != nil {
			assert.True(t, frag.ID > 0)
			assert.True(t, frag.RWordID > 0)
			assert.True(t, frag.LWordID > 0)
			assert.True(t, frag.SentenceID > 0)
		}
	}

	frag := getFragmentFromWord("english")

	assertValidFragment(frag)

	frag = getFragmentFromWord("english")

	assertValidFragment(frag)

	word := "http://www.gutenberg.org"
	frag = getFragmentFromWord(word)

	assert.True(t, utils.IsURL(word))
	assertValidFragment(frag)

	word = "http://pglaf.org/donate"
	frag = getFragmentFromWord(word)

	assert.True(t, utils.IsURL(word))
	assertValidFragment(frag)
}

func TestQueryChain(t *testing.T) {
	t.SkipNow()

	random.UseTestSeed()
	connect()

	var count int64
	var frag *fragment

	for count == 0 {
		frag  = getFragmentFromWord("english")

		if frag == nil {
			break
		}

		count = queryChainCount(frag, "RWordID", "LWordID", nil)
	}

	assert.NotNil(t, frag)
	assert.True(t, count >= 1)

	if frag == nil {
		assert.Fail(t, "Fragment is nil so we cannot proceed further")
		return
	}

	frag2 := queryChain(frag, "RWordID", "LWordID", nil)

	assert.NotNil(t, frag2)
	assert.True(t, frag2.ID > 0)
}

func TestQueryChain2(t *testing.T) {
	t.SkipNow()

	random.UseTestSeed()
	connect()
	frag := &fragment{
		1,
		2,
		3,
		4,
		5,
		"",
	}

	queryBuilder := func(fragment *fragment, source, target string) (string, int64, int64) {
		query := `
			SELECT *
			FROM fragment f
			WHERE f.word_id = $1
			AND f.%s = $2
			LIMIT 1 OFFSET %d
		`

		out := fmt.Sprintf(
			query,
			// getattr(fragment, s_word_id)
			frag.column(source), 0,
		)

		return out, frag.intAttr(target), frag.WordID
	}

	lrSQL, lrID, lrWordID := queryBuilder(frag, "LWordID", "RWordID")

	fmt.Println(lrSQL)
	fmt.Println(lrID)
	fmt.Println(lrWordID)

	rlSQL, rlID, rlWordID := queryBuilder(frag, "RWordID", "LWordID")

	fmt.Println(rlSQL)
	fmt.Println(rlID)
	fmt.Println(rlWordID)
}

func TestGenerate(t *testing.T) {
	t.SkipNow()

	random.UseTestSeed()
	connect()

	assert.True(t, "" != Generate("dovi"))
	assert.True(t, "" != Generate("rand.Int64n(-1)"))
}
