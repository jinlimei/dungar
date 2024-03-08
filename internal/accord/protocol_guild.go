package accord

import "github.com/bwmarrin/discordgo"

// Guild manages a specific discord guild's member/channel/emoji structs
type Guild struct {
	guildID      string
	memberCache  map[string]*discordgo.Member
	channelCache map[string]*discordgo.Channel
	roleCache    map[string]*discordgo.Role
	emojiCache   map[string]*discordgo.Emoji
}

func (g *Guild) init() {
	g.memberCache = make(map[string]*discordgo.Member, 70)
	g.channelCache = make(map[string]*discordgo.Channel, 70)
	g.roleCache = make(map[string]*discordgo.Role, 70)
	g.emojiCache = make(map[string]*discordgo.Emoji, 70)
}

// GuildManager manages the set of guilds that the bot is actively connected to
type GuildManager map[string]*Guild

func (gm GuildManager) getChannels() map[string]*discordgo.Channel {
	channels := make(map[string]*discordgo.Channel)

	for _, guild := range gm {
		for _, channel := range guild.channelCache {
			channels[channel.ID] = channel
		}
	}

	return channels
}

func (gm GuildManager) getChannelByID(channelID string) (*discordgo.Channel, bool) {
	for _, guild := range gm {
		channel, ok := guild.channelCache[channelID]
		if ok {
			return channel, true
		}
	}

	return nil, false
}

func (gm GuildManager) getUserByID(userID string) []*discordgo.Member {
	user := make([]*discordgo.Member, 0)

	for _, guild := range gm {
		if member, ok := guild.memberCache[userID]; ok {
			user = append(user, member)
		}
	}

	return user
}
