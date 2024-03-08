package taut

import (
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) getUserReal(id string) (core2.User, error) {
	slackUser, err := d.Con.GetUser(id)

	if err != nil {
		return core2.User{}, err
	}

	d.userCache[id] = slackUser

	return d.translateSlackUser(slackUser), nil
}

func (d *Driver) translateSlackUser(user *slack.User) core2.User {
	name := user.Name

	if user.Profile.DisplayNameNormalized != "" {
		name = user.Profile.DisplayNameNormalized
	}

	return core2.User{
		ID:       user.ID,
		ServerID: d.teamID,
		Name:     name,
		IsBot:    user.IsBot,
		IsAdmin:  user.IsAdmin,
	}
}
