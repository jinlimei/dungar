package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestLotsOfSpaceHandler(t *testing.T) {
	svc := initMockServices()

	rsp := lotsOfSpaceHandler(svc, core2.MakeSingleRsp("  hello         world "))
	assert.Equal(t, "hello world", rsp[0].Contents)
}

func TestRemoveGarbagePrefixedHandler(t *testing.T) {
	svc := initMockServices()

	rsp := removeGarbagePrefixedHandler(svc, core2.MakeSingleRsp("<Some Name> thats not true"))
	assert.Equal(t, "thats not true", rsp[0].Contents)
}
