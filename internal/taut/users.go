package taut

import (
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

// GetUserName returns a user name based off of their user id
func (d *Driver) GetUserName(userID, _ string) string {
	return d.Con.GetUserName(userID)
}

// GetBotUser retrieves the active bot user from the connection
func (d *Driver) GetBotUser() core2.BotUser {
	if d.botUser != nil {
		return *d.botUser
	}

	return core2.BotUser{}
}

// GetUser attempts to return a user struct and bool of success
func (d *Driver) GetUser(userID, _ string) (core2.User, error) {
	if userID == "" {
		return core2.User{}, core2.ErrInvalidUserID
	}

	if user, ok := d.userCache[userID]; ok {
		return d.translateSlackUser(user), nil
	}

	return d.getUserReal(userID)
}

// GetUsers returns a map of all Users with their UIDs as the key
func (d *Driver) GetUsers(_ string) map[string]core2.User {
	users := d.Con.GetUsers()
	out := make(map[string]core2.User, len(users))

	for _, user := range users {
		out[user.ID] = d.translateSlackUser(user)
	}

	return out
}

// SetBotUser sets the bot user
func (d *Driver) SetBotUser(bu *core2.BotUser) {
	d.botUser = bu
}
