package triggers

import (
	"testing"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestPreBadWordsHandler(t *testing.T) {
	random.UseTestSeed()
	//db.TestDatabaseConnect()

	badWords = []string{"foo", "butt", "bu"}
	badWordCheck = time.Now().Unix()

	msg := makeMessage("hey butt!", "", "")
	rsp := preBadWordsHandler(msg)

	assert.True(t, rsp[0].IsEmpty())
	assert.True(t, rsp[0].IsCancelled())
}

func TestPostBadWordsHandler(t *testing.T) {
	svc := initMockServices()
	random.UseTestSeed()
	//db.TestDatabaseConnect()

	badWords = []string{"foo", "butt", "bu"}
	badWordCheck = time.Now().Unix()

	rsp := core2.MakeSingleRsp("foo buz baz bar")
	nex := postBadWordsHandler(svc, rsp)

	assert.False(t, nex[0].IsEmpty())
	assert.True(t, nex[0].IsCancelled())

	rsp = core2.MakeSingleRsp("FoO buz baz bar")
	nex = postBadWordsHandler(svc, rsp)

	assert.False(t, nex[0].IsEmpty())
	assert.True(t, nex[0].IsCancelled())

	rsp = core2.MakeSingleRsp("buz baz bar")
	nex = postBadWordsHandler(svc, rsp)

	assert.False(t, nex[0].IsEmpty())
	assert.Equal(t, rsp, nex)
}

func TestShouldTossMessage(t *testing.T) {
	random.UseTestSeed()

	badWordCheck = time.Now().Unix()
	badWords = []string{"foo", "bar"}

	assert.True(t, shouldTossMessage("hello foo"))
	assert.True(t, shouldTossMessage("foo"))
	assert.True(t, shouldTossMessage("But what about FoO"))
	assert.True(t, shouldTossMessage("Bar Foo"))
	assert.True(t, shouldTossMessage("Bar, Foo"))
	assert.True(t, shouldTossMessage("Baz, Foo!"))
	assert.False(t, shouldTossMessage("Hello, World!"))
}

func TestBadWordListCheck(t *testing.T) {
	random.UseTestSeed()
	db.TestDatabaseConnect()

	badWordCheck = 0
	badWords = []string{"b"}

	badWordListCheck()

	assert.NotEqual(t, []string{"b"}, badWords)
	assert.NotEqual(t, 0, badWordCheck)
}
