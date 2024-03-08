package triggers

import (
	"regexp"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

var (
	questionStartRegex = regexp.MustCompile("^(@[^ ]+:?|[^ ]+:)\\s+.+\\??")
)

// isDirectedAtDungar provides a direct yes/no for being directed at dungar
// criteria requires it be the start of the message.
// with no changes of other scenarios (no RNG to decide things)
func isDirectedAtDungar(svc *core2.Service, msg *core2.IncomingMessage) bool {
	// Get the bot name
	bot := svc.GetBotUser()

	if bot.ID == "" || bot.Name == "" {
		return false
	}

	if msg.ChannelID == bot.ID {
		return true
	}

	return checkDirectedPrefix(msg.Contents, bot.Name)
}

// isMentioningDungar checks if dungar is being mentioned in a message, not
// just at the beginning of the message but anywhere else.
func isMentioningDungar(svc *core2.Service, msg *core2.IncomingMessage) bool {
	bot := svc.GetBotUser()

	if bot.ID == "" || bot.Name == "" {
		return false
	}

	if msg.ChannelID == bot.ID {
		return true
	}

	return strings.Contains(
		strings.ToLower(msg.Contents),
		strings.ToLower(bot.Name),
	)
}

// isDirectedAtDungarRandom provides a chance for the message to turn
// randomly true even if it's not being directed at dungar
func isDirectedAtDungarRandom(svc *core2.Service, msg *core2.IncomingMessage, chance float64) bool {
	if fromDefinedChance("directedAtDungar--random", chance) {
		return true
	}

	return isDirectedAtDungar(svc, msg)
}

// isMentioningDungarRandom provides a chance for the message to turn
// randomly true even if it's not mentioning dungar
func isMentioningDungarRandom(svc *core2.Service, msg *core2.IncomingMessage, chance float64) bool {
	if fromDefinedChance("mentioningDungar--random", chance) {
		return true
	}

	return isMentioningDungar(svc, msg)
}

func extractUserNameTarget(str string, reg *regexp.Regexp) string {
	matches := reg.FindStringSubmatch(str)

	if len(matches) < 2 {
		return ""
	}

	return strings.Trim(matches[1], "@:")
}

func checkDirectedPrefix(msg, name string) bool {
	msg = strings.TrimSpace(msg)
	name = strings.TrimSpace(name)

	var (
		msgRunes  = []rune(msg)
		nameRunes = []rune(name)
		msgRLen   = len(msgRunes)
		nameRLen  = len(nameRunes)

		pos        = 0
		firstSpace = 0
	)

	if msgRLen < nameRLen {
		return false
	}

	// if first char is a '@'
	if msgRunes[0] == '@' {
		pos++
		firstSpace++
	}

	for ; firstSpace < msgRLen; firstSpace++ {
		if msgRunes[firstSpace] == ' ' {
			break
		}
	}

	if firstSpace > 0 && msgRunes[firstSpace-1] == ':' {
		firstSpace--
	}

	word := string(msgRunes[pos:firstSpace])
	//fmt.Printf("word: '%s', name: '%s'\n", word, name)

	if strings.ToLower(word) == strings.ToLower(name) {
		return true
	}

	return false
}
