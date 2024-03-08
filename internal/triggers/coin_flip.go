package triggers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/random"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var coinFlipRegexp = regexp.MustCompile("flip (a coin|([\\d,]+) coins)")

var coinFlipBadRangeResponses = []weightedChoice{
	{0.20, "i have no idea what you want"},
	{0.20, "heads. no wait, tails"},
	{0.20, "have you considered just deciding yourself?"},
	{0.20, "error: i have no coins, am bot"},
	{0.20, "what about"},
}

func coinFlipHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if !isDirectedAtDungar(svc, msg) {
		return core2.EmptyRsp()
	}

	if !coinFlipRegexp.MatchString(msg.Contents) {
		return core2.EmptyRsp()
	}

	var (
		match      = coinFlipRegexp.FindAllStringSubmatch(msg.Contents, 1)
		sideChance = random.Float64Range(0.01, 0.09)
	)

	if match[0][1] == "a coin" {
		var (
			coin = fromBasicChance("coinFlipHandler--singleFlip")
			str  = "tails"
		)

		if coin {
			str = "heads"
		}

		if fromDefinedChance("coinFlipHandler--singleSide", sideChance) {
			str = "oh shit the coin landed on its side"
		}

		return core2.PrefixedSingleRsp(str)
	}

	var (
		strNum      = strings.ReplaceAll(match[0][2], ",", "")
		intNum, err = strconv.Atoi(strNum)
	)

	if err != nil {
		log.Printf("encountered bad number '%s': %v\n", strNum, intNum)
		return core2.PrefixedSingleRsp("i have no idea what you want")
	}

	if intNum <= 0 || intNum >= 1_000_000 {
		return core2.PrefixedSingleRsp(pickWeightedChoice(coinFlipBadRangeResponses))
	}

	var (
		heads = 0
		tails = 0
		sides = 0
	)

	for k := 0; k < intNum; k++ {
		if fromDefinedChance("coinFlipHandler--sides", sideChance) {
			sides++
		} else if fromBasicChance("coinFlipHandler--flip") {
			heads++
		} else {
			tails++
		}
	}

	coinWord := simplePluralize("coin", intNum)
	headWord := simplePluralize("head", heads)
	tailWord := simplePluralize("tail", tails)
	sideWord := simplePluralize("coin", sides)

	response := fmt.Sprintf(
		"flipped %d %s: %d %s, %d %s",
		intNum,
		coinWord,
		heads,
		headWord,
		tails,
		tailWord,
	)

	if sides > 0 {
		response += fmt.Sprintf(
			", and somehow %d %s ended up on a side",
			sides,
			sideWord,
		)
	}

	return core2.PrefixedSingleRsp(response)
}

func simplePluralize(word string, number int) string {
	return pluralize(word, "s", number)
}

func pluralize(word, suffix string, number int) string {
	if number != 1 {
		word += suffix
	}

	return word
}
