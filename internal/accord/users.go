package accord

import (
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
	"log"
)

// SetBotUser assigns the bot user
func (d *Driver) SetBotUser(bu *core2.BotUser) {
	d.botUser = bu
}

// GetUserName returns a users nick based off of their guild membership and user ID
func (d *Driver) GetUserName(userID, serverID string) string {
	guild := d.getOrMakeGuild(serverID)

	member, ok := guild.memberCache[userID]

	if !ok {
		return ""
	}

	if member.Nick != "" {
		return member.Nick
	}

	return member.User.Username
}

// GetBotUser returns the bot user if it is valid or an empty bot user
func (d *Driver) GetBotUser() core2.BotUser {
	if d.botUser != nil {
		return *d.botUser
	}

	return core2.BotUser{}
}

// GetUser translates a discord member to a core2.User
func (d *Driver) GetUser(userID, serverID string) (core2.User, error) {
	if userID == "" {
		return core2.User{}, core2.ErrInvalidUserID
	}

	guild := d.getOrMakeGuild(serverID)

	log.Printf("GetUser(userID='%s',serverID='%s')", userID, serverID)

	if member, ok := guild.memberCache[userID]; ok {
		return d.translateDiscordMember(member), nil
	}

	log.Printf("user '%s' not in membercache (size = %d)", userID, len(guild.memberCache))
	//spew.Dump(guild.memberCache)

	member, err := d.Con.GetSession().GuildMember(serverID, userID)

	if err != nil {
		log.Printf("ERROR: failed to get user from session: %v", err)
		return core2.User{}, err
	}

	guild.memberCache[member.User.ID] = member

	return d.translateDiscordMember(member), nil
}

// GetUsers translate all members to core2.User with their user ID as the key
func (d *Driver) GetUsers(serverID string) map[string]core2.User {
	guild := d.getOrMakeGuild(serverID)
	users := make(map[string]core2.User, len(guild.memberCache))

	for _, member := range guild.memberCache {
		users[member.User.ID] = d.translateDiscordMember(member)
	}

	return users
}
