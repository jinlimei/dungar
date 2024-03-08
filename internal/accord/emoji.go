package accord

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (d *Driver) handleGuildEmojisUpdate(s *discordgo.Session, ev *discordgo.GuildEmojisUpdate) {
	// We're going to now clear our and recreate our emoji handling
	updated := make(map[string]*discordgo.Emoji, len(ev.Emojis))

	for _, emoji := range ev.Emojis {
		name := strings.ToLower(emoji.Name)

		updated[name] = emoji
	}

	guild := d.getOrMakeGuild(ev.GuildID)
	guild.emojiCache = updated
}

func (d *Driver) retrieveAllGuildEmoji(guildID string) {
	guild := d.getOrMakeGuild(guildID)
	if guild.emojiCache == nil {
		guild.emojiCache = make(map[string]*discordgo.Emoji)
	}

	sess := d.Con.GetSession()
	emojis, err := sess.GuildEmojis(guildID)
	if err != nil {
		log.Printf("ERROR: failed to retrieve guild emojis: %v", err)
		return
	}

	for _, emoji := range emojis {
		name := strings.ToLower(emoji.Name)

		guild.emojiCache[name] = emoji
	}
}
