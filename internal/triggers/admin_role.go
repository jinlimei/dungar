package triggers

import (
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

const discordAdminRoleName = "admin"

type roleChecker interface {
	HasRole(userID, serverID, roleName string) bool
}

func canUseDiscordAdminCommand(svc *core2.Service, msg *core2.IncomingMessage) bool {
	if svc == nil || msg == nil || msg.UserID == "" || msg.ServerID == "" {
		return false
	}

	checker, ok := svc.Driver().(roleChecker)
	if !ok {
		return utils.InTestEnv()
	}

	return checker.HasRole(msg.UserID, msg.ServerID, discordAdminRoleName)
}

func commandVerbMatches(input, expected string) bool {
	return strings.EqualFold(strings.TrimSpace(input), strings.TrimSpace(expected))
}
