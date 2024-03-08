package triggers

import (
	"fmt"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func postWeirdOutputFilterHandler(svc *core2.Service, rsps []*core2.Response) []*core2.Response {
	for _, rsp := range rsps {
		if isJulietNumberSequence(rsp.Contents) {
			rsp.Cancel()
		}
	}

	return rsps
}

func isJulietNumberSequence(str string) bool {
	digits := 0
	nonDigits := 0

	for _, chr := range str {
		if strings.IndexAny("0123456789,.", string(chr)) >= 0 {
			digits++
		} else {
			nonDigits++
		}
	}

	if nonDigits == 0 || digits == 0 {
		return false
	}

	ratio := float64(digits / nonDigits)

	if ratio > 2.0 {
		log.Printf("Ratio is >0: %v\n", ratio)

		db.LogIssue(
			"juliet_ratio_check",
			"ratio met, cancelling msg",
			fmt.Sprintf("ratio: %f, digits: %d, non-digits: %d, msg: '%s'", ratio, digits, nonDigits, str),
		)

		return true
	}

	return false
}
