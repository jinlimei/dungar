package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordRawDiscordMessageQueryUsesAuthorUserID(t *testing.T) {
	assert.Contains(t, recordRawDiscordMessageQuery, "author_user_id")
	assert.Contains(t, recordRawDiscordMessageQuery, "WITH updated AS")
	assert.Contains(t, recordRawDiscordMessageQuery, "UPDATE raw_messages_discord")
	assert.Contains(t, recordRawDiscordMessageQuery, "author_user_id IS NULL")
	assert.Contains(t, recordRawDiscordMessageQuery, "WHERE message_id = CAST($1 AS varchar(32))")
}

func TestSyncRawDiscordMessageQueryUpdatesMessage(t *testing.T) {
	assert.Contains(t, syncRawDiscordMessageQuery, "UPDATE raw_messages_discord")
	assert.Contains(t, syncRawDiscordMessageQuery, "message = CAST($5 AS text)")
	assert.Contains(t, syncRawDiscordMessageQuery, "INSERT INTO raw_messages_discord")
}
