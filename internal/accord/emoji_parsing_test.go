package accord

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDiscordEmoji(t *testing.T) {
	driver := &Driver{}

	parsed := driver.parseDiscordEmoji(":blobyes:1133574218303410176")
	assert.Equal(t, "blobyes", parsed.name)
	assert.Equal(t, "1133574218303410176", parsed.id)

	parsed = driver.parseDiscordEmoji("blobyes:1133574218303410176")
	assert.Equal(t, "blobyes", parsed.name)
	assert.Equal(t, "1133574218303410176", parsed.id)

	parsed = driver.parseDiscordEmoji("blobyes")
	assert.Equal(t, "blobyes", parsed.name)
	assert.Equal(t, "", parsed.id)

	parsed = driver.parseDiscordEmoji(":blobyes:")
	assert.Equal(t, "blobyes", parsed.name)
	assert.Equal(t, "", parsed.id)
}
