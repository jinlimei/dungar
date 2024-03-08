package accord

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
	"strings"
)

func (d *Driver) handleChannelCreateEvent(s *discordgo.Session, ev *discordgo.ChannelCreate) {
	guild := d.getOrMakeGuild(ev.GuildID)
	guild.channelCache[ev.Channel.ID] = ev.Channel

	d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.ID,
		UserID:    ev.OwnerID,
		ServerID:  ev.GuildID,
		ChannelID: ev.Channel.ID,
		Text:      ev.Channel.Name,
		Type:      core2.EventChannelJoin,
	})
}

func (d *Driver) handleChannelUpdateEvent(s *discordgo.Session, ev *discordgo.ChannelUpdate) {
	guild := d.getOrMakeGuild(ev.GuildID)
	guild.channelCache[ev.Channel.ID] = ev.Channel
}

func (d *Driver) handleChannelDeleteEvent(s *discordgo.Session, ev *discordgo.ChannelDelete) {
	guild := d.getOrMakeGuild(ev.GuildID)
	delete(guild.channelCache, ev.Channel.ID)

	d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.ID,
		UserID:    ev.OwnerID,
		ServerID:  ev.GuildID,
		ChannelID: ev.Channel.ID,
		Text:      ev.Channel.Name,
		Type:      core2.EventChannelLeave,
	})
}

func (d *Driver) handleThreadCreateEvent(s *discordgo.Session, ev *discordgo.ThreadCreate) {
	guild := d.getOrMakeGuild(ev.GuildID)
	guild.channelCache[ev.Channel.ID] = ev.Channel

	d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.ID,
		UserID:    ev.OwnerID,
		ServerID:  ev.GuildID,
		ChannelID: ev.Channel.ID,
		Text:      ev.Channel.Name,
		Type:      core2.EventChannelJoin,
	})
}

func (d *Driver) handleThreadUpdateEvent(s *discordgo.Session, ev *discordgo.ThreadUpdate) {
	guild := d.getOrMakeGuild(ev.GuildID)
	guild.channelCache[ev.Channel.ID] = ev.Channel
}

func (d *Driver) handleThreadDeleteEvent(s *discordgo.Session, ev *discordgo.ThreadDelete) {
	guild := d.getOrMakeGuild(ev.GuildID)
	delete(guild.channelCache, ev.Channel.ID)

	d.core.HandleIncomingEvent(&core2.IncomingEvent{
		ID:        ev.ID,
		UserID:    ev.OwnerID,
		ServerID:  ev.GuildID,
		ChannelID: ev.Channel.ID,
		Text:      ev.Channel.Name,
		Type:      core2.EventChannelLeave,
	})
}

func (d *Driver) getCoreChannels() map[string]core2.Channel {
	channels := d.guilds.getChannels()
	converted := make(map[string]core2.Channel, len(channels))

	for _, channel := range channels {
		converted[channel.ID] = d.translateDiscordChannel(channel)
	}

	return converted
}

func (d *Driver) translateDiscordChannel(channel *discordgo.Channel) core2.Channel {
	return core2.Channel{
		ID:             channel.ID,
		ServerID:       channel.GuildID,
		Name:           channel.Name,
		NameNormalized: strings.ToLower(channel.Name),
		Topic:          channel.Topic,
		Type:           d.translateDiscordChannelType(channel.Type),
		Archetype:      d.inferChannelArchetype(channel),
		PreviousNames:  []string{},
	}
}

func (d *Driver) inferChannelArchetype(channel *discordgo.Channel) string {
	if channel.IsThread() {
		return "thread"
	}

	if channel.DefaultForumLayout != 0 {
		return "forum"
	}

	return "public"
}

func (d *Driver) translateDiscordChannelType(chanType discordgo.ChannelType) core2.ChannelType {
	switch chanType {
	case discordgo.ChannelTypeDM:
		return core2.ChannelDirectMessage

	case discordgo.ChannelTypeGroupDM:
		return core2.ChannelGroupMessage

	case discordgo.ChannelTypeGuildText:
		return core2.ChannelPublic

	case discordgo.ChannelTypeGuildPublicThread:
		return core2.ChannelPublic

	case discordgo.ChannelTypeGuildForum:
		return core2.ChannelForum

	case discordgo.ChannelTypeGuildPrivateThread:
		return core2.ChannelPrivileged

	case discordgo.ChannelTypeGuildNewsThread:
		return core2.ChannelReadOnly

	case discordgo.ChannelTypeGuildCategory:
		return core2.ChannelReadOnly

	case discordgo.ChannelTypeGuildNews:
		return core2.ChannelReadOnly

	case discordgo.ChannelTypeGuildStore:
		return core2.ChannelReadOnly

	case discordgo.ChannelTypeGuildVoice, discordgo.ChannelTypeGuildStageVoice:
		return core2.ChannelVoice
	}

	log.Printf("Unknown channel type '%d' provided", chanType)
	panic("unknown channel type provided")
}
