package db

import (
	"strings"
)

// LegacyRecordRawMessage handles recording/ingesting messages that are
// completely and totally raw.
func LegacyRecordRawMessage(msg, source string) {
	source = strings.ToLower(source)

	sql := `
		INSERT INTO raw_messages (message, source, created_at)
		VALUES($1, $2, CURRENT_TIMESTAMP)
	`

	ConMustExec(sql, msg, source)
}

// RecordRawDiscordMessage is intended to record discord messages, rawly
func RecordRawDiscordMessage(messageID, serverID, channelID, message string) {
	qry := `
		INSERT INTO raw_messages_discord(message_id, server_id, channel_id, message, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`

	ConMustExec(qry, messageID, serverID, channelID, message)
}
