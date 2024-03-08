package triggers

import (
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestNoYouRegexp(t *testing.T) {
	assert.True(t, noYouFuckYou.MatchString("fuck you"))
}

func TestNoYouHandler(t *testing.T) {
	random.UseTestSeed()

	masterChanceList["noYouHandler--regexp"] = 0.5
	masterChanceList["noYouHandler--normal"] = 0.5

	testRoutine(t, "fuck you")
	testRoutine(t, "@bob fuck you")
	testRoutine(t, "fuck u")
	testRoutine(t, "@bob fuck u")
	testRoutine(t, "fuck ur house")
	testRoutine(t, "@bob fuck ur house")
	testRoutine(t, "ur mom is a phone")
	testRoutine(t, "@bob ur mom is a phone")
	testRoutine(t, "ur a phone")
	testRoutine(t, "@bob ur a phone")
}

func testRoutine(t *testing.T, str string) {
	msg := makeMessage(str, "", "")
	pass := false

	i := 0
	for ; i < 100; i++ {
		rsp := noYouHandler(msg)

		if utils.StringInSlice(rsp[0].Contents, noYouResponses) {
			pass = true
			break
		}
	}

	//if pass {
	//	log.Printf("noYouHandler passed against '%s' message after %d iterations\n", str, i)
	//}

	assert.Truef(t, pass,
		"noYouHandler failed against '%s' after %d iterations\n", str, i)
}
