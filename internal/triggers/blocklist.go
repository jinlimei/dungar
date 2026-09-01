package triggers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func blocklistHandler(svc *core2.Service, msg *core2.IncomingMessage) []*core2.Response {
	txt := strings.TrimSpace(msg.Contents)
	if !strings.HasPrefix(strings.ToLower(txt), "!blocklist") {
		return core2.EmptyRsp()
	}

	if !canUseBlocklistCommand(svc, msg) {
		return core2.EmptyRsp()
	}

	pieces := strings.Fields(txt)
	if len(pieces) < 3 {
		return core2.MakeSingleRsp("Usage: !blocklist <add|remove> @nick")
	}

	command := strings.ToLower(pieces[1])
	if command != "add" && command != "remove" {
		return core2.MakeSingleRsp("Usage: !blocklist <add|remove> @nick")
	}

	if svc.DriverName() != "discord" && !utils.InTestEnv() {
		return core2.MakeSingleRsp("!blocklist is only implemented for Discord right now")
	}

	targetText := strings.TrimSpace(strings.Join(pieces[2:], " "))
	targetUser, err := resolveBlocklistTargetUser(svc, msg, targetText)
	if err != nil {
		return core2.MakeSingleRsp(fmt.Sprintf("Could not resolve blocklist target: %v", err))
	}

	if command == "remove" {
		if err := db.RemoveRawMessageUserBlock("discord", msg.ServerID, targetUser.ID); err != nil {
			return core2.MakeSingleRsp(fmt.Sprintf("Failed to update blocklist: %v", err))
		}

		markovV3RebuildAsync("blocklist-remove-" + targetUser.ID)
		log.Printf(
			"[blocklistHandler] Removed blocklist entry for userID='%s' nick='%s'; Markov rebuild queued",
			targetUser.ID,
			targetUser.Name,
		)

		return core2.MakeSingleRsp(fmt.Sprintf(
			"Removed @%s (%s) from the blocklist; Markov rebuild started.",
			targetUser.Name,
			targetUser.ID,
		))
	}

	if err := db.UpsertRawMessageUserBlock("discord", msg.ServerID, targetUser.ID, targetUser.Name); err != nil {
		return core2.MakeSingleRsp(fmt.Sprintf("Failed to update blocklist: %v", err))
	}

	deleted, err := db.DeleteDiscordRawMessagesByAuthor(msg.ServerID, targetUser.ID)
	if err != nil {
		return core2.MakeSingleRsp(fmt.Sprintf("Failed to delete stored messages: %v", err))
	}

	markovV3RebuildAsync("blocklist-" + targetUser.ID)
	log.Printf(
		"[blocklistHandler] Blocklisted userID='%s' nick='%s' deleted=%d; Markov rebuild queued",
		targetUser.ID,
		targetUser.Name,
		deleted,
	)

	return core2.MakeSingleRsp(fmt.Sprintf(
		"Blocklisted @%s (%s), deleted %d stored messages; Markov rebuild started.",
		targetUser.Name,
		targetUser.ID,
		deleted,
	))
}

func canUseBlocklistCommand(svc *core2.Service, msg *core2.IncomingMessage) bool {
	return canUseDiscordAdminCommand(svc, msg)
}

func resolveBlocklistTargetUser(svc *core2.Service, msg *core2.IncomingMessage, targetText string) (core2.User, error) {
	if svc == nil || msg == nil {
		return core2.User{}, errors.New("missing service or message context")
	}

	if userID := extractBlocklistTargetUserID(msg); userID != "" {
		return resolveUserByIDOrFallback(svc, userID, msg.ServerID, targetText), nil
	}

	targetText = strings.TrimSpace(strings.TrimPrefix(targetText, "@"))
	if targetText == "" {
		return core2.User{}, errors.New("missing target")
	}

	if user, err := resolveUserByNickFromCache(svc, msg.ServerID, targetText); err == nil {
		return user, nil
	}

	tracking, err := db.GetLastSeen(targetText, msg.ServerID)
	if err == nil && tracking != nil {
		return core2.User{
			ID:       tracking.UserID,
			ServerID: msg.ServerID,
			Name:     tracking.Nick,
		}, nil
	}

	return core2.User{}, fmt.Errorf("could not find user '%s'", targetText)
}

func extractBlocklistTargetUserID(msg *core2.IncomingMessage) string {
	return extractTargetUserID(msg)
}

func extractTargetUserID(msg *core2.IncomingMessage, ignoreUserIDs ...string) string {
	if msg == nil || msg.ParsedContents == nil {
		return ""
	}

	ignored := make(map[string]struct{}, len(ignoreUserIDs))
	for _, userID := range ignoreUserIDs {
		userID = strings.TrimSpace(userID)
		if userID != "" {
			ignored[userID] = struct{}{}
		}
	}

	for _, tok := range msg.ParsedContents.IDTokens() {
		if tok.Type == core2.TokenUserID && tok.Value != nil && *tok.Value != "" {
			if _, skip := ignored[*tok.Value]; skip {
				continue
			}

			return *tok.Value
		}
	}

	return ""
}

func resolveUserByIDOrFallback(svc *core2.Service, userID, serverID, fallbackName string) core2.User {
	user, err := svc.GetUser(userID, serverID)
	if err == nil {
		return user
	}

	name := strings.TrimSpace(strings.TrimPrefix(fallbackName, "@"))
	if name == "" {
		name = svc.GetUserName(userID, serverID)
	}
	if name == "" {
		name = userID
	}

	return core2.User{
		ID:       userID,
		ServerID: serverID,
		Name:     name,
	}
}

func resolveUserByNickFromCache(svc *core2.Service, serverID, nick string) (core2.User, error) {
	users := svc.GetUsers(serverID)
	matches := make([]core2.User, 0, 1)

	for _, user := range users {
		if strings.EqualFold(user.Name, nick) || user.ID == nick {
			matches = append(matches, user)
		}
	}

	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return core2.User{}, core2.ErrUserNotFound
	default:
		return core2.User{}, fmt.Errorf("multiple users matched '%s'; use a direct mention", nick)
	}
}
