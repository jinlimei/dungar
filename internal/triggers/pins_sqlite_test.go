package triggers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"strings"
	"testing"
)

func TestPinDBService(t *testing.T) {
	random.UseTestSeed()

	svc, err := startPinService()
	assert.Nilf(t, err, "err was not nil: %v", err)
	assert.NotNil(t, svc, "svc was nil")
	if err != nil || svc == nil {
		return
	}

	// Requesting All Pins
	pins, err := svc.requestAllPins()
	assert.Nilf(t, err, "request all pins er was not nil: %v", err)
	assert.Truef(t, len(pins) > 0, "pins returned had len '%d' (nil? %v)", len(pins), pins)
	fmt.Println(len(pins))
}

func TestPinDBHandler(t *testing.T) {
	random.UseTestSeed()

	svc, err := startPinService()
	assert.Nilf(t, err, "err was not nil: %v", err)
	assert.NotNil(t, svc, "svc was nil")
	if err != nil || svc == nil {
		return
	}

	var rsp []*core2.Response

	rsp = pinDBHandler(makeMessage("hey how are you", "", ""))
	assert.Truef(t, isEmptyRsp(rsp), "rsp is not empty: %+v", rsp)

	rsp = pinDBHandler(makeMessage("!pins", "", ""))
	assert.Truef(t, !isEmptyRsp(rsp), "rsp is empty: %+v", rsp)
	assert.Truef(t, strings.Contains(rsp[0].Contents, "tastosis"), "rsp does not contain tastosis: %+v", rsp)

	rsp = pinDBHandler(makeMessage("!pins tastosis", "", ""))
	assert.Truef(t, !isEmptyRsp(rsp), "rsp is empty: %+v", rsp)
	assert.Truef(t, strings.Contains(rsp[0].Contents, "tastosis"), "rsp does not contain tastosis: %+v", rsp)

	rsp = pinDBHandler(makeMessage("!pins @dungar", "", ""))
	assert.Truef(t, !isEmptyRsp(rsp), "rsp is empty: %+v", rsp)
	assert.Truef(t, strings.Contains(rsp[0].Contents, "bruh"), "rsp does not contain tastosis: %+v", rsp)
}
