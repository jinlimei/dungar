package accord

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
)

func (d *Driver) handleGuildMemberAddEvent(s *discordgo.Session, ev *discordgo.GuildMemberAdd) {
	guild := d.getOrMakeGuild(ev.GuildID)
	guild.memberCache[ev.User.ID] = ev.Member
}

func (d *Driver) handleGuildMemberUpdateEvent(s *discordgo.Session, ev *discordgo.GuildMemberUpdate) {
	guild := d.getOrMakeGuild(ev.GuildID)
	guild.memberCache[ev.User.ID] = ev.Member
}

func (d *Driver) handleGuildMemberRemoveEvent(s *discordgo.Session, ev *discordgo.GuildMemberRemove) {
	guild := d.getOrMakeGuild(ev.GuildID)
	delete(guild.memberCache, ev.User.ID)
}

func (d *Driver) translateDiscordMember(member *discordgo.Member) core2.User {
	nickName := member.User.Username

	if member.Nick != "" {
		nickName = member.Nick
	}

	log.Printf("translateDiscordMember (nickName='%s')", nickName)

	return core2.User{
		ID:       member.User.ID,
		ServerID: member.GuildID,
		Name:     nickName,
		IsBot:    member.User.Bot || member.User.System,
		IsAdmin:  false,
	}
}
