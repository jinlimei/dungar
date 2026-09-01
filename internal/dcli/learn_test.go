package dcli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeLearnProtocolMode(t *testing.T) {
	mode, err := normalizeLearnProtocolMode("discord")
	assert.NoError(t, err)
	assert.Equal(t, "discord", mode)

	mode, err = normalizeLearnProtocolMode("Slack")
	assert.NoError(t, err)
	assert.Equal(t, "slack", mode)
}

func TestNormalizeLearnProtocolModeUnknown(t *testing.T) {
	mode, err := normalizeLearnProtocolMode("xmpp")
	assert.Error(t, err)
	assert.Equal(t, "", mode)
}
