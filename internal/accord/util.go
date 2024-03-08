package accord

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	mockConnection *MockDiscordConnection
	mockDriver     *Driver
)

func initMockDriver() *Driver {
	mockConnection = &MockDiscordConnection{
		handlers: make([]DiscordHandlerFn, 0),
		Members:  make(map[string]*discordgo.Member),
		Channels: make(map[string]*discordgo.Channel),
	}

	var err error

	mockDriver, err = New(mockConnection)

	if err != nil {
		log.Fatalf("ERROR: Failed to init mock driver: %v", err)
	}

	mockDriver.registerHandlers()
	mockDriver.guilds["mock"] = &Guild{
		guildID: "mock",
		memberCache: map[string]*discordgo.Member{
			"569251658773692538": makeMockMember("mock", "569251658773692538", "Dungar", ""),
			"62657919375646720":  makeMockMember("mock", "62657919375646720", "jinli.mei", ""),
		},
		roleCache: map[string]*discordgo.Role{
			"766488443248836609": makeMockRole("766488443248836609", "example"),
			"574140080751378432": makeMockRole("574140080751378432", "Dungar"),
		},
		channelCache: map[string]*discordgo.Channel{
			"574133717648408596":  makeMockChannel("mock", "574133717648408596", "general"),
			"1132905045961224313": makeMockChannel("mock", "1132905045961224313", "a message that starts a thread"),
		},
		emojiCache: map[string]*discordgo.Emoji{
			"1133574218303410176": makeMockEmoji("1133574218303410176", "blobyes"),
		},
	}

	return mockDriver
}

func makeMockMember(guildID, userID, userName, nick string) *discordgo.Member {
	return &discordgo.Member{
		GuildID: guildID,
		Nick:    nick,
		User: &discordgo.User{
			ID:       userID,
			Username: userName,
		},
	}
}

func makeMockChannel(guildID, channelID, name string) *discordgo.Channel {
	return &discordgo.Channel{
		ID:      channelID,
		GuildID: guildID,
		Name:    name,
	}
}

func makeMockRole(roleID, name string) *discordgo.Role {
	return &discordgo.Role{
		ID:   roleID,
		Name: name,
	}
}

func makeMockEmoji(emojiID, name string) *discordgo.Emoji {
	return &discordgo.Emoji{
		ID:   emojiID,
		Name: name,
	}
}
