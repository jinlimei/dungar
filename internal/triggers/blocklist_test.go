package triggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func TestCanUseBlocklistCommandRequiresAdminRole(t *testing.T) {
	utils.WithCICDEnvVars(func() {
		svc := initMockServices()
		mockDriver.SetUser("admin-user", "Admin User")
		mockDriver.AddUserRole("admin-user", "admin")

		assert.True(t, canUseBlocklistCommand(svc, &core2.IncomingMessage{
			UserID:   "admin-user",
			ServerID: "arena",
		}))
		assert.False(t, canUseBlocklistCommand(svc, &core2.IncomingMessage{
			UserID:   "not-admin",
			ServerID: "arena",
		}))
	})
}

func TestBlocklistHandlerAddsUserFromMention(t *testing.T) {
	utils.WithCICDEnvVars(func() {
		svc := initMockServices()
		mockDriver.SetUser("admin-user", "Admin User")
		mockDriver.AddUserRole("admin-user", "admin")
		mockDriver.SetUser("blocked-user", "Furrymatic")

		userID := "blocked-user"
		msg := &core2.IncomingMessage{
			UserID:   "admin-user",
			ServerID: "arena",
			Contents: "!blocklist add @Furrymatic",
			ParsedContents: &core2.ParsedMessage{
				Tokens: []core2.MessageToken{
					{Token: "!blocklist", Type: core2.TokenWord},
					{Token: " ", Type: core2.TokenSpace},
					{Token: "add", Type: core2.TokenWord},
					{Token: " ", Type: core2.TokenSpace},
					{Token: "<@blocked-user>", Type: core2.TokenUserID, Value: &userID},
				},
			},
		}

		rsp := blocklistHandler(svc, msg)

		assert.Len(t, rsp, 1)
		assert.Contains(t, rsp[0].Contents, "Blocklisted @Furrymatic")
		assert.True(t, db.IsRawMessageUserBlocked("discord", "arena", "blocked-user"))
	})
}

func TestBlocklistHandlerRemovesUserFromMention(t *testing.T) {
	utils.WithCICDEnvVars(func() {
		svc := initMockServices()
		mockDriver.SetUser("admin-user", "Admin User")
		mockDriver.AddUserRole("admin-user", "admin")
		mockDriver.SetUser("blocked-user", "Furrymatic")

		userID := "blocked-user"
		assert.NoError(t, db.UpsertRawMessageUserBlock("discord", "arena", userID, "Furrymatic"))
		msg := &core2.IncomingMessage{
			UserID:   "admin-user",
			ServerID: "arena",
			Contents: "!blocklist remove @Furrymatic",
			ParsedContents: &core2.ParsedMessage{
				Tokens: []core2.MessageToken{
					{Token: "!blocklist", Type: core2.TokenWord},
					{Token: " ", Type: core2.TokenSpace},
					{Token: "remove", Type: core2.TokenWord},
					{Token: " ", Type: core2.TokenSpace},
					{Token: "<@blocked-user>", Type: core2.TokenUserID, Value: &userID},
				},
			},
		}

		rsp := blocklistHandler(svc, msg)

		assert.Len(t, rsp, 1)
		assert.Contains(t, rsp[0].Contents, "Removed @Furrymatic")
		assert.False(t, db.IsRawMessageUserBlocked("discord", "arena", userID))
	})
}
