package taut

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
	"gitlab.int.magneato.site/dungar/prototype/internal/db"
	"gitlab.int.magneato.site/dungar/prototype/internal/utils"
	"gitlab.int.magneato.site/dungar/prototype/library/core2"
)

func (d *Driver) isFromBotUser(ev *slack.MessageEvent) bool {
	botID := d.botUser.ID

	if ev.Msg.User == botID {
		return true
	}

	if ev.SubMessage != nil && ev.SubMessage.User == botID {
		return true
	}

	if ev.PreviousMessage != nil && ev.PreviousMessage.User == botID {
		return true
	}

	return false
}

func (d *Driver) handleMessageUser(ev *slack.MessageEvent) (core2.User, error) {
	if ev.SubMessage != nil {
		return d.handleSubMessage(ev)
	}

	return d.handleMainMessage(ev)
}

func (d *Driver) handleSubMessage(ev *slack.MessageEvent) (core2.User, error) {
	id, ok := utils.CoalesceStr(ev.Msg.User, ev.SubMessage.User)

	if !ok {
		db.LogIssue(
			"invalid_user_id",
			fmt.Sprintf("user_id '%s' from slack msg (%s, %s) resulted in an issue.", id, ev.Msg.Type, ev.Msg.SubType),
			spew.Sdump(ev),
		)

		return core2.User{}, core2.ErrInvalidUserID
	}

	user, err := d.GetUser(id, "")

	// Let us track when the user is doing weirdness.
	if err != nil || user.ID == "" {
		db.LogIssue(
			"unknown_user_id",
			fmt.Sprintf("unknown user_id '%s' in msg (%s, %s)", id, ev.Msg.Type, ev.Msg.SubType),
			spew.Sdump(ev),
		)

		return core2.User{}, core2.ErrUserNotFound
	}

	return user, nil
}

func (d *Driver) handleMainMessage(ev *slack.MessageEvent) (core2.User, error) {
	if ev.Msg.User == "" {
		db.LogIssue(
			"invalid_user_id",
			fmt.Sprintf("user_id '%s' from slack msg (%s, %s) resulted in an issue.", ev.Msg.User, ev.Msg.Type, ev.Msg.SubType),
			spew.Sdump(ev),
		)

		return core2.User{}, core2.ErrInvalidUserID
	}

	user, err := d.GetUser(ev.Msg.User, "")

	if err != nil || user.ID == "" {
		return core2.User{}, core2.ErrUserNotFound
	}

	return user, nil
}
