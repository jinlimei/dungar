package triggers

import (
	"regexp"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var negativeMentionSnarkChoices = []weightedChoice{
	{0.10, "no u"},
	{0.10, "cope"},
	{0.10, "skill issue"},
	{0.10, "cry about it"},
	{0.10, "didn't ask"},
	{0.10, "ok and?"},
	{0.10, "ratio"},
	{0.10, "touch grass"},
	{0.10, "keep coping"},
	{0.10, "mad cuz bad"},
}

var negativeMentionSnarkRegex = regexp.MustCompile(`\b(` +
	`bad bot|shut up|stfu|fuck you|f you|you suck|suck|annoying|hate you|dumb|idiot|garbage|trash|worst` +
	`)\b`)

func negativeMentionSnarkHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	if svc == nil || msg == nil || msg.UserID == "" {
		return core2.EmptyRsp()
	}

	if svc.GetBotUser().ID != "" && msg.UserID == svc.GetBotUser().ID {
		return core2.EmptyRsp()
	}

	if strings.HasPrefix(strings.TrimSpace(msg.Contents), "!") {
		return core2.EmptyRsp()
	}

	if !isMentioningDungar(svc, msg) && !isDirectedAtDungar(svc, msg) {
		return core2.EmptyRsp()
	}

	contents := msg.Contents
	if isDirectedAtDungar(svc, msg) {
		contents = normalizeDirectedContents(svc, msg.ServerID, contents, svc.GetBotUser())
	}

	contents = strings.ToLower(utils.CleanSpaces(contents))
	if !negativeMentionSnarkRegex.MatchString(contents) {
		return core2.EmptyRsp()
	}

	if !fromBasicChance("negativeMentionSnarkHandler--respond") {
		return core2.EmptyRsp()
	}

	return core2.MakeSingleRsp(pickWeightedChoice(negativeMentionSnarkChoices))
}
