package triggers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var rangeRegexp = regexp.MustCompile("^(@[^: ]+|[^: ]+):? pick a (.+) between ([\\d,.-]+)(?: and |\\s*-\\s*)([\\d,.-]+)\\s*$")

var badMinMaxResponses = []weightedChoice{
	{0.40, "I want you to think on that one buddy"},
	{0.40, "What"},
	{0.10, "wanna try that again?"},
	{0.05, "Fucking WHAT"},
	{0.05, ":thonkeng:"},
}

func rangeHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if !isDirectedAtDungar(svc, msg) {
		return core2.EmptyRsp()
	}

	if !rangeRegexp.MatchString(msg.Contents) {
		return core2.EmptyRsp()
	}

	matches := rangeRegexp.FindStringSubmatch(msg.Contents)

	minStr := strings.Replace(matches[3], ",", "", -1)
	maxStr := strings.Replace(matches[4], ",", "", -1)

	// We're dealing with floats for the float gods.
	if strings.Contains(minStr, ".") || strings.Contains(maxStr, ".") {
		min := mustParseFloat(minStr)
		max := mustParseFloat(maxStr)

		if max <= min {
			return core2.PrefixedSingleRsp(pickWeightedChoice(badMinMaxResponses))
		}

		rsp := fmt.Sprintf("I picked the %s %0.4f", matches[2], random.Float64Range(min, max))

		return core2.PrefixedSingleRsp(rsp)
	}

	min := mustParseInt(minStr)
	max := mustParseInt(maxStr)

	if max <= min {
		return core2.PrefixedSingleRsp(pickWeightedChoice(badMinMaxResponses))
	}

	rsp := fmt.Sprintf("I picked the %s %d", matches[2], random.Int64Range(min, max))
	return core2.PrefixedSingleRsp(rsp)
}

func mustParseInt(s string) int64 {
	out, _ := strconv.ParseInt(s, 10, 64)

	return out
}

func mustParseFloat(s string) float64 {
	out, _ := strconv.ParseFloat(s, 64)

	return out
}
