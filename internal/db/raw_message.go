package db

import (
	"database/sql"
	"strings"
	"time"
)

const recordRawDiscordMessageQuery = `
	WITH updated AS (
		UPDATE raw_messages_discord
		SET author_user_id = NULLIF(CAST($4 AS varchar(32)), '')
		WHERE message_id = CAST($1 AS varchar(32))
		  AND author_user_id IS NULL
		  AND NULLIF(CAST($4 AS varchar(32)), '') IS NOT NULL
		RETURNING 1
	)
	INSERT INTO raw_messages_discord(message_id, server_id, channel_id, author_user_id, message, created_at)
	SELECT
		CAST($1 AS varchar(32)),
		CAST($2 AS varchar(32)),
		CAST($3 AS varchar(32)),
		NULLIF(CAST($4 AS varchar(32)), ''),
		CAST($5 AS text),
		CAST($6 AS timestamp with time zone)
	WHERE NOT EXISTS (
		SELECT 1
		FROM raw_messages_discord
		WHERE message_id = CAST($1 AS varchar(32))
	)
	  AND NOT EXISTS (
		SELECT 1
		FROM updated
	)
`

const syncRawDiscordMessageQuery = `
	WITH updated AS (
		UPDATE raw_messages_discord
		SET server_id = CAST($2 AS varchar(32)),
		    channel_id = CAST($3 AS varchar(32)),
		    author_user_id = COALESCE(NULLIF(CAST($4 AS varchar(32)), ''), author_user_id),
		    message = CAST($5 AS text)
		WHERE message_id = CAST($1 AS varchar(32))
		RETURNING 1
	)
	INSERT INTO raw_messages_discord(message_id, server_id, channel_id, author_user_id, message, created_at)
	SELECT
		CAST($1 AS varchar(32)),
		CAST($2 AS varchar(32)),
		CAST($3 AS varchar(32)),
		NULLIF(CAST($4 AS varchar(32)), ''),
		CAST($5 AS text),
		CAST($6 AS timestamp with time zone)
	WHERE NOT EXISTS (
		SELECT 1
		FROM updated
	)
`

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
func RecordRawDiscordMessage(messageID, serverID, channelID, authorUserID, message string) bool {
	return RecordRawDiscordMessageAt(
		messageID,
		serverID,
		channelID,
		authorUserID,
		message,
		time.Now().UTC(),
	)
}

// RecordRawDiscordMessageAt records a Discord message at a specific timestamp.
// This is used by history backfills so the stored row reflects the message's
// original creation time instead of the import time. Returns true when a row
// was inserted or an existing row was enriched with a newly known author ID.
func RecordRawDiscordMessageAt(messageID, serverID, channelID, authorUserID, message string, createdAt time.Time) bool {
	res := ConMustExec(
		recordRawDiscordMessageQuery,
		messageID,
		serverID,
		channelID,
		authorUserID,
		message,
		createdAt,
	)
	return rowsAffectedAtLeastOne("RecordRawDiscordMessageAt", res)
}

// SyncRawDiscordMessage updates an existing stored Discord message or inserts it
// if it did not previously exist.
func SyncRawDiscordMessage(messageID, serverID, channelID, authorUserID, message string, createdAt time.Time) bool {
	res := ConMustExec(
		syncRawDiscordMessageQuery,
		messageID,
		serverID,
		channelID,
		authorUserID,
		message,
		createdAt,
	)
	return rowsAffectedAtLeastOne("SyncRawDiscordMessage", res)
}

// DeleteRawDiscordMessage removes a stored Discord message by message ID.
func DeleteRawDiscordMessage(messageID string) bool {
	if strings.TrimSpace(messageID) == "" {
		return false
	}

	res := ConMustExec(`
		DELETE FROM raw_messages_discord
		WHERE message_id = CAST($1 AS varchar(32))
	`, messageID)

	return rowsAffectedAtLeastOne("DeleteRawDiscordMessage", res)
}

func rowsAffectedAtLeastOne(loc string, res sql.Result) bool {
	if res == nil {
		return false
	}

	rows, err := res.RowsAffected()
	handleError(loc+".RowsAffected", err)
	return rows > 0
}
