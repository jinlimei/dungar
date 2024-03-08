package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestNameCheckHandler(t *testing.T) {
	random.UseTestSeed()

	svc := initMockServices()
	msg := makeMessage("@Dungar do you like butts?", "", "")
	pass := false
	i := 0
	for ; i < 500; i++ {
		rsp := nameCheckHandler(svc, msg)

		if strings.Contains(rsp[0].Contents, "name") {
			pass = true
			break
		}
	}

	//if pass {
	//	log.Printf("TestNameCheckHandler passed after %d iterations\n", i)
	//}

	assert.True(t, pass)
}
