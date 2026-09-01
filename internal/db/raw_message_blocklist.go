package db

import (
	"fmt"
	"strings"
	"time"

	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
)

type RawMessageUserBlock struct {
	ProtocolDriver string
	ServerID       string
	UserID         string
	Nick           string
	CreatedAt      time.Time
}

var testRawMessageUserBlocks = map[string]RawMessageUserBlock{}

func rawMessageUserBlockKey(protocolDriver, serverID, userID string) string {
	return strings.ToLower(protocolDriver) + "::" + serverID + "::" + userID
}

func UpsertRawMessageUserBlock(protocolDriver, serverID, userID, nick string) error {
	if utils.InTestEnv() {
		testRawMessageUserBlocks[rawMessageUserBlockKey(protocolDriver, serverID, userID)] = RawMessageUserBlock{
			ProtocolDriver: strings.ToLower(protocolDriver),
			ServerID:       serverID,
			UserID:         userID,
			Nick:           nick,
			CreatedAt:      time.Now(),
		}

		return nil
	}

	qry := `
		INSERT INTO raw_message_user_blocklist(protocol_driver, server_id, user_id, nick, created_at)
		VALUES(
			CAST($1 AS varchar(20)),
			CAST($2 AS varchar(32)),
			CAST($3 AS varchar(32)),
			CAST($4 AS varchar(100)),
			CURRENT_TIMESTAMP
		)
		ON CONFLICT (protocol_driver, server_id, user_id) DO UPDATE
		SET nick = EXCLUDED.nick
	`

	_, err := ConExec(qry, strings.ToLower(protocolDriver), serverID, userID, nick)
	return err
}

func RemoveRawMessageUserBlock(protocolDriver, serverID, userID string) error {
	if utils.InTestEnv() {
		delete(testRawMessageUserBlocks, rawMessageUserBlockKey(protocolDriver, serverID, userID))
		return nil
	}

	qry := `
		DELETE FROM raw_message_user_blocklist
		WHERE protocol_driver = $1
		  AND server_id = $2
		  AND user_id = $3
	`

	_, err := ConExec(qry, strings.ToLower(protocolDriver), serverID, userID)
	return err
}

func IsRawMessageUserBlocked(protocolDriver, serverID, userID string) bool {
	if userID == "" {
		return false
	}

	if utils.InTestEnv() {
		_, ok := testRawMessageUserBlocks[rawMessageUserBlockKey(protocolDriver, serverID, userID)]
		return ok
	}

	qry := `
		SELECT 1
		FROM raw_message_user_blocklist
		WHERE protocol_driver = $1
		  AND server_id = $2
		  AND user_id = $3
	`

	row := ConQueryRow(qry, strings.ToLower(protocolDriver), serverID, userID)
	if row == nil {
		return false
	}

	var exists int
	if err := row.Scan(&exists); err != nil {
		return false
	}

	return exists == 1
}

func DeleteDiscordRawMessagesByAuthor(serverID, userID string) (int64, error) {
	if userID == "" {
		return 0, nil
	}

	if utils.InTestEnv() {
		return 0, nil
	}

	qry := `
		DELETE FROM raw_messages_discord
		WHERE server_id = $1
		  AND author_user_id = $2
	`

	res, err := ConExec(qry, serverID, userID)
	if err != nil {
		return 0, err
	}

	if res == nil {
		return 0, nil
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %w", err)
	}

	return rows, nil
}
