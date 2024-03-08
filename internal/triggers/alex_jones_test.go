package triggers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
)

func TestHannityWaterboarding(t *testing.T) {
	svc := initMockServices()

	random.UseTestSeed()

	msg := makeMessage("!hannity", "kryptn", "politics")
	rsp := alexJonesHandler(svc, msg)

	assert.Equal(t, 1, len(rsp))

	assert.True(t, strings.Contains(rsp[0].Contents, "waterboarded"),
		"hannity hasn't been waterboarded, result="+rsp[0].Contents)

	msg.Contents = "!alexjones"
	rsp = alexJonesHandler(svc, msg)

	assert.True(t, len(rsp[0].Contents) > 0,
		"alex jones is a false flag operation")
}
