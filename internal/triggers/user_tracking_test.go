package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestLastSeenHandler(t *testing.T) {
	svc := initMockServices()

	random.UseTestSeed()
	db.TestDatabaseConnect()

	msg := makeMessage("!seen foo", "foo", "bar")
	rsp := userTrackingHandler(svc, msg)
	assert.Truef(t, strings.Contains(rsp[0].Contents, "could not find user"),
		"Failed to get valid results for '!seen foo', got '%+v' instead", rsp)

	msg.Contents = "!seen v"

	rsp = userTrackingHandler(svc, msg)
	assert.Truef(t, strings.Contains(rsp[0].Contents, "was last seen"),
		"Failed to get valid results for '!seen v', got '%+v' instead", rsp)

	msg = makeMessage("!top", "foo", "bar")
	rsp = userTrackingHandler(svc, msg)
	//spew.Dump(rsp[0].Contents)
	assert.Truef(t, strings.Contains(rsp[0].Contents, "01: foo with 2 posts"),
		"Failed to get valid results for '!top', got '%+v' instead", rsp)
}
