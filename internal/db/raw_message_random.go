package db

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

// RandomRawMessageForProtocol returns a single stored raw message chosen at random
// for the active protocol mode. Discord uses the same scope as Markov learning:
// the last 30 days, allowlisted channels, and blocked-user exclusion.
func RandomRawMessageForProtocol(protocol string, discordAllowed map[string]struct{}) string {
	protocol = strings.ToLower(strings.TrimSpace(protocol))

	if utils.InTestEnv() {
		return ""
	}

	EnsureDatabaseConnection()

	switch protocol {
	case "discord":
		return randomDiscordRawMessage(discordAllowed)
	case "slack":
		return randomSlackRawMessage()
	default:
		return ""
	}
}

func randomSlackRawMessage() string {
	row := ConQueryRow(`
		SELECT message
		FROM raw_messages
		ORDER BY RANDOM()
		LIMIT 1
	`)

	if row == nil {
		return ""
	}

	var message string
	if err := row.Scan(&message); err != nil {
		return ""
	}

	return message
}

func randomDiscordRawMessage(allowed map[string]struct{}) string {
	query := `
		SELECT message
		FROM raw_messages_discord
	`

	args := make([]any, 0, len(allowed))
	where := make([]string, 0, 3)

	where = append(
		where,
		fmt.Sprintf("created_at >= (CURRENT_TIMESTAMP - INTERVAL '%d days')", utils.DiscordLearningLookbackDays),
	)

	where = append(where, `
		NOT EXISTS (
			SELECT 1
			FROM raw_message_user_blocklist
			WHERE protocol_driver = 'discord'
			  AND server_id = raw_messages_discord.server_id
			  AND user_id = raw_messages_discord.author_user_id
		)
	`)

	channelIDs := make([]string, 0, len(allowed))
	for channelID := range allowed {
		channelIDs = append(channelIDs, channelID)
	}

	sort.Strings(channelIDs)

	if len(channelIDs) > 0 {
		placeholders := make([]string, 0, len(channelIDs))

		for _, channelID := range channelIDs {
			args = append(args, channelID)
			placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)))
		}

		where = append(where, fmt.Sprintf("channel_id IN (%s)", strings.Join(placeholders, ", ")))
	}

	if len(where) > 0 {
		query += `
			WHERE ` + strings.Join(where, " AND ") + `
		`
	}

	query += `
		ORDER BY RANDOM()
		LIMIT 1
	`

	row := ConQueryRow(query, args...)
	if row == nil {
		return ""
	}

	var message string
	if err := row.Scan(&message); err != nil {
		return ""
	}

	return message
}
