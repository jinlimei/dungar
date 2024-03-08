package triggers

import (
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func isRepeatableMessage(msg string) bool {
	if len(msg) < 3 {
		return false
	}

	if strings.Contains(msg, " ") {
		return false
	}

	if utils.IsURL(msg) {
		return false
	}

	runes := []rune(msg)

	firstRune := runes[0]
	lastRune := runes[len(runes)-1]

	// skip out on pings and unfurled urls or w/e they're called idk
	if firstRune == '<' && lastRune == '>' {
		return false
	}

	if firstRune != ':' && lastRune != ':' {
		return false
	}

	return true
}

func repeatableHandler(msg *core2.IncomingMessage) []*core2.Response {
	if !isRepeatableMessage(msg.Contents) {
		return core2.EmptyRsp()
	}

	if fromBasicChance("repeatableHandler--literally") {
		return core2.MakeSingleRsp("literally " + msg.Contents)
	}

	if fromBasicChance("repeatableHandler--fingerGuns") {
		return core2.MakeSingleRsp(":point_right: " + msg.Contents + " :point_right:")
	}

	if fromBasicChance("repeatableHandler--basic") {
		return core2.MakeSingleRsp(msg.Contents)
	}

	return core2.EmptyRsp()
}
