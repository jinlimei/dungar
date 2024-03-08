package triggers

import (
	"log"
	"regexp"
	"strings"
)

var phraseRegexp = regexp.MustCompile("\\$([a-zA-Z0-9]+)\\$")

type buildHandlerFn func(serverID string) string

var postBuildHandlers = map[string]buildHandlerFn{
	"randomNoun":     randomNoun,
	"randomFood":     randomFood,
	"randomPlatform": randomPlatform,
	"randomNick":     randomNick,
}

func buildPostMessage(str, serverID string) string {
	if len(str) == 0 {
		return str
	}

	if !phraseRegexp.MatchString(str) {
		return str
	}

	matches := phraseRegexp.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		callback, ok := postBuildHandlers[match[1]]

		// we'll replace and log
		if !ok {
			log.Print("Cannot find callback for '" + match[1] + "'")
			str = strings.Replace(str, match[0], "", -1)
			continue
		}

		str = strings.Replace(str, match[0], callback(serverID), 1)
	}

	return strings.TrimSpace(str)
}
